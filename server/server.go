package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/vilterp/go-parserlib/examples/treesql"
	"github.com/vilterp/treesql-to-sql/live_queries"
	"github.com/vilterp/treesql-to-sql/querygen"
	"github.com/vilterp/treesql-to-sql/schema"
)

type Server struct {
	conn   *sql.DB
	mux    *http.ServeMux
	schema schema.Schema
}

func NewServer(connParams string) (*Server, error) {
	conn, err := sql.Open("postgres", connParams)
	if err != nil {
		return nil, err
	}

	dbSchema, err := schema.LoadSchema(conn)
	if err != nil {
		return nil, err
	}

	log.Println("creating live query schema...")
	if err := live_queries.CreateSchema(conn); err != nil {
		return nil, err
	}
	log.Println("created live query schema")

	mux := http.NewServeMux()

	s := &Server{
		conn:   conn,
		mux:    mux,
		schema: dbSchema,
	}

	mux.Handle("/query", http.HandlerFunc(s.serveSQL))
	mux.Handle("/schema", http.HandlerFunc(s.serveSchema))
	mux.Handle("/validate", http.HandlerFunc(s.serveValidate))
	mux.Handle("/", http.FileServer(http.Dir("ui/build")))

	events, err := live_queries.LiveQuery(conn, dbSchema)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			evt := <-events
			log.Printf("changefeed event %#v", evt)
		}
	}()

	return s, nil
}

func (s *Server) serveSQL(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")

	query, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("error reading body:", err)
		w.WriteHeader(500)
		return
	}

	out, qErr := s.runQuery(string(query))
	if qErr != nil {
		// wtf, why does this not let me get the code
		http.Error(w, qErr.Error(), qErr.code)
		log.Println("error running query:", qErr.Error())
		return
	}

	// TODO(vilterp): isn't there a stdlib method that just writes an entire string to a writer? jeez.
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(out); err != nil {
		log.Println("error writing result:", err)
		w.WriteHeader(500)
		return
	}
}

type queryError struct {
	msg  string
	code int
}

func mkQueryError(code int, msg string) *queryError {
	return &queryError{
		msg:  msg,
		code: code,
	}
}

func (qe *queryError) Error() string {
	return qe.msg
}

type QueryResult struct {
	Res              string
	SQL              string
	FormattedTreeSQL string
}

func (s *Server) runQuery(query string) (*QueryResult, *queryError) {
	plLang := treesql.MakeLanguage(s.schemaForLanguage())

	psiTree, err := plLang.Parse(query)
	if err != nil {
		log.Println("parse error:", err)
		return nil, &queryError{
			code: http.StatusBadRequest,
			msg:  fmt.Sprintf("parse error: %v", err),
		}
	}

	sqlQuery, err := querygen.Generate(psiTree.(*treesql.Select), s.schema)
	if err != nil {
		return nil, mkQueryError(http.StatusBadRequest, fmt.Sprintf("generating query: %v", err.Error()))
	}

	log.Println("SQL query", sqlQuery)

	queryStartTime := time.Now()
	rows, err := s.conn.Query(string(sqlQuery))
	queryEndTime := time.Now()
	if err != nil {
		log.Println("error running query:", err)
		return nil, mkQueryError(http.StatusInternalServerError, fmt.Sprintf("running query: %v", err.Error()))
	}
	log.Println("query time:", queryEndTime.Sub(queryStartTime))

	var out string
	rows.Next()
	if err := rows.Scan(&out); err != nil {
		return nil, mkQueryError(http.StatusInternalServerError, fmt.Sprintf("scanning rows: %v", err.Error()))
	}

	return &QueryResult{
		Res: out,
		SQL: sqlQuery,
		//FormattedTreeSQL: stmt.Pretty().String(),
	}, nil
}

func (s *Server) schemaForLanguage() *treesql.SchemaDesc {
	sd := &treesql.SchemaDesc{
		Tables: map[string]*treesql.TableDesc{},
	}
	for name := range s.schema {
		tbl := &treesql.TableDesc{
			Columns: map[string]*treesql.ColDesc{},
		}
		sd.Tables[name] = tbl
	}
	return sd
}

func (s *Server) serveSchema(w http.ResponseWriter, req *http.Request) {
	schemaDesc := s.schemaForLanguage()
	bytes, err := json.Marshal(schemaDesc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err := w.Write(bytes); err != nil {
		log.Println(err)
	}
}

func (s *Server) serveValidate(w http.ResponseWriter, req *http.Request) {
	l := treesql.MakeLanguage(s.schemaForLanguage())
	l.ServeCompletions(w, req)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}
