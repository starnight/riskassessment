package risk_assessment

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/starnight/riskassessment/backend/database"

  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

var asset_utils = AssetUtils{
  DB_Client: database.ConnectDB(database.GetDBStr("")),
}

func TestAddAsset(t *testing.T) {
  expected_blg_ctg := "big"
  expected_small_ctg := "small"
  expected_name := "test name"

  asset := Asset{
    Scope: primitive.NewObjectID(),
    BigCategory: expected_blg_ctg,
    SmallCategory: expected_small_ctg,
    Name: expected_name,
  }

  err1 := asset_utils.AddAsset(&asset)
  assert.Nil(t, err1)

  assets, err2 := asset_utils.GetAssets(0, 0)
  assert.Nil(t, err2)
  asset2 := assets[len(assets)-1]
  assert.NotNil(t, asset2.ID)
  assert.NotNil(t, asset2.CreateTime)
  assert.Equal(t, expected_blg_ctg, asset2.BigCategory)
  assert.Equal(t, expected_small_ctg, asset2.SmallCategory)
  assert.Equal(t, expected_name, asset2.Name)
  /* Have not set the asset's Value. So, values should be zeros. */
  assert.Zero(t, asset2.Value.Confidentiality)
  assert.Zero(t, asset2.Value.Integrity)
  assert.Zero(t, asset2.Value.Availability)
  assert.Zero(t, len(asset2.Risks))

  asset3, err3 := asset_utils.GetAssetByID(asset2.ID)
  assert.Nil(t, err3)
  assert.Equal(t, asset2, asset3)

  assets4, err4 := asset_utils.GetAssetsByScopeID(asset2.Scope)
  assert.Nil(t, err4)
  assert.Equal(t, assets, assets4)
}

func TestSetAssetValue(t *testing.T) {
  expected_c := uint(4)
  expected_i := uint(3)
  expected_a := uint(2)

  assets, _ := asset_utils.GetAssets(0, 0)
  orig_asset := assets[len(assets)-1]

  err1 := asset_utils.SetAssetValue(orig_asset.ID.Hex(), expected_c, expected_i, expected_a)
  assert.Nil(t, err1)

  assets2, err2 := asset_utils.GetAssets(0, 0)
  assert.Nil(t, err2)
  new_asset := assets2[len(assets2)-1]
  assert.Equal(t, orig_asset.ID, new_asset.ID)
  assert.Equal(t, expected_c, new_asset.Value.Confidentiality)
  assert.Equal(t, expected_i, new_asset.Value.Integrity)
  assert.Equal(t, expected_a, new_asset.Value.Availability)
}

func TestUpdteAsset(t *testing.T) {
  expected_c := uint(2)
  expected_i := uint(4)
  expected_a := uint(3)

  assets, _ := asset_utils.GetAssets(0, 0)
  orig_asset := assets[len(assets)-1]

  orig_asset.Value.Confidentiality = expected_c
  orig_asset.Value.Integrity = expected_i
  orig_asset.Value.Availability = expected_a
  risk := Risk{
    Threat: "Test threat",
    Vulnerability: "Test vulnerability",
    CurrentControl: "Test control",
    Possibility: uint(2),
    Impact: uint(3),
  }
  orig_asset.Risks = append(orig_asset.Risks, risk)

  err1 := asset_utils.UpdateAsset(&orig_asset)
  assert.Nil(t, err1)

  saved_asset, err2 := asset_utils.GetAssetByID(orig_asset.ID)
  assert.Nil(t, err2)
  assert.Equal(t, orig_asset, saved_asset)
}

func TestDeleteAsset(t *testing.T) {
  assets, _ := asset_utils.GetAssets(0, 1)
  assert.Equal(t, 1, len(assets))
  delete_asset := assets[0]

  err1 := asset_utils.DeleteAsset(delete_asset.ID)
  assert.Nil(t, err1)

  _, err2 := asset_utils.GetAssetByID(delete_asset.ID)
  assert.Equal(t, mongo.ErrNoDocuments, err2)
}
