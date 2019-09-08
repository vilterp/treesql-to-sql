package live_queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	schema2 "github.com/vilterp/treesql-to-sql/schema"
)

type RawEvent struct {
	Table string
	Key   string
	Value string
}

type Event struct {
	Table   string
	Key     []interface{}
	Payload EventPayload
}

type EventPayload struct {
	Before map[string]interface{}
	After  map[string]interface{}
}

func LiveQuery(conn *sql.DB, dbSchema schema2.Schema) (chan *Event, error) {
	res := conn.QueryRow("SELECT cluster_logical_timestamp()")
	var timestamp string
	if err := res.Scan(&timestamp); err != nil {
		return nil, fmt.Errorf("getting timestamp to start changefeeds: %v", err)
	}

	c := make(chan *Event)
	for tn := range dbSchema {
		tableName := tn // because Go closures and loops interact weirdly
		go func() {
			res, err := conn.Query(fmt.Sprintf("CREATE CHANGEFEED FOR TABLE %s WITH cursor = '%s'", tableName, timestamp))
			if err != nil {
				log.Printf("creating changefeed for table %s: %v", tableName, err)
			}
			for {
				res.Next()
				rawEvt := &RawEvent{}
				if err := res.Scan(&rawEvt.Table, &rawEvt.Key, &rawEvt.Value); err != nil {
					log.Println("err reading from changefeed:", err)
				}
				evt, err := decodeEvent(rawEvt)
				if err != nil {
					log.Println("error decoding event", rawEvt, ":", err)
				}
				c <- evt
			}
		}()
	}
	log.Printf("opened changefeeds for %d tables", len(dbSchema))
	return c, nil
}

func decodeEvent(event *RawEvent) (*Event, error) {
	ret := &Event{
		Table: event.Table,
	}
	if err := json.Unmarshal([]byte(event.Key), &ret.Key); err != nil {
		return nil, fmt.Errorf("decoding key: %v", err)
	}
	if err := json.Unmarshal([]byte(event.Value), &ret.Payload); err != nil {
		return nil, fmt.Errorf("decoding value: %v", err)
	}
	return ret, nil
}
