package main

import (
  "net/http"
  "net/http/httptest"

  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
)

func _SetupMockSession(c *gin.Context) gin.HandlerFunc {
  session_store := prepareSessionStore(nil)
  return sessions.Sessions("sessionid", session_store)
}

func GetMockContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder, sessions.Session) {
  w := httptest.NewRecorder()
  c, _ := gin.CreateTestContext(w)
  c.Request = req
  cb := _SetupMockSession(c)
  cb(c)
  session := sessions.Default(c)
  return c, w, session
}
