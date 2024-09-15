package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/starnight/riskassessment/backend/auth"
  "github.com/starnight/riskassessment/backend/middleware"
  "github.com/starnight/riskassessment/backend/risk_assessment"
)

type IAssetsApp interface {
  GetAssets(c *gin.Context)
  AddAsset(c *gin.Context)
  UpdateAsset(c *gin.Context)
  DeleteAsset(c *gin.Context)
}

type AssetsApp struct {
  User_utils auth.IUserUtils
  Csrf_utils middleware.ICSRFUtils
  Asset_utils risk_assessment.IAssetUtils
}

type AssetPage struct {
  UserInfo auth.UserInfo
  Assets []risk_assessment.Asset
}

type assetReq struct {
  ScopeID primitive.ObjectID `binding:"required"`
  Asset risk_assessment.Asset
}

func (ap *AssetsApp) GetAssets(c *gin.Context) {
  var asset_page AssetPage
  var authorized bool
  var u_id primitive.ObjectID
  var s_id primitive.ObjectID
  var err error

  session := sessions.Default(c)
  userID := session.Get("id").(string)
  u_id, _ = primitive.ObjectIDFromHex(userID)
  asset_page.UserInfo.Role = session.Get("role").(uint)

  scopeID := c.Param("scopeID")
  s_id, err = primitive.ObjectIDFromHex(scopeID)
  if (err != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  authorized, err = ap.User_utils.UserHasScopeID(u_id, s_id)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  } else if (!authorized) {
    c.AbortWithStatus(http.StatusForbidden)
    return
  }

  asset_page.Assets, err = ap.Asset_utils.GetAssetsByScopeID(s_id)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  ap.Csrf_utils.AddCSRFToken(c)
  c.JSON(http.StatusOK, asset_page)
}

func (ap *AssetsApp) AddAsset(c *gin.Context) {
  var asset risk_assessment.Asset

  session := sessions.Default(c)
  userID := session.Get("id").(string)
  u_id, _ := primitive.ObjectIDFromHex(userID)

  if (c.BindJSON(&asset) != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  authorized, err := ap.User_utils.UserHasScopeID(u_id, asset.Scope)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  } else if (!authorized) {
    c.AbortWithStatus(http.StatusForbidden)
    return
  }

  err = ap.Asset_utils.AddAsset(&asset)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }
  c.Status(http.StatusOK)
}

func (ap *AssetsApp) UpdateAsset(c *gin.Context) {
  var asset risk_assessment.Asset
  var authorized bool

  session := sessions.Default(c)
  userID := session.Get("id").(string)
  u_id, _ := primitive.ObjectIDFromHex(userID)

  if (c.BindJSON(&asset) != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  orig_asset, err := ap.Asset_utils.GetAssetByID(asset.ID)
  if (err != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  authorized, err = ap.User_utils.UserHasScopeID(u_id, orig_asset.Scope)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  } else if (!authorized) {
    c.AbortWithStatus(http.StatusForbidden)
    return
  }

  asset.CreateTime = orig_asset.CreateTime
  asset.Scope = orig_asset.Scope
  err = ap.Asset_utils.UpdateAsset(&asset)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusOK)
}

type ID struct {
  Id primitive.ObjectID `json:"id"`
}

func (ap *AssetsApp) DeleteAsset(c *gin.Context) {
  var id ID
  var authorized bool

  session := sessions.Default(c)
  userID := session.Get("id").(string)
  u_id, _ := primitive.ObjectIDFromHex(userID)

  if (c.BindJSON(&id) != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  asset, err := ap.Asset_utils.GetAssetByID(id.Id)
  if (err != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  authorized, err = ap.User_utils.UserHasScopeID(u_id, asset.Scope)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  } else if (!authorized) {
    c.AbortWithStatus(http.StatusForbidden)
    return
  }

  err = ap.Asset_utils.DeleteAsset(id.Id)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusOK)
}
