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

func LiveQuery(conn *sql.DB, dbSchema schema2.Schema) (chan *Event, error) {
	res := conn.QueryRow("SELECT cluster_logical_timestamp()")
	var timestamp string
	if err := res.Scan(&timestamp); err != nil {
		return nil, fmt.Errorf("getting timestamp to start changefeeds: %v", err)
	}

	c := make(chan *Event)
	for tableName := range dbSchema {
		go func() {
			res, err := conn.Query(fmt.Sprintf("CREATE CHANGEFEED FOR TABLE %s WITH cursor = '%s'", tableName, timestamp))
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
	log.Printf("opened changefeeds for %d tables", len(dbSchema))
	return c, nil
}