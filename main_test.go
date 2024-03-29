package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetArguments() {
	httpConfig.url = ""
	databaseConfig.host = ""
	databaseConfig.port = ""
	databaseConfig.user = ""
	databaseConfig.pass = ""
	databaseConfig.database = ""
	databaseConfig.collection = ""
}

func TestUrlIsRequired(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		if r != nil {
			assert.NotNil(t, r)
			assert.Contains(t, r, "-request.url is required")
		}
	}()
	os.Args = []string{"cmd"}
	resetArguments()

	checkArguments()
}

func TestDBDatabaseIsRequired(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		if r != nil {
			assert.NotNil(t, r)
			assert.Contains(t, r, "-db.databse is required")
		}
	}()
	os.Args = []string{"cmd", "-request.url=http://host:port/path?a=1&b=2"}
	resetArguments()

	checkArguments()
}

func TestDBCollectionIsRequired(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		if r != nil {
			assert.NotNil(t, r)
			assert.Contains(t, r, "-db.collection is required")
		}
	}()
	os.Args = []string{"cmd", "-request.url=http://host:port/path?a=1&b=2", "-db.database=foo"}
	resetArguments()

	checkArguments()
}
