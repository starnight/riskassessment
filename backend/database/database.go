package database

import (
  "context"
  "os"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var _client *mongo.Client = nil

func GetDBStr(uri string) string {
  if uri != "" {
     return uri
  }
  uri = os.Getenv("MONGODB_URI")
  if uri == "" {
    uri = "mongodb://localhost:27017"
  }

  return uri
}

func ConnectDB(uri string) *mongo.Client {
  if (_client != nil) {
    return _client
  }

  var err error
  _client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
  if(err != nil) {
    panic(err)
  }

  return _client
}
