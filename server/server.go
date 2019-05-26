package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vilterp/treesql-to-sql/querygen"

	"github.com/vilterp/treesql-to-sql/parse"

	_ "github.com/lib/pq"
)

type Server struct {
	conn *sql.DB
	mux  *http.ServeMux
}

func NewServer(connParams string) (*Server, error) {
	conn, err := sql.Open("postgres", connParams)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	s := &Server{
		conn: conn,
		mux:  mux,
	}

	mux.Handle("/query", http.HandlerFunc(s.serveSQL))
	mux.Handle("/", http.FileServer(http.Dir("ui/build")))

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

	log.Println(out.SQL)

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
	Res string
	SQL string
}

func (s *Server) runQuery(query string) (*QueryResult, *queryError) {
	stmt, err := parse.Parse(string(query))
	if err != nil {
		log.Println("parse error", err)
		return nil, mkQueryError(http.StatusBadRequest, fmt.Sprintf("parse error: %v", err))
	}

	sqlQuery, err := querygen.Generate(stmt.Select)
	if err != nil {
		return nil, mkQueryError(http.StatusBadRequest, fmt.Sprintf("generating query: %v", err.Error()))
	}

	rows, err := s.conn.Query(string(sqlQuery))
	if err != nil {
		log.Println("error running query:", err)
		return nil, mkQueryError(http.StatusInternalServerError, fmt.Sprintf("running query: %v", err.Error()))
	}

	var out string
	rows.Next()
	if err := rows.Scan(&out); err != nil {
		return nil, mkQueryError(http.StatusInternalServerError, fmt.Sprintf("scanning rows: %v", err.Error()))
	}

	return &QueryResult{
		Res: out,
		SQL: sqlQuery,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}
