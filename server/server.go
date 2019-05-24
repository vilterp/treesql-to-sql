package server

import (
	"bufio"
	"database/sql"
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

	stmt, err := parse.Parse(string(query))
	if err != nil {
		log.Println("parse error", err)
		w.Write([]byte(fmt.Sprintf("parse error: %v", err)))
		http.Error(w, fmt.Sprintf("parse error"), http.StatusBadRequest)
	}

	sqlQuery := querygen.Generate(stmt.Select)

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
