package snowflake

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	sf "github.com/snowflakedb/gosnowflake"
	"github.com/vilterp/treesql-to-sql/live_queries"
)

func NewSnowflakeListener(cfg *sf.Config, tableName string) (*live_queries.Listener, error) {
	dsn, err := sf.DSN(cfg)
	if err != nil {
		return nil, err
	}
	conn, err := sql.Open("snowflake", dsn)
	if err != nil {
		return nil, err
	}

	_, err = conn.Query("SELECT 1")
	if err != nil {
		return nil, err
	}

	log.Println("connected to Snowflake")

	return &live_queries.Listener{
		Insert: func(r live_queries.Row) error {
			rowJSON, err := json.Marshal(r)
			if err != nil {
				return fmt.Errorf("error marshalling: %v", err)
			}
			// TODO: use placeholders...
			insStatement := fmt.Sprintf("INSERT INTO %s SELECT parse_json('%s')", tableName, rowJSON)
			log.Println("inserting into snowflake", insStatement)
			_, err = conn.Exec(insStatement)
			return err
		},
	}, nil
}
