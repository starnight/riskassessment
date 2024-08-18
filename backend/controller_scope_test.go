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

  "github.com/starnight/risk_assessment/auth"
  "github.com/starnight/risk_assessment/risk_assessment"
)

type mockScopeUtils struct {
  mock.Mock
}

func (m *mockScopeUtils) AddScope(scope *risk_assessment.Scope) (primitive.ObjectID, error) {
  args := m.Called(scope)
  return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *mockScopeUtils) GetScopes() ([]risk_assessment.Scope, error) {
  args := m.Called()
  return args.Get(0).([]risk_assessment.Scope), args.Error(1)
}

func (m *mockScopeUtils) GetScopeByID(id primitive.ObjectID) (risk_assessment.Scope, error) {
  args := m.Called(id)
  return args.Get(0).(risk_assessment.Scope), args.Error(1)
}

func (m *mockScopeUtils) GetScopeByIDs(ids []primitive.ObjectID) ([]risk_assessment.Scope, error) {
  args := m.Called(ids)
  return args.Get(0).([]risk_assessment.Scope), args.Error(1)
}

func (m *mockScopeUtils) HasScopeID(id primitive.ObjectID) (bool, error) {
  args := m.Called(id)
  return args.Get(0).(bool), args.Error(1)
}

func (m *mockScopeUtils) UpdateScope(scope *risk_assessment.Scope) (error) {
  args := m.Called(scope)
  return args.Error(0)
}

func TestGetScopes(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScopes := []risk_assessment.Scope {
    {
      ID: primitive.NewObjectID(),
      CreateTime: time.Now().UTC(),
      Name: "Test scope1",
    },
    {
      ID: primitive.NewObjectID(),
      CreateTime: time.Now().UTC(),
      Name: "Test scope2",
    },
  }
  scope_util_mck.On("GetScopes", mock.Anything, mock.Anything).Return(mockScopes, nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("role", uint(auth.Administrator))
  session.Save()

  ap.GetScopes(c)

  var scope_page ScopePage
  json.Unmarshal(w.Body.Bytes(), &scope_page)
  assert.Equal(t, http.StatusOK, w.Code)
  var reduced_scopes []ReducedScope
  for _, scope := range mockScopes {
    reduced_scopes = append(reduced_scopes,
                            ReducedScope {
                              ID: scope.ID,
			      Name: scope.Name,
                            })
  }
  assert.Equal(t, reduced_scopes, scope_page.Scopes)
}

func TestGetScopesFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScopes := []risk_assessment.Scope{}

  scope_util_mck.On("GetScopes", mock.Anything, mock.Anything).Return(mockScopes, nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("role", uint(auth.NormalUser))
  session.Save()

  ap.GetScopes(c)

  assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetScopesFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScopes := []risk_assessment.Scope{}

  err := errors.New("Get failed")
  scope_util_mck.On("GetScopes", mock.Anything, mock.Anything).Return(mockScopes, err)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("role", uint(auth.Administrator))
  session.Save()

  ap.GetScopes(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetScopesByUser(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockUser := auth.User{
    Scopes: []primitive.ObjectID {
      primitive.NewObjectID(),
      primitive.NewObjectID(),
    },
  }
  auth_util_mck.On("GetUserByID", mock.Anything).Return(mockUser, nil)
  mockScopes := []risk_assessment.Scope {
    {
      ID: mockUser.Scopes[0],
      CreateTime: time.Now().UTC(),
      Name: "Test scope1",
    },
    {
      ID: mockUser.Scopes[1],
      CreateTime: time.Now().UTC(),
      Name: "Test scope2",
    },
  }
  scope_util_mck.On("GetScopeByIDs", mock.Anything).Return(mockScopes, nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", primitive.NewObjectID().Hex())
  session.Set("role", uint(auth.Administrator))
  session.Save()

  ap.GetScopesByUser(c)

  var scope_page ScopePage
  json.Unmarshal(w.Body.Bytes(), &scope_page)
  assert.Equal(t, http.StatusOK, w.Code)
  var reduced_scopes []ReducedScope
  for _, scope := range mockScopes {
    reduced_scopes = append(reduced_scopes,
                            ReducedScope {
                              ID: scope.ID,
			      Name: scope.Name,
                            })
  }
  assert.Equal(t, reduced_scopes, scope_page.Scopes)
}

func TestGetScopesByUserFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockUser := auth.User{
    Scopes: []primitive.ObjectID {
      primitive.NewObjectID(),
      primitive.NewObjectID(),
    },
  }
  err := errors.New("Get failed")
  auth_util_mck.On("GetUserByID", mock.Anything).Return(mockUser, err)
  mockScopes := []risk_assessment.Scope {
    {
      ID: mockUser.Scopes[0],
      CreateTime: time.Now().UTC(),
      Name: "Test scope1",
    },
    {
      ID: mockUser.Scopes[1],
      CreateTime: time.Now().UTC(),
      Name: "Test scope2",
    },
  }
  scope_util_mck.On("GetScopeByIDs", mock.Anything).Return(mockScopes, nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", primitive.NewObjectID().Hex())
  session.Set("role", uint(auth.Administrator))
  session.Save()

  ap.GetScopesByUser(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetScopesByUserFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockUser := auth.User{
    Scopes: []primitive.ObjectID {
      primitive.NewObjectID(),
      primitive.NewObjectID(),
    },
  }
  auth_util_mck.On("GetUserByID", mock.Anything).Return(mockUser, nil)
  mockScopes := []risk_assessment.Scope {
    {
      ID: mockUser.Scopes[0],
      CreateTime: time.Now().UTC(),
      Name: "Test scope1",
    },
    {
      ID: mockUser.Scopes[1],
      CreateTime: time.Now().UTC(),
      Name: "Test scope2",
    },
  }
  err := errors.New("Get failed")
  scope_util_mck.On("GetScopeByIDs", mock.Anything).Return(mockScopes, err)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, session := GetMockContext(req)

  session.Set("id", primitive.NewObjectID().Hex())
  session.Set("role", uint(auth.Administrator))
  session.Save()

  ap.GetScopesByUser(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddScope(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScopeID := primitive.NewObjectID()
  scope_util_mck.On("AddScope", mock.Anything).Return(mockScopeID, nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {
    Name: "Test Scope",
  }
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.AddScope(c)

  assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddScopeFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScopeID := primitive.NewObjectID()
  scope_util_mck.On("AddScope", mock.Anything).Return(mockScopeID, nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {}
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.AddScope(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddScopeFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScopeID := primitive.NewObjectID()
  scope_util_mck.On("AddScope", mock.Anything).Return(mockScopeID, nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {
    Name: " ",
  }
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.AddScope(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddScopeFailed3(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScopeID := primitive.NewObjectID()
  err := errors.New("Add failed")
  scope_util_mck.On("AddScope", mock.Anything).Return(mockScopeID, err)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {
    Name: "Test Scope",
  }
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.AddScope(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateScope(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScope := risk_assessment.Scope {
    ID: primitive.NewObjectID(),
    CreateTime: time.Now().UTC(),
    Name: "Original Scope",
  }
  scope_util_mck.On("GetScopeByID", mock.Anything).Return(mockScope, nil)
  scope_util_mck.On("UpdateScope", mock.Anything).Return(nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {
    ID: mockScope.ID,
    CreateTime: time.Now().UTC(),
    Name: "Test Scope",
  }
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.UpdateScope(c)

  assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateScopeFailed1(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScope := risk_assessment.Scope {
    ID: primitive.NewObjectID(),
    CreateTime: time.Now().UTC(),
    Name: "Original Scope",
  }
  scope_util_mck.On("GetScopeByID", mock.Anything).Return(mockScope, nil)
  scope_util_mck.On("UpdateScope", mock.Anything).Return(nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {}
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.UpdateScope(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateScopeFailed2(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScope := risk_assessment.Scope {
    ID: primitive.NewObjectID(),
    CreateTime: time.Now().UTC(),
    Name: "Original Scope",
  }
  scope_util_mck.On("GetScopeByID", mock.Anything).Return(mockScope, nil)
  scope_util_mck.On("UpdateScope", mock.Anything).Return(nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {
    ID: mockScope.ID,
    CreateTime: time.Now().UTC(),
    Name: " ",
  }
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.UpdateScope(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateScopeFailed3(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScope := risk_assessment.Scope {
    ID: primitive.NewObjectID(),
    CreateTime: time.Now().UTC(),
    Name: "Original Scope",
  }
  err := errors.New("Get failed")
  scope_util_mck.On("GetScopeByID", mock.Anything).Return(mockScope, err)
  scope_util_mck.On("UpdateScope", mock.Anything).Return(nil)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {
    ID: mockScope.ID,
    CreateTime: time.Now().UTC(),
    Name: "Test Scope",
  }
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.UpdateScope(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateScopeFailed4(t *testing.T) {
  auth_util_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  scope_util_mck := new(mockScopeUtils)
  mockScope := risk_assessment.Scope {
    ID: primitive.NewObjectID(),
    CreateTime: time.Now().UTC(),
    Name: "Original Scope",
  }
  scope_util_mck.On("GetScopeByID", mock.Anything).Return(mockScope, nil)
  err := errors.New("Update failed")
  scope_util_mck.On("UpdateScope", mock.Anything).Return(err)
  ap := ScopesApp{User_utils: auth_util_mck, Csrf_utils: csrf_util_mck, Scope_utils: scope_util_mck}

  testScope := risk_assessment.Scope {
    ID: mockScope.ID,
    CreateTime: time.Now().UTC(),
    Name: "Test Scope",
  }
  json_bytes, _ := json.Marshal(testScope)
  json_buf := bytes.NewBuffer(json_bytes)

  gin.SetMode(gin.TestMode)
  req := httptest.NewRequest("POST", "/", json_buf)
  req.Header.Add("Content-Type", binding.MIMEJSON)
  c, w, _ := GetMockContext(req)

  ap.UpdateScope(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}
