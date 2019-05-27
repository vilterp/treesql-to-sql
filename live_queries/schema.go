package live_queries

import (
	"database/sql"
	"fmt"
)

var schema = []string{
	`CREATE DATABASE IF NOT EXISTS treesql_live_queries;`,

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.connections (
		id          uuid PRIMARY KEY,
		remote_addr text
	);`,

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.queries (
  	id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  	text          text,
  	connection_id uuid REFERENCES treesql_live_queries.connections
	);`,

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.col_listeners (
		id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
		table_name  text,
		column_name text,
		value       text
	);`,

	// TODO(vilterp): what about full table listeners?

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.record_listeners (
		id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
		table_name  text,
		pk_value    text
	);`,

	`CREATE TABLE IF NOT EXISTS treesql_live_queries.query_listeners (
		id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    col_listener    uuid REFERENCES treesql_live_queries.col_listeners,
    record_listener uuid REFERENCES treesql_live_queries.record_listeners,
    query           uuid REFERENCES treesql_live_queries.queries,
    path            jsonb
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
