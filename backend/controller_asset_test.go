package main

import (
  "bytes"
  "errors"
  "testing"
  "time"
  "encoding/json"
  "net/http"
  "net/http/httptest"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/starnight/riskassessment/backend/risk_assessment"
)

type mockCsrtUtils struct {}

func (ap *mockCsrtUtils) AddCSRFToken(c *gin.Context) {
  return
}

type mockAssetUtils struct {
  mock.Mock
}

func (m *mockAssetUtils) AddAsset(asset *risk_assessment.Asset) (error) {
  args := m.Called(asset)
  return args.Error(0)
}

func (m *mockAssetUtils) GetAssetByID(id primitive.ObjectID) (risk_assessment.Asset, error) {
  args := m.Called(id)
  return args.Get(0).(risk_assessment.Asset), args.Error(1)
}

func (m *mockAssetUtils) GetAssetsByScopeID(id primitive.ObjectID) ([]risk_assessment.Asset, error) {
  args := m.Called(id)
  return args.Get(0).([]risk_assessment.Asset), args.Error(1)
}

func (m *mockAssetUtils) GetAssets(offset int64, amount int64) ([]risk_assessment.Asset, error) {
  args := m.Called(offset, amount)
  return args.Get(0).([]risk_assessment.Asset), args.Error(1)
}

func (m *mockAssetUtils) SetAssetValue(id string, c uint, i uint, a uint) (error) {
  args := m.Called(id, c, i, a)
  return args.Error(0)
}

func (m *mockAssetUtils) UpdateAsset(asset *risk_assessment.Asset) (error) {
  args := m.Called(asset)
  return args.Error(0)
}

func (m *mockAssetUtils) DeleteAsset(id primitive.ObjectID) (error) {
  args := m.Called(id)
  return args.Error(0)
}

func TestGetAssets(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  ast_util_mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAssets := []risk_assessment.Asset {
    {
      ID: primitive.NewObjectID(),
      CreateTime: time.Now().UTC(),
      BigCategory: "Big Category",
      SmallCategory: "Small Category",
      Name: "Test asset1",
    },
    {
      ID: primitive.NewObjectID(),
      CreateTime: time.Now().UTC(),
      BigCategory: "Big Category",
      SmallCategory: "Small Category",
      Name: "Test asset2",
    },
  }
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  ast_util_mck.On("GetAssetsByScopeID", mock.Anything).Return(mockAssets, nil)
  ap := AssetsApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Asset_utils: ast_util_mck}

  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Set("role", uint(0))
  session.Save()

  c.Params = append(c.Params, gin.Param{Key: "scopeID", Value: scopeID.Hex()})

  ap.GetAssets(c)

  var asset_page AssetPage
  json.Unmarshal(w.Body.Bytes(), &asset_page)
  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, mockAssets, asset_page.Assets)
}

func TestGetAssetsAuthorizeFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  ast_util_mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  err := errors.New("Get failed")
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, err)
  ap := AssetsApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Asset_utils: ast_util_mck}

  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Set("role", uint(0))
  session.Save()

  c.Params = append(c.Params, gin.Param{Key: "scopeID", Value: scopeID.Hex()})

  ap.GetAssets(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAssetsAuthorizeFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  ast_util_mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(false, nil)
  ap := AssetsApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Asset_utils: ast_util_mck}
  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Set("role", uint(0))
  session.Save()

  c.Params = append(c.Params, gin.Param{Key: "scopeID", Value: scopeID.Hex()})

  ap.GetAssets(c)

  assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetAssetsFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  ast_util_mck := new(mockAssetUtils)
  mockAssets := []risk_assessment.Asset {}
  err := errors.New("Get failed")
  ast_util_mck.On("GetAssetsByScopeID", mock.Anything).Return(mockAssets, err)
  ap := AssetsApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Asset_utils: ast_util_mck}

  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", primitive.NewObjectID().Hex())
  session.Set("role", uint(0))
  session.Save()

  ap.GetAssets(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
  assert.Nil(t, w.Body.Bytes())
}

func TestGetAssetsFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  ast_util_mck := new(mockAssetUtils)
  mockAssets := []risk_assessment.Asset {}
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  err := errors.New("Get failed")
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  ast_util_mck.On("GetAssetsByScopeID", mock.Anything).Return(mockAssets, err)
  ap := AssetsApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Asset_utils: ast_util_mck}

  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Set("role", uint(0))
  session.Save()

  c.Params = append(c.Params, gin.Param{Key: "scopeID", Value: scopeID.Hex()})

  ap.GetAssets(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
  assert.Nil(t, w.Body.Bytes())
}

func TestAddAsset(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  testAsset := risk_assessment.Asset {
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Test name",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }
  mck.On("AddAsset", &testAsset).Return(nil)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.AddAsset(c)

  assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddAssetFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  userID := primitive.NewObjectID()
  mck := new(mockAssetUtils)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  req := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.AddAsset(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAssetAuthorizeFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mck := new(mockAssetUtils)
  err := errors.New("Get failed")
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  testAsset := risk_assessment.Asset {
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Not Found",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }

  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)
  req := httptest.NewRequest("POST", "/", json_buf)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.AddAsset(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddAssetAuthorizeFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mck := new(mockAssetUtils)
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(false, nil)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  testAsset := risk_assessment.Asset {
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Not Found",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }

  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)
  req := httptest.NewRequest("POST", "/", json_buf)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.AddAsset(c)

  assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAddAssetFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mck := new(mockAssetUtils)
  err := errors.New("Add failed")
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  testAsset := risk_assessment.Asset {
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Add failed",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }
  mck.On("AddAsset", &testAsset).Return(err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.AddAsset(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateAsset(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {
    ID: primitive.NewObjectID(),
    CreateTime: time.Now().UTC(),
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Test asset1",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }
  mck.On("GetAssetByID", mockAsset.ID).Return(mockAsset, nil)
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  mck.On("UpdateAsset", &mockAsset).Return(nil)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  testAsset := mockAsset
  testAsset.CreateTime = time.Now()
  assert.NotEqual(t, mockAsset.CreateTime, testAsset.CreateTime)
  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.UpdateAsset(c)

  assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateAssetBadReq(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  req := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", primitive.NewObjectID().Hex())
  session.Save()

  ap.UpdateAsset(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAssetNotFound(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  mockAsset := risk_assessment.Asset {}
  assetID := primitive.NewObjectID()
  err := errors.New("Get failed")
  mck.On("GetAssetByID", assetID).Return(mockAsset, err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  testAsset := risk_assessment.Asset {
    ID: assetID,
    CreateTime: time.Now().UTC(),
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Not Found",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }
  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", primitive.NewObjectID().Hex())
  session.Save()

  ap.UpdateAsset(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAssetAuthorizeFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {
    ID: assetID,
    CreateTime: time.Now().UTC(),
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Update failed",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }
  mck.On("GetAssetByID", assetID).Return(mockAsset, nil)
  err := errors.New("Get failed")
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  testAsset := mockAsset
  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.UpdateAsset(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateAssetAuthorizeFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {
    ID: assetID,
    CreateTime: time.Now().UTC(),
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Update failed",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }
  mck.On("GetAssetByID", assetID).Return(mockAsset, nil)
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(false, nil)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  testAsset := mockAsset
  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.UpdateAsset(c)

  assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateAssetFailed(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {
    ID: assetID,
    CreateTime: time.Now().UTC(),
    Scope: scopeID,
    BigCategory: "Big Category",
    SmallCategory: "Small Category",
    Name: "Update failed",
    Owner: "Owner name",
    Value: risk_assessment.Value{
      Confidentiality: 4,
      Integrity: 3,
      Availability: 2,
    },
  }
  mck.On("GetAssetByID", assetID).Return(mockAsset, nil)
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  err := errors.New("Update failed")
  mck.On("UpdateAsset", &mockAsset).Return(err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  testAsset := mockAsset
  json_bytes, _ := json.Marshal(testAsset)
  json_buf := bytes.NewBuffer(json_bytes)

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.UpdateAsset(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteAsset(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {Scope: scopeID}
  mck.On("GetAssetByID", assetID).Return(mockAsset, nil)
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  mck.On("DeleteAsset", assetID).Return(nil)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  json_buf := bytes.NewBufferString("{\"id\":\"" + assetID.Hex() + "\"}")

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.DeleteAsset(c)

  assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteAssetBadReq(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  req := httptest.NewRequest("POST", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.DeleteAsset(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteAssetNotFound(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {Scope: scopeID}
  err := errors.New("Not Found")
  mck.On("GetAssetByID", assetID).Return(mockAsset, err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  json_buf := bytes.NewBufferString("{\"id\":\"" + assetID.Hex() + "\"}")

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.DeleteAsset(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteAuthorizeFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {Scope: scopeID}
  mck.On("GetAssetByID", assetID).Return(mockAsset, nil)
  err := errors.New("Get failed")
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  json_buf := bytes.NewBufferString("{\"id\":\"" + assetID.Hex() + "\"}")

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.DeleteAsset(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteAuthorizeFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {Scope: scopeID}
  mck.On("GetAssetByID", assetID).Return(mockAsset, nil)
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(false, nil)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  json_buf := bytes.NewBufferString("{\"id\":\"" + assetID.Hex() + "\"}")

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.DeleteAsset(c)

  assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteAssetFailed(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  mck := new(mockAssetUtils)
  userID := primitive.NewObjectID()
  assetID := primitive.NewObjectID()
  scopeID := primitive.NewObjectID()
  mockAsset := risk_assessment.Asset {Scope: scopeID}
  mck.On("GetAssetByID", assetID).Return(mockAsset, nil)
  auth_util_mck.On("UserHasScopeID", userID, scopeID).Return(true, nil)
  err := errors.New("Not Found")
  mck.On("DeleteAsset", assetID).Return(err)
  ap := AssetsApp{User_utils: auth_util_mck, Asset_utils: mck}

  json_buf := bytes.NewBufferString("{\"id\":\"" + assetID.Hex() + "\"}")

  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, session := GetMockContext(req)

  session.Set("id", userID.Hex())
  session.Save()

  ap.DeleteAsset(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}
