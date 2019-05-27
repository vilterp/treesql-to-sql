package live_queries

import (
	"database/sql"
	"fmt"
)

var schema = []string{
	`CREATE DATABASE IF NOT EXISTS treesql_live_queries;`,

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.listeners (
		id          uuid PRIMARY KEY,
		table_name  text,
		column_name text,
		value       text,
		rowid       int default unique_rowid() not null
	);`,

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.connections (
		id          uuid PRIMARY KEY,
		remote_addr text,
		rowid       int
	);`,

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.queries (
  	id uuid PRIMARY KEY,
  	text text,
  	connection_id uuid REFERENCES treesql_live_queries.connections
	);`,
}

func CreateSchema(db *sql.DB) error {
	for _, stmt := range schema {
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error running %s: %v", stmt, err)
		}
	}
	return nil
}
