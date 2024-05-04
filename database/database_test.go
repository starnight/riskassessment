package database

import (
  "os"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestGetDBStr(t *testing.T) {
  expect_uri := "test uri"
  uri := GetDBStr(expect_uri)
  assert.Equal(t, expect_uri, uri)

  expect_uri = "test ENV uri"
  os.Setenv("MONGODB_URI", expect_uri)
  uri = GetDBStr("")
  assert.Equal(t, expect_uri, uri)

  expect_uri = "mongodb://localhost:27017"
  os.Setenv("MONGODB_URI", "")
  uri = GetDBStr("")
  assert.Equal(t, expect_uri, uri)
}

func TestConnectDB(t *testing.T) {
  mongo_uri := "mongodb://localhost:27017"
  client := ConnectDB(mongo_uri)
  assert.NotNil(t, client)

  expect_client := client
  client = ConnectDB(mongo_uri)
  assert.Equal(t, expect_client, client)
}
