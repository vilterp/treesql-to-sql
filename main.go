package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/envy"
	"github.com/vilterp/treesql-to-sql/server"
	"github.com/vilterp/treesql-to-sql/util"
)

func main() {
	connParams := envy.Get("CONN_PARAMS", "user=root dbname=blog sslmode=disable port=26257")

	host := envy.Get("HOST", "0.0.0.0")
	port := envy.Get("PORT", "9001")
	addr := fmt.Sprintf("%s:%s", host, port)

	s, err := server.NewServer(connParams)
	if err != nil {
		log.Fatal("couldn't start up:", err)
	}

	log.Printf("starting server at http://%s/", addr)
	if err := http.ListenAndServe(addr, util.Logger(s)); err != nil {
		log.Fatal(err)
	}
}
