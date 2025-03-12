package drivers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlite(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	sdb := Sqlite("sqlite:./test.db")

	err := sdb.Open()
	require.NoError(err)
	defer sdb.Close()
	defer os.Remove("./test.db")

	err = sdb.Exec(`CREATE TABLE IF NOT Exists test (id INTEGER PRIMARY KEY,title text NOT NULL,body text)`, nil)

	require.NoError(err)

	err = sdb.Exec(`insert into test (title, body) values (?, ?)`, "test", "ing")

	assert.NoError(err)

	cols, data, err := sdb.Query(`select * from test`)
	assert.NoError(err)
	assert.Len(cols, 3)
	assert.Len(data, 1)
	rd := data[0].(map[string]interface{})
	assert.Equal("test", rd["title"])
	assert.Equal("ing", rd["body"])
	assert.Equal("1", rd["id"])
}
