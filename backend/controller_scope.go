package main

import (
  "strings"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/starnight/risk_assessment/auth"
  "github.com/starnight/risk_assessment/middleware"
  "github.com/starnight/risk_assessment/risk_assessment"
)

type IScopesApp interface {
  GetScopes(c *gin.Context)
  GetScopesByUser(c *gin.Context)
  AddScope(c *gin.Context)
  UpdateScope(c *gin.Context)
}

type ScopesApp struct {
  User_utils auth.IUserUtils
  Csrf_utils middleware.ICSRFUtils
  Scope_utils risk_assessment.IScopeUtils
}

type ReducedScope struct {
  ID primitive.ObjectID `bson:"_id"`
  Name string `binding:"required"`
}

type ScopePage struct {
  UserInfo auth.UserInfo
  Scopes []ReducedScope
}

func (ap *ScopesApp) GetScopes(c *gin.Context) {
  var scope_page ScopePage

  session := sessions.Default(c)
  scope_page.UserInfo.Role = session.Get("role").(uint)
  if (scope_page.UserInfo.Role != auth.Administrator) {
    c.AbortWithStatus(http.StatusForbidden)
    return
  }

  scopes, err := ap.Scope_utils.GetScopes()
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  for _, scope := range scopes {
    scope_page.Scopes = append(scope_page.Scopes,
                               ReducedScope{
                                 ID: scope.ID,
				 Name: scope.Name,
                               })
  }

  ap.Csrf_utils.AddCSRFToken(c)
  c.JSON(http.StatusOK, scope_page)
}

func (ap *ScopesApp) GetScopesByUser(c *gin.Context) {
  session := sessions.Default(c)
  id_str := session.Get("id")
  id, _ := primitive.ObjectIDFromHex(id_str.(string))

  user, err := ap.User_utils.GetUserByID(id)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  var scope_page ScopePage
  var scopes []risk_assessment.Scope
  scope_page.UserInfo.Role = session.Get("role").(uint)
  scopes, err = ap.Scope_utils.GetScopeByIDs(user.Scopes)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  for _, scope := range scopes {
    scope_page.Scopes = append(scope_page.Scopes,
                               ReducedScope{
                                 ID: scope.ID,
				 Name: scope.Name,
                               })
  }

  ap.Csrf_utils.AddCSRFToken(c)
  c.JSON(http.StatusOK, scope_page)
}

func (ap *ScopesApp) AddScope(c *gin.Context) {
  var scope risk_assessment.Scope

  if (c.BindJSON(&scope) != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  scope.Name = strings.TrimSpace(scope.Name)
  if (len(scope.Name) == 0)  {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  _, err := ap.Scope_utils.AddScope(&scope)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusOK)
}

func (ap *ScopesApp) UpdateScope(c *gin.Context) {
  var scope risk_assessment.Scope

  if (c.BindJSON(&scope) != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  scope.Name = strings.TrimSpace(scope.Name)
  if (len(scope.Name) == 0)  {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  orig_scope, err := ap.Scope_utils.GetScopeByID(scope.ID)
  if (err != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  scope.CreateTime = orig_scope.CreateTime
  err = ap.Scope_utils.UpdateScope(&scope)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusOK)
}
