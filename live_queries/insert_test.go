package live_queries

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	"github.com/vilterp/treesql-to-sql/parse"
)

func TestInsert(t *testing.T) {
	conn, err := sql.Open("postgres", "user=root dbname=treesql_live_queries sslmode=disable port=26257")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec("DROP DATABASE treesql_live_queries"); err != nil {
		t.Fatal(err)
	}

	if err := CreateSchema(conn); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.Exec(
		"INSERT INTO connections (id, remote_addr) VALUES ($1, $2)",
		"8985c599-a6a9-4ae7-beab-a73ef30085ba", "1.2.3.4",
	); err != nil {
		t.Fatal(err)
	}

	queryText := "MANY posts { id, comments: MANY comments { id } }"
	parsedQuery, err := parse.Parse(queryText)
	if err != nil {
		t.Fatal(err)
	}
	queryCtx := &QueryContext{
		connID:    "8985c599-a6a9-4ae7-beab-a73ef30085ba",
		queryText: queryText,
		conn:      conn,
		query:     parsedQuery.Select,
	}
	err = InsertListeners(queryCtx)
	if err != nil {
		t.Fatal(err)
	}
}
