package live_queries

import (
	"database/sql"
	"fmt"

	"github.com/vilterp/treesql-to-sql/parse"
)

type QueryContext struct {
	conn      *sql.DB
	connID    string
	query     *parse.Select
	queryText string
}

func InsertListeners(queryCtx *QueryContext) error {
	row := queryCtx.conn.QueryRow(
		"INSERT INTO queries (text, connection_id) VALUES ($1, $2) RETURNING (id)",
		queryCtx.queryText, queryCtx.connID,
	)
	var queryID string
	err := row.Scan(&queryID)
	if err != nil {
		return fmt.Errorf("inserting query: %v", err)
	}

	// TODO(vilterp): walk result set, putting these in
	row = queryCtx.conn.QueryRow(
		"INSERT INTO col_listeners (table_name, column_name, value) VALUES ($1, $2, $3) RETURNING (id);",
		"comments", "post_id", "5",
	)
	var colListenerID string
	err = row.Scan(&colListenerID)
	if err != nil {
		return fmt.Errorf("inserting col listener: %v", err)
	}

	// insert query listener
	if _, err := queryCtx.conn.Exec(
		"INSERT INTO query_listeners (col_listener, query, path) VALUES ($1, $2, $3);",
		colListenerID, queryID, `["some", "path"]`,
	); err != nil {
		return fmt.Errorf("inserting query listener: %v", err)
	}

	return nil
}
