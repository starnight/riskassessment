package main

import (
  "bytes"
  "encoding/json"
  "errors"
  "net/http"
  "net/http/httptest"
  "net/url"
  "strings"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/starnight/risk_assessment/auth"
)

type mockUserUtils struct {
  mock.Mock
}

func (m *mockUserUtils) AddUser(user *auth.User) (primitive.ObjectID, error) {
  args := m.Called(user)
  return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *mockUserUtils) GetUserByID(id primitive.ObjectID) (auth.User, error) {
  args := m.Called(id)
  return args.Get(0).(auth.User), args.Error(1)
}

func (m *mockUserUtils) GetUserByAccount(account string) (auth.User, error) {
  args := m.Called(account)
  return args.Get(0).(auth.User), args.Error(1)
}

func (m *mockUserUtils) GetUserByAccountPwd(account string, password string) (auth.User, error) {
  args := m.Called(account, password)
  return args.Get(0).(auth.User), args.Error(1)
}

func (m *mockUserUtils) HasUser() (bool, error){
  args := m.Called()
  return args.Get(0).(bool), args.Error(1)
}

func (m *mockUserUtils) UserHasScopeID(u_id primitive.ObjectID, s_id primitive.ObjectID) (bool, error) {
  args := m.Called(u_id, s_id)
  return args.Get(0).(bool), args.Error(1)
}

func (m *mockUserUtils) UpdateUser(user *auth.User) (error){
  args := m.Called(user)
  return args.Error(0)
}

func TestGetLogin(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  csrf_util_mck := new(mockCsrtUtils)
  ap := AuthApp{User_utils: auth_utils_mck, Csrf_utils: csrf_util_mck}

  req := httptest.NewRequest("Get", "/", bytes.NewBufferString(""))
  c, w, _ := GetMockContext(req)

  ap.GetLogin(c)

  assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoLoginNoPwd(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  ap := AuthApp{User_utils: auth_utils_mck}

  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "")
  data.Set("_csrf", "csrf=")

  req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  c, w, _ := GetMockContext(req)

  ap.DoLogin(c)

  assert.Equal(t, http.StatusForbidden, w.Code)
  assert.Equal(t, "Wrong account or password", w.Body.String())
}

func TestDoLoginWronPwd(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User{}
  err := errors.New("Get failed")
  auth_utils_mck.On("GetUserByAccountPwd", mock.Anything, mock.Anything).Return(mockUser, err)
  ap := AuthApp{User_utils: auth_utils_mck}

  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  data.Set("_csrf", "csrf=")

  req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  c, w, _ := GetMockContext(req)

  ap.DoLogin(c)

  assert.Equal(t, http.StatusForbidden, w.Code)
  assert.Equal(t, "Wrong account or password", w.Body.String())
}

func TestDoLogin(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User{
    ID: primitive.NewObjectID(),
    Role: auth.NormalUser,
  }
  auth_utils_mck.On("GetUserByAccountPwd", mock.Anything, mock.Anything).Return(mockUser, nil)
  ap := AuthApp{User_utils: auth_utils_mck}

  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  data.Set("_csrf", "csrf=")

  req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  c, w, _ := GetMockContext(req)

  ap.DoLogin(c)

  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, "/assets/assets.html", w.Header().Get("Location"))
}

func TestLogout(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  ap := AuthApp{User_utils: auth_utils_mck}

  req, _ := http.NewRequest("GET", "/", bytes.NewBufferString(""))
  c, w, _ := GetMockContext(req)

  ap.Logout(c)

  assert.Equal(t, http.StatusFound, w.Code)
  assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestAddUserNoPwd(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  ap := AuthApp{User_utils: auth_utils_mck}

  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "")
  data.Set("_csrf", "csrf=")

  req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  c, w, _ := GetMockContext(req)

  ap.AddUser(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
  assert.Equal(t, "Wrong account or password", w.Body.String())
}

func TestAddUserHasUserFailed(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  err := errors.New("Get failed")
  auth_utils_mck.On("HasUser", mock.Anything, mock.Anything).Return(false, err)
  mockID := primitive.NewObjectID()
  auth_utils_mck.On("AddUser", mock.Anything, mock.Anything).Return(mockID, err)
  ap := AuthApp{User_utils: auth_utils_mck}

  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  data.Set("_csrf", "csrf=")

  req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  c, w, _ := GetMockContext(req)

  ap.AddUser(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddUserFailed(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  auth_utils_mck.On("HasUser", mock.Anything, mock.Anything).Return(false, nil)
  mockID := primitive.NewObjectID()
  err := errors.New("Get failed")
  auth_utils_mck.On("AddUser", mock.Anything, mock.Anything).Return(mockID, err)
  ap := AuthApp{User_utils: auth_utils_mck}

  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  data.Set("_csrf", "csrf=")

  req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  c, w, _ := GetMockContext(req)

  ap.AddUser(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddUser(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  auth_utils_mck.On("HasUser", mock.Anything, mock.Anything).Return(false, nil)
  mockID := primitive.NewObjectID()
  auth_utils_mck.On("AddUser", mock.Anything, mock.Anything).Return(mockID, nil)
  ap := AuthApp{User_utils: auth_utils_mck}

  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  data.Set("_csrf", "csrf=")

  req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  c, w, _ := GetMockContext(req)

  ap.AddUser(c)

  assert.Equal(t, "/", w.Header().Get("Location"))
}

func TestGetUser_By_AccountNoAccount(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User {
    ID: primitive.NewObjectID(),
    Account: "foo",
    Scopes: []primitive.ObjectID{},
  }
  auth_utils_mck.On("GetUserByAccount", mock.Anything).Return(mockUser, nil)
  ap := AuthApp{User_utils: auth_utils_mck}

  req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(""))
  c, w, _ := GetMockContext(req)

  ap.GetUser_By_Account(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUser_By_AccountFailed(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User {
    ID: primitive.NewObjectID(),
    Account: "foo",
    Scopes: []primitive.ObjectID{},
  }
  err := errors.New("Get failed")
  auth_utils_mck.On("GetUserByAccount", mock.Anything).Return(mockUser, err)
  ap := AuthApp{User_utils: auth_utils_mck}

  req, _ := http.NewRequest("POST", "/?account=foo", bytes.NewBufferString(""))
  c, w, _ := GetMockContext(req)

  ap.GetUser_By_Account(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUser_By_Account(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User {
    ID: primitive.NewObjectID(),
    Account: "foo",
    Scopes: []primitive.ObjectID{},
  }
  auth_utils_mck.On("GetUserByAccount", mock.Anything).Return(mockUser, nil)
  ap := AuthApp{User_utils: auth_utils_mck}

  req, _ := http.NewRequest("POST", "/?account=foo", bytes.NewBufferString(""))
  c, w, _ := GetMockContext(req)

  ap.GetUser_By_Account(c)

  assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser_ScopesWrongJSON(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User {
    ID: primitive.NewObjectID(),
    Account: "foo",
    Scopes: []primitive.ObjectID{},
  }
  auth_utils_mck.On("GetUserByID", mock.Anything).Return(mockUser, nil)
  auth_utils_mck.On("UpdateUser", mock.Anything).Return(nil)
  ap := AuthApp{User_utils: auth_utils_mck}

  req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(""))
  c, w, _ := GetMockContext(req)

  ap.UpdateUser_Scopes(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_ScopesGetUserFailed(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User {
    ID: primitive.NewObjectID(),
    Account: "foo",
    Scopes: []primitive.ObjectID{},
  }
  err := errors.New("Get failed")
  auth_utils_mck.On("GetUserByID", mock.Anything).Return(mockUser, err)
  auth_utils_mck.On("UpdateUser", mock.Anything).Return(nil)
  ap := AuthApp{User_utils: auth_utils_mck}

  testUser := ReducedUser {
    ID: mockUser.ID,
    Scopes: []primitive.ObjectID{
      primitive.NewObjectID(),
      primitive.NewObjectID(),
    },
  }
  json_bytes, _ := json.Marshal(testUser)
  json_buf := bytes.NewBuffer(json_bytes)

  req, _ := http.NewRequest("POST", "/", json_buf)
  c, w, _ := GetMockContext(req)

  ap.UpdateUser_Scopes(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser_ScopesFailed(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User {
    ID: primitive.NewObjectID(),
    Account: "foo",
    Scopes: []primitive.ObjectID{},
  }
  auth_utils_mck.On("GetUserByID", mock.Anything).Return(mockUser, nil)
  err := errors.New("Update failed")
  auth_utils_mck.On("UpdateUser", mock.Anything).Return(err)
  ap := AuthApp{User_utils: auth_utils_mck}

  testUser := ReducedUser {
    ID: mockUser.ID,
    Scopes: []primitive.ObjectID{
      primitive.NewObjectID(),
      primitive.NewObjectID(),
    },
  }
  json_bytes, _ := json.Marshal(testUser)
  json_buf := bytes.NewBuffer(json_bytes)

  req, _ := http.NewRequest("POST", "/", json_buf)
  c, w, _ := GetMockContext(req)

  ap.UpdateUser_Scopes(c)

  assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateUser_Scopes(t *testing.T) {
  auth_utils_mck := new(mockUserUtils)
  mockUser := auth.User {
    ID: primitive.NewObjectID(),
    Account: "foo",
    Scopes: []primitive.ObjectID{},
  }
  auth_utils_mck.On("GetUserByID", mock.Anything).Return(mockUser, nil)
  auth_utils_mck.On("UpdateUser", mock.Anything).Return(nil)
  ap := AuthApp{User_utils: auth_utils_mck}

  testUser := ReducedUser {
    ID: mockUser.ID,
    Scopes: []primitive.ObjectID{
      primitive.NewObjectID(),
      primitive.NewObjectID(),
    },
  }
  json_bytes, _ := json.Marshal(testUser)
  json_buf := bytes.NewBuffer(json_bytes)

  req, _ := http.NewRequest("POST", "/", json_buf)
  c, w, _ := GetMockContext(req)

  ap.UpdateUser_Scopes(c)

  assert.Equal(t, http.StatusOK, w.Code)
}
