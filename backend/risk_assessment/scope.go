package risk_assessment

import (
  "context"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"

  "github.com/starnight/riskassessment/backend/config"
)

type Scope struct {
  ID primitive.ObjectID `bson:"_id"`
  CreateTime time.Time
  Name string `binding:"required"`
}

type IScopeUtils interface {
  AddScope(scope *Scope) (primitive.ObjectID, error)
  GetScopes() ([]Scope, error)
  GetScopeByID(id primitive.ObjectID) (Scope, error)
  GetScopeByIDs(ids []primitive.ObjectID) ([]Scope, error)
  HasScopeID(id primitive.ObjectID) (bool, error)
  UpdateScope(scope *Scope) (error)
}

type ScopeUtils struct {
  DB_Client *mongo.Client
}

var SCOPE_MONGO_DB string = config.DB_NAME
const SCOPE_COLLECTION = "scopes"

func (utils *ScopeUtils) AddScope(scope *Scope) (primitive.ObjectID, error) {
  scope.ID = primitive.NewObjectID()
  scope.CreateTime = time.Now().UTC()

  coll := utils.DB_Client.Database(SCOPE_MONGO_DB).Collection(SCOPE_COLLECTION)
  res, err := coll.InsertOne(context.TODO(), scope)
  id := res.InsertedID.(primitive.ObjectID)
  return id, err
}

func (utils *ScopeUtils) GetScopes() ([]Scope, error) {
  var scopes []Scope

  coll := utils.DB_Client.Database(SCOPE_MONGO_DB).Collection(SCOPE_COLLECTION)
  filter := bson.D{{}}
  cur, err := coll.Find(context.TODO(), filter)
  if err != nil {
    return scopes, err
  }

  err = cur.All(context.TODO(), &scopes)
  return scopes, err
}

func (utils *ScopeUtils) GetScopeByID(id primitive.ObjectID) (Scope, error) {
  var scope Scope

  coll := utils.DB_Client.Database(SCOPE_MONGO_DB).Collection(SCOPE_COLLECTION)
  filter := bson.D{{ "_id", id }}
  err := coll.FindOne(context.TODO(), filter).Decode(&scope)
  return scope, err
}

func (utils *ScopeUtils) GetScopeByIDs(ids []primitive.ObjectID) ([]Scope, error) {
  var scopes []Scope

  coll := utils.DB_Client.Database(SCOPE_MONGO_DB).Collection(SCOPE_COLLECTION)
  filter := bson.M{"_id": bson.M{"$in": ids}}
  cur, err := coll.Find(context.TODO(), filter)
  if err != nil {
    return scopes, err
  }

  err = cur.All(context.TODO(), &scopes)
  return scopes, err
}

func (utils *ScopeUtils) HasScopeID(id primitive.ObjectID) (bool, error) {
  coll := utils.DB_Client.Database(SCOPE_MONGO_DB).Collection(SCOPE_COLLECTION)
  filter := bson.D{{"_id", id}}
  count, err := coll.CountDocuments(context.TODO(), filter)
  return count > 0, err
}

func (utils *ScopeUtils) UpdateScope(scope *Scope) (error) {
  filter := bson.D{{ "_id", scope.ID }}

  coll := utils.DB_Client.Database(SCOPE_MONGO_DB).Collection(SCOPE_COLLECTION)
  _, err := coll.ReplaceOne(context.TODO(), filter, scope)
  return err
}
