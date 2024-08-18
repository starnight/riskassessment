package main

import (
  "github.com/gin-gonic/gin"
)

func PublicAuthRoutes (g *gin.RouterGroup, ap IAuthApp) {
  g.GET("/api/login", ap.GetLogin)
  g.POST("/api/login", ap.DoLogin)
  g.POST("/api/register", ap.AddUser)
}

func PrivateAuthRoutes (g *gin.RouterGroup, ap IAuthApp) {
  g.GET("/api/logout", ap.Logout)
}

func ScopesRoutes (g *gin.RouterGroup, ap IScopesApp) {
  g.GET("/api/getscopes", ap.GetScopes)
  g.GET("/api/getscopesbyuser", ap.GetScopesByUser)
  g.POST("/api/addscope", ap.AddScope)
  g.POST("/api/updatescope", ap.UpdateScope)
}

func AssetsRoutes (g *gin.RouterGroup, ap IAssetsApp) {
  g.GET("/api/getassets/:scopeID", ap.GetAssets)
  g.POST("/api/addasset", ap.AddAsset)
  g.POST("/api/updateasset", ap.UpdateAsset)
  g.POST("/api/deleteasset", ap.DeleteAsset)
}

func PrivilegeAuthRoutes (g *gin.RouterGroup, ap IAuthApp) {
  g.GET("/api/getuser_by_account", ap.GetUser_By_Account)
  g.POST("/api/updateuser_scopes", ap.UpdateUser_Scopes)
}
