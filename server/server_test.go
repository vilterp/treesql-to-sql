package server

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vilterp/treesql-to-sql/util"
)

func TestServer(t *testing.T) {
	s, err := NewServer("user=root dbname=management_console_dev sslmode=disable port=26257")
	assert.NoError(t, err)

	ts := httptest.NewServer(util.Logger(s))

	body := strings.NewReader("SELECT * FROM clusters")
	resp, err := ts.Client().Post(ts.URL + "/sql", "application/x-sql", body)
	assert.NoError(t, err)

	fmt.Println("resp:", resp)

	assert.Equal(t, 200, resp.StatusCode)
	respBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	fmt.Println("respBytes:", respBytes)

	assert.Contains(t, string(respBytes), "Rows")
}
