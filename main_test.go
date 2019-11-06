package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetArguments() {
	urlToGather = ""
	hostMongo = ""
	portMongo = ""
	userMongo = ""
	passMongo = ""
	databseMongo = ""
	collectionMongo = ""
}

func TestUrlIsRequired(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		if r != nil {
			assert.NotNil(t, r)
			assert.Contains(t, r, "-url is required")
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
	os.Args = []string{"cmd", "-url=http://host:port/path?a=1&b=2"}
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
	os.Args = []string{"cmd", "-url=http://host:port/path?a=1&b=2", "-db.database=foo"}
	resetArguments()

	checkArguments()
}
