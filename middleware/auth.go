package middleware

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"

  "github.com/starnight/risk_assessment/auth"
)

func AuthenticationRequired(c *gin.Context) {
  session := sessions.Default(c)

  id := session.Get("id")
  if (id == nil) {
    c.String(http.StatusUnauthorized, "Please login first")
    c.Abort()
    return
  }

  c.Next()
}

func AuthorizationRequired(c *gin.Context) {
  session := sessions.Default(c)

  role := session.Get("role")
  if (role == nil || role.(uint) != auth.Administrator) {
    c.String(http.StatusForbidden, "Permission denied")
    c.Abort()
    return
  }

  c.Next()
}
