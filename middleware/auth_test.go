package middleware

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"

  "strings"
  "testing"
  "net/http/httptest"
  "github.com/stretchr/testify/assert"

  "github.com/starnight/risk_assessment/auth"
)

func PublicRoutes (g *gin.RouterGroup) {
  g.POST("/login", func (c *gin.Context) {
    session := sessions.Default(c)
    session.Set("account", "foo")
    session.Set("id", uint(1))
    session.Set("role", uint(auth.Administrator))
    session.Save()
    c.Status(http.StatusOK)
  })
}

func PrivateRoutes (g *gin.RouterGroup) {
  g.GET("/private", func (c *gin.Context){
    session := sessions.Default(c)
    account := session.Get("account").(string)
    c.String(http.StatusOK, account)
  })
}

func copyCookies(req *http.Request, res *httptest.ResponseRecorder) {
  req.Header.Set("Cookie", strings.Join(res.Header().Values("Set-Cookie"), "; "))
}

func TestAuthenticationRequired(t *testing.T) {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

  public := r.Group("/")
  PublicRoutes(public)

  private := r.Group("/")
  private.Use(AuthenticationRequired)
  PrivateRoutes(private)

  /* Must get forbidden, because has not login */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/private", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusUnauthorized, res1.Code)
  assert.Equal(t, "Please login first", res1.Body.String())

  /* Login */
  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("POST", "/login", nil)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusOK, res2.Code)

  /* Then, request again with the cookie for the session */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/private", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.Equal(t, "foo", res3.Body.String())
}

func TestAuthorizationRequired(t *testing.T) {
  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))

  public := r.Group("/")
  PublicRoutes(public)

  private := r.Group("/")
  private.Use(AuthorizationRequired)
  PrivateRoutes(private)

  /* Must get forbidden, because has not login */
  res1 := httptest.NewRecorder()
  req1, _ := http.NewRequest("GET", "/private", nil)
  r.ServeHTTP(res1, req1)

  assert.Equal(t, http.StatusForbidden, res1.Code)
  assert.Equal(t, "Permission denied", res1.Body.String())

  /* Login */
  res2 := httptest.NewRecorder()
  req2, _ := http.NewRequest("POST", "/login", nil)
  r.ServeHTTP(res2, req2)

  assert.Equal(t, http.StatusOK, res2.Code)

  /* Then, request again with the cookie for the session */
  res3 := httptest.NewRecorder()
  req3, _ := http.NewRequest("GET", "/private", nil)
  copyCookies(req3, res2)
  r.ServeHTTP(res3, req3)

  assert.Equal(t, http.StatusOK, res3.Code)
  assert.Equal(t, "foo", res3.Body.String())
}
