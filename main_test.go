package main

import (
  "net/http"
  "net/http/httptest"
  "net/url"
  "os"
  "strings"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"

  "github.com/utrack/gin-csrf"

  "github.com/starnight/risk_assessment/auth"
  "github.com/starnight/risk_assessment/middleware"
  "github.com/starnight/risk_assessment/risk_assessment"
)

type mockAuthApp struct {}

func (m *mockAuthApp) GetLogin(c *gin.Context) {
  c.String(http.StatusOK, csrf.GetToken(c))
}

func (m *mockAuthApp) DoLogin(c *gin.Context) {
  session := sessions.Default(c)
  session.Set("id", "foo")
  session.Set("role", uint(auth.Administrator))
  session.Save()
  c.Redirect(http.StatusFound, "/assets/assets.html")
}

func (m *mockAuthApp) Logout(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()
  session.Options(sessions.Options{Path: "/", MaxAge: -1})
  session.Save()
  c.Redirect(http.StatusFound, "/")
}

func (m *mockAuthApp) AddUser(c *gin.Context){
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockAuthApp) GetUser_By_Account(c *gin.Context){
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockAuthApp) UpdateUser_Scopes(c *gin.Context){
  c.String(http.StatusOK, c.Request.URL.Path)
}

type mockScopesApp struct {
  Scope_utils risk_assessment.IScopeUtils
}

func (m *mockScopesApp) GetScopes(c *gin.Context) {
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockScopesApp) GetScopesByUser(c *gin.Context) {
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockScopesApp) AddScope(c *gin.Context) {
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockScopesApp) UpdateScope(c *gin.Context) {
  c.String(http.StatusOK, c.Request.URL.Path)
}

type mockAssetsApp struct {
  Asset_utils risk_assessment.IAssetUtils
}

func (m *mockAssetsApp) GetAssets(c *gin.Context) {
  csrf_utils := middleware.CsrfUtils{}
  csrf_utils.AddCSRFToken(c)
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockAssetsApp) AddAsset(c *gin.Context) {
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockAssetsApp) UpdateAsset(c *gin.Context) {
  c.String(http.StatusOK, c.Request.URL.Path)
}

func (m *mockAssetsApp) DeleteAsset(c *gin.Context) {
  c.String(http.StatusOK, c.Request.URL.Path)
}

func copyCookies(req *http.Request, res *httptest.ResponseRecorder) {
  req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func TestSetupRouter(t *testing.T) {
  session_store := prepareSessionStore(nil)
  auth_ap := mockAuthApp{}
  scope_ap := mockScopesApp{}
  assets_ap := mockAssetsApp{}
  apps := Apps{AuthApp: &auth_ap, ScopesApp: &scope_ap, AssetsApp: &assets_ap}
  r := setupRouter(&apps, session_store)

  /* Get CSRF token for Login */
  w := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/api/login", nil)
  r.ServeHTTP(w, req)
  assert.Equal(t, http.StatusOK, w.Code)
  csrf_token := w.Body.String()
  assert.True(t, len(csrf_token) > 0)

  /* Login */
  data := url.Values{}
  data.Set("account", "foo")
  data.Set("passwd", "bar")
  data.Set("_csrf", csrf_token)

  w1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("POST", "/api/login", strings.NewReader(data.Encode()))
  req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  copyCookies(req1, w)
  r.ServeHTTP(w1, req1)
  assert.Equal(t, http.StatusFound, w1.Code)
  assert.Equal(t, "/assets/assets.html", w1.Header().Get("Location"))

  /* Get Assets */
  w2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("GET", "/api/getassets/xxxaa", nil)
  copyCookies(req2, w1)
  r.ServeHTTP(w2, req2)
  assert.Equal(t, http.StatusOK, w2.Code)
  assert.Equal(t, "/api/getassets/xxxaa", w2.Body.String())
  csrf_token = w2.Header().Get("X-CSRF-TOKEN")

  /* Register an User */
  w3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("POST", "/api/register", nil)
  req3.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req3, w1)
  r.ServeHTTP(w3, req3)
  assert.Equal(t, http.StatusOK, w3.Code)
  assert.Equal(t, "/api/register", w3.Body.String())

  /* Add an Asset */
  w4 := httptest.NewRecorder()
  req4, _ := http.NewRequest("POST", "/api/addasset", nil)
  req4.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req4, w1)
  r.ServeHTTP(w4, req4)
  assert.Equal(t, http.StatusOK, w4.Code)
  assert.Equal(t, "/api/addasset", w4.Body.String())

  /* Update an Asset */
  w5 := httptest.NewRecorder()
  req5, _ := http.NewRequest("POST", "/api/updateasset", nil)
  req5.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req5, w1)
  r.ServeHTTP(w5, req5)
  assert.Equal(t, http.StatusOK, w5.Code)
  assert.Equal(t, "/api/updateasset", w5.Body.String())

  /* Delete an Asset */
  w6 := httptest.NewRecorder()
  req6, _ := http.NewRequest("POST", "/api/deleteasset", nil)
  req6.Header.Set("X-CSRF-TOKEN", csrf_token)
  copyCookies(req6, w1)
  r.ServeHTTP(w6, req6)
  assert.Equal(t, http.StatusOK, w6.Code)
  assert.Equal(t, "/api/deleteasset", w6.Body.String())

  /* Logout */
  w7 := httptest.NewRecorder()
  req7, _ := http.NewRequest("GET", "/api/logout", nil)
  copyCookies(req7, w1)
  r.ServeHTTP(w7, req7)
  assert.Equal(t, http.StatusFound, w7.Code)
  assert.Equal(t, "/", w7.Header().Get("Location"))

  /* Get User by Account */
  w8 := httptest.NewRecorder()
  req8, _ := http.NewRequest("GET", "/api/getuser_by_account?account=foo", nil)
  copyCookies(req8, w1)
  r.ServeHTTP(w8, req8)
  assert.Equal(t, http.StatusOK, w8.Code)
  assert.Equal(t, "/api/getuser_by_account", w8.Body.String())
}

func TestGetPort(t *testing.T) {
  os.Setenv("FUNCTIONS_CUSTOMHANDLER_PORT", "9090")
  res1 := getPort()
  assert.Equal(t, ":9090", res1)

  os.Unsetenv("FUNCTIONS_CUSTOMHANDLER_PORT")
  res2 := getPort()
  assert.Equal(t, ":8080", res2)
}
