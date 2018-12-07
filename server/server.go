package server

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/vilterp/treesql-to-sql/parserlib"
)

type Server struct {
	conn    *sql.DB
	mux     *http.ServeMux
	grammar *parserlib.Grammar
}

func NewServer(connParams string) (*Server, error) {
	conn, err := sql.Open("postgres", connParams)
	if err != nil {
		return nil, err
	}

	gram, err := parserlib.TestTreeSQLGrammar()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	s := &Server{
		conn:    conn,
		mux:     mux,
		grammar: gram,
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

	traceTree, err := s.grammar.Parse("select", string(query))
	if err != nil {
		// TODO(vilterp): really need to wrap this so I can just return an error
		msg := fmt.Sprintf("parsing query: %v", err)
		log.Println("error", msg)
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Println("error writing error:", err)
		}
		return
	}

	log.Println("query:", traceTree.Format().String())

	// TODO: generate
	sqlQuery := "select json_agg(json_build_object('id', id, 'name', name, 'created_at', created_at)) FROM clusters"

	rows, err := s.conn.Query(string(sqlQuery))
	if err != nil {
		log.Println("error running query:", err)
		w.WriteHeader(http.StatusBadRequest)
		// TODO(vilterp): return this as JSON
		if _, err := w.Write([]byte(fmt.Sprintf("error running query: %v", err))); err != nil {
			log.Println("error writing error:", err)
		}
		return
	}

	var out string
	rows.Next()
	if err := rows.Scan(&out); err != nil {
		log.Println("error scanning rows:", err)
		w.WriteHeader(500)
		return
	}

	// TODO(vilterp): isn't there a stdlib method that just writes an entire string to a writer? jeez.
	b := bufio.NewWriter(w)
	if _, err := b.WriteString(out); err != nil {
		log.Println("error writing result:", err)
		w.WriteHeader(500)
		return
	}
	if err := b.Flush(); err != nil {
		log.Println("error writing result:", err)
		w.WriteHeader(500)
		return
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}
