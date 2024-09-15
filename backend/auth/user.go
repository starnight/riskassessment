package auth

import (
  "context"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"

  "github.com/starnight/riskassessment/backend/config"
)

const (
  NormalUser = 0
  Administrator = 1
)

type UserInfo struct {
  Role uint
}

type User struct {
  ID primitive.ObjectID `bson:"_id"`
  CreateTime time.Time
  Account string `binding:"required"`
  Password string `binding:"required"`
  Role uint
  Scopes []primitive.ObjectID
}

type IUserUtils interface {
  AddUser(user *User) (primitive.ObjectID, error)
  GetUserByID(id primitive.ObjectID) (User, error)
  GetUserByAccount(account string) (User, error)
  GetUserByAccountPwd(account string, password string) (User, error)
  HasUser() (bool, error)
  UserHasScopeID(u_id primitive.ObjectID, s_id primitive.ObjectID) (bool, error)
  UpdateUser (user *User) (error)
}

type UserUtils struct {
  DB_Client *mongo.Client
}

var USER_MONGO_DB string = config.DB_NAME
var USER_COLLECTION string = "users"

func (utils *UserUtils) AddUser(user *User) (primitive.ObjectID, error) {
  var id primitive.ObjectID

  user.ID = primitive.NewObjectID()
  user.CreateTime = time.Now().UTC()

  if user.Scopes == nil {
    user.Scopes = []primitive.ObjectID{}
  }

  coll := utils.DB_Client.Database(USER_MONGO_DB).Collection(USER_COLLECTION)
  res, err := coll.InsertOne(context.TODO(), user)
  id = res.InsertedID.(primitive.ObjectID)
  return id, err
}

func (utils *UserUtils) GetUserByID(id primitive.ObjectID) (User, error) {
  var user User

  coll := utils.DB_Client.Database(USER_MONGO_DB).Collection(USER_COLLECTION)
  filter := bson.D{{ "_id", id }}
  err := coll.FindOne(context.TODO(), filter).Decode(&user)
  return user, err
}

func (utils *UserUtils) GetUserByAccount(account string) (User, error) {
  var user User

  coll := utils.DB_Client.Database(USER_MONGO_DB).Collection(USER_COLLECTION)
  filter := bson.D{{"account", account}}
  err := coll.FindOne(context.TODO(), filter).Decode(&user)
  return user, err
}

func (utils *UserUtils) GetUserByAccountPwd(account string, password string) (User, error) {
  var user User

  coll := utils.DB_Client.Database(USER_MONGO_DB).Collection(USER_COLLECTION)
  filter := bson.D{{"account", account}, {"password", password}}
  err := coll.FindOne(context.TODO(), filter).Decode(&user)
  return user, err
}

func (utils *UserUtils) HasUser() (bool, error) {
  coll := utils.DB_Client.Database(USER_MONGO_DB).Collection(USER_COLLECTION)
  filter := bson.D{{}}
  count, err := coll.CountDocuments(context.TODO(), filter)
  return count > 0, err
}

func (utils *UserUtils) UserHasScopeID(u_id primitive.ObjectID, s_id primitive.ObjectID) (bool, error) {
  coll := utils.DB_Client.Database(USER_MONGO_DB).Collection(USER_COLLECTION)
  filter := bson.M{ "_id": u_id, "scopes": bson.M{ "$elemMatch": bson.M{ "$eq": s_id }}}
  count, err := coll.CountDocuments(context.TODO(), filter)
  return count > 0, err
}

func (utils *UserUtils) UpdateUser(user *User) (error) {
  filter := bson.D{{ "_id", user.ID }}

  coll := utils.DB_Client.Database(USER_MONGO_DB).Collection(USER_COLLECTION)
  _, err := coll.ReplaceOne(context.TODO(), filter, user)
  return err
}
