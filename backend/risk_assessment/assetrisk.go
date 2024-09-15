/*
 * https://files.chcg.gov.tw/files/13_20220419130506915_2.%E8%B3%87%E7%94%A2%E7%9B%A4%E9%BB%9E%E8%88%87%E9%A2%A8%E9%9A%AA%E8%A9%95%E9%91%91%E8%A8%93%E7%B7%B4_1110425.pdf
 */

package risk_assessment

import (
  "context"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  "github.com/starnight/riskassessment/backend/config"
)

type Asset struct {
  ID primitive.ObjectID `bson:"_id"`
  CreateTime time.Time
  Scope primitive.ObjectID
  BigCategory string
  SmallCategory string
  Name string
  Owner string
  Value Value
  Risks []Risk
}

type Value struct {
  Confidentiality uint `binding:"required"`
  Integrity uint `binding:"required"`
  Availability uint `binding:"required"`
}

type Risk struct {
  Threat string
  Vulnerability string
  CurrentControl string
  Possibility uint `binding:"required"`
  Impact uint `binding:"required"`
}

type IAssetUtils interface {
  AddAsset(asset *Asset) (error)
  GetAssetByID(id primitive.ObjectID) (Asset, error)
  GetAssetsByScopeID(id primitive.ObjectID) ([]Asset, error)
  GetAssets(offset int64, amount int64) ([]Asset, error)
  SetAssetValue(id string, c uint, i uint, a uint) (error)
  UpdateAsset(asset *Asset) (error)
  DeleteAsset(id primitive.ObjectID) (error)
}

type AssetUtils struct {
  DB_Client *mongo.Client
}

var _MONGO_DB string = config.DB_NAME
var _COLLECTION string = "assets"

func (utils *AssetUtils) AddAsset(asset *Asset) (error) {
  asset.ID = primitive.NewObjectID()
  asset.CreateTime = time.Now().UTC()
  if asset.Risks == nil {
    asset.Risks = []Risk{}
  }

  coll := utils.DB_Client.Database(_MONGO_DB).Collection(_COLLECTION)
  _, err := coll.InsertOne(context.TODO(), asset)
  return err
}

func (utils *AssetUtils) GetAssetByID(id primitive.ObjectID) (Asset, error) {
  var asset Asset

  coll := utils.DB_Client.Database(_MONGO_DB).Collection(_COLLECTION)
  filter := bson.D{{ "_id", id }}
  err := coll.FindOne(context.TODO(), filter).Decode(&asset)
  return asset, err
}

func (utils *AssetUtils) GetAssetsByScopeID(id primitive.ObjectID) ([]Asset, error) {
  var assets []Asset

  coll := utils.DB_Client.Database(_MONGO_DB).Collection(_COLLECTION)
  filter := bson.D{{ "scope", id }}
  cur, err := coll.Find(context.TODO(), filter)
  if err != nil {
    return assets, err
  }

  err = cur.All(context.TODO(), &assets)
  return assets, err
}

func (utils *AssetUtils) GetAssets(offset int64, amount int64) ([]Asset, error) {
  var assets []Asset

  coll := utils.DB_Client.Database(_MONGO_DB).Collection(_COLLECTION)
  filter := bson.D{{}}
  opts := options.Find().SetLimit(amount).SetSkip(offset)
  cur, err := coll.Find(context.TODO(), filter, opts)
  if err != nil {
    return assets, err
  }

  err = cur.All(context.TODO(), &assets)
  return assets, err
}

func (utils *AssetUtils) SetAssetValue(id string, c uint, i uint, a uint) (error) {
  a_id, _ := primitive.ObjectIDFromHex(id)
  filter := bson.M{ "_id": a_id }
  update := bson.M{
    "$set": bson.M{
      "value": bson.M{
        "confidentiality": c,
        "integrity": i,
        "availability": a,
      },
    },
  }

  coll := utils.DB_Client.Database(_MONGO_DB).Collection(_COLLECTION)
  _, err := coll.UpdateOne(context.TODO(), filter, update)
  return err
}

func (utils *AssetUtils) UpdateAsset(asset *Asset) (error) {
  filter := bson.D{{ "_id", asset.ID }}

  coll := utils.DB_Client.Database(_MONGO_DB).Collection(_COLLECTION)
  _, err := coll.ReplaceOne(context.TODO(), filter, asset)
  return err
}

func (utils *AssetUtils) DeleteAsset(id primitive.ObjectID) (error) {
  filter := bson.D{{ "_id", id }}

  coll := utils.DB_Client.Database(_MONGO_DB).Collection(_COLLECTION)
  _, err := coll.DeleteOne(context.TODO(), filter)
  return err
}
