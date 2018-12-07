package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	conn *sqlx.DB
	mux  *http.ServeMux
}

func NewServer(connParams string) (*Server, error) {
	conn, err := sqlx.Connect("postgres", connParams)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	s := &Server{
		conn: conn,
		mux:  mux,
	}

	mux.Handle("/sql", http.HandlerFunc(s.serveSQL))

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

	cols, err := rows.Columns()
	if err != nil {
		log.Println("error getting cols:", err)
		w.WriteHeader(500)
		return
	}

	var rowsOut [][]interface{}
	for rows.Next() {
		row := make([]interface{}, len(cols))

		if err := rows.Scan(row...); err != nil {
			log.Println("error scanning:", err)
			w.WriteHeader(500)
			return
		}
		rowsOut = append(rowsOut, row)
	}

	res := map[string]interface{}{
		"cols": cols,
		"rows": rowsOut,
	}

	bytes, err := json.Marshal(res)
	if err != nil {
		log.Println("error writing response:", err)
		w.WriteHeader(500)
		return
	}

	// TODO: buffer output???
	if n, err := w.Write(bytes); err != nil {
		log.Println("error writing result:", err)
		w.WriteHeader(500)
		return
	} else {
		fmt.Println("la")
		log.Printf("wrote %d bytes", n)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}
