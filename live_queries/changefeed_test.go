package live_queries

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeEvent(t *testing.T) {
	rawEvent := &RawEvent{
		Table: "blog",
		Key:   "[8]",
		Value: `{"after": {"id": 8, "title": "bar", "body": "baz"}}`,
	}
	evt, err := decodeEvent(rawEvent)
	if err != nil {
		t.Fatal(err)
	}
	payload := map[string]interface{}{"body": "baz", "id": float64(8), "title": "bar"}
	require.Equal(t, payload, evt.Payload.After)
	require.Equal(t, []interface{}{float64(8)}, evt.Key)
}
