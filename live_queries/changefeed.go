package live_queries

import (
	"database/sql"
	"fmt"
	"log"

	schema2 "github.com/vilterp/treesql-to-sql/schema"
)

type Event struct {
	Table string
	Key   string
	Value string
}

func LiveQuery(conn *sql.DB, schema schema2.Schema) (chan *Event, error) {
	c := make(chan *Event)
	for tableName := range schema {
		go func() {
			res, err := conn.Query(fmt.Sprintf("CREATE CHANGEFEED FOR TABLE %s", tableName))
			if err != nil {
				log.Printf("creating changefeed for table %s: %v", tableName, err)
			}
			for {
				res.Next()
				evt := &Event{}
				if err := res.Scan(&evt.Table, &evt.Key, &evt.Value); err != nil {
					log.Println("err reading from changefeed:", err)
				}
				c <- evt
			}
		}()
	}
	return c, nil
}
