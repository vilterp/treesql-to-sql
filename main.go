package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/envy"
	sf "github.com/snowflakedb/gosnowflake"
	. "github.com/vilterp/treesql-to-sql/live_queries"
	"github.com/vilterp/treesql-to-sql/server"
	"github.com/vilterp/treesql-to-sql/snowflake"
	"github.com/vilterp/treesql-to-sql/util"
)

func main() {
	connParams := envy.Get("CONN_PARAMS", "user=root dbname=management_console_production sslmode=disable port=26257")
	snowflakePw := envy.Get("SNOWFLAKE_PW", "")

	host := envy.Get("HOST", "0.0.0.0")
	port := envy.Get("PORT", "9001")
	addr := fmt.Sprintf("%s:%s", host, port)

	sl, err := snowflake.NewSnowflakeListener(&sf.Config{
		Account:  "mv63954",
		Region:   "us-east-1",
		User:     "vilterp",
		Password: snowflakePw,
		Database: "console",
	}, "audit_log")
	if err != nil {
		panic(err)
	}

	s, err := server.NewServer(connParams, []*server.ListenerSpec{
		{
			TableName: "audit_log",
			Name:      "log",
			Listener: &Listener{
				Insert: func(r Row) error {
					log.Println("listener: insert row:", r)
					return nil
				},
			},
		},
		{
			TableName: "audit_log",
			Name:      "snowflake",
			Listener:  sl,
		},
	})
	if err != nil {
		log.Fatal("couldn't start up:", err)
	}

	log.Printf("starting server at http://%s/", addr)
	if err := http.ListenAndServe(addr, util.Logger(s)); err != nil {
		log.Fatal(err)
	}
}
