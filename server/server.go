package server

import (
	"bufio"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

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
	query, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("error reading body:", err)
		w.WriteHeader(500)
		return
	}

	rows, err := s.conn.Query(string(query))
	if err != nil {
		log.Println("error running query:", err)
		w.WriteHeader(500)
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
