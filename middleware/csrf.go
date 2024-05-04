package middleware

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/utrack/gin-csrf"
)

type ICSRFUtils interface {
  AddCSRFToken(c *gin.Context)
}

type CsrfUtils struct {}

func (ap *CsrfUtils) AddCSRFToken(c *gin.Context) {
  c.Writer.Header().Set("X-CSRF-TOKEN", csrf.GetToken(c))
}

func CSRFError(c *gin.Context) {
  c.AbortWithStatus(http.StatusBadRequest)
}
