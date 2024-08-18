package middleware

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/utrack/gin-csrf"

  "testing"
  "net/http/httptest"
  "github.com/stretchr/testify/assert"
)

func TestAddCSRFToken(t *testing.T) {
  csrf_utils := CsrfUtils{}

  r := gin.Default()
  store := cookie.NewStore([]byte("secret"))
  r.Use(sessions.Sessions("sessionid", store))
  r.Use(csrf.Middleware(csrf.Options{
    Secret: "secret123",
  }))
  r.Use(csrf_utils.AddCSRFToken)

  r.GET("/", func (c *gin.Context) {
    c.Status(http.StatusOK)
  })

  res := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "/", nil)
  r.ServeHTTP(res, req)

  assert.Equal(t, http.StatusOK, res.Code)
  assert.True(t, len(res.Header().Get("X-CSRF-TOKEN")) > 0)
}

func TestCSRFError(t *testing.T) {
  gin.SetMode(gin.TestMode)
  w := httptest.NewRecorder()
  c, _ := gin.CreateTestContext(w)
  c.Request = httptest.NewRequest("POST", "/", nil)

  CSRFError(c)

  assert.Equal(t, http.StatusBadRequest, w.Code)
}
