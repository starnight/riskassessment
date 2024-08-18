package main

import (
  "encoding/hex"
  "net/http"
  "crypto/sha256"
  "strings"

  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"

  "go.mongodb.org/mongo-driver/bson/primitive"

  "github.com/starnight/risk_assessment/auth"
  "github.com/starnight/risk_assessment/middleware"
)

type IAuthApp interface {
  GetLogin(c *gin.Context)
  DoLogin(c *gin.Context)
  Logout(c *gin.Context)
  AddUser(c *gin.Context)
  GetUser_By_Account(c *gin.Context)
  UpdateUser_Scopes(c *gin.Context)
}

type AuthApp struct {
  User_utils auth.IUserUtils
  Csrf_utils middleware.ICSRFUtils
}

func (ap *AuthApp) GetLogin(c *gin.Context) {
  ap.Csrf_utils.AddCSRFToken(c)
  c.Status(http.StatusOK)
}

func (ap *AuthApp) DoLogin(c *gin.Context) {
  account := strings.TrimSpace(c.PostForm("account"))
  passwd := strings.TrimSpace(c.PostForm("passwd"))

  if (len(account) == 0 || len(passwd) == 0) {
    c.String(http.StatusForbidden, "Wrong account or password")
    return
  }

  pwd := sha256.Sum256([]byte(passwd))
  password := hex.EncodeToString(pwd[:])
  user, err := ap.User_utils.GetUserByAccountPwd(account, password)
  if (err != nil) {
    c.String(http.StatusForbidden, "Wrong account or password")
    return
  }

  session := sessions.Default(c)
  session.Set("id", user.ID.Hex())
  session.Set("role", user.Role)
  session.Save()

  c.Status(http.StatusOK)
}

func (ap *AuthApp) Logout(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()
  session.Options(sessions.Options{Path: "/", MaxAge: -1})
  session.Save()
  c.Redirect(http.StatusFound, "/")
}

func (ap *AuthApp) AddUser(c *gin.Context) {
  var user auth.User

  account := strings.TrimSpace(c.PostForm("account"))
  passwd := strings.TrimSpace(c.PostForm("passwd"))

  if (len(account) == 0 || len(passwd) == 0) {
    c.String(http.StatusBadRequest, "Wrong account or password")
    return
  }

  user.Account = account
  pwd := sha256.Sum256([]byte(passwd))
  user.Password = hex.EncodeToString(pwd[:])

  has, err := ap.User_utils.HasUser()
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  user.Role = auth.NormalUser
  if (!has) {
    user.Role = auth.Administrator
  }

  _, err = ap.User_utils.AddUser(&user)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  c.Redirect(http.StatusFound, "/")
}

type ReducedUser struct {
  ID primitive.ObjectID `bson:"_id"`
  Account string
  Role uint
  Scopes []primitive.ObjectID
}

func (ap *AuthApp) GetUser_By_Account(c *gin.Context) {
  account := strings.TrimSpace(c.Query("account"))
  if (account == "") {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  user, err := ap.User_utils.GetUserByAccount(account)
  if (err != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  var reduced_user ReducedUser
  reduced_user.ID = user.ID
  reduced_user.Account = user.Account
  reduced_user.Role = user.Role
  reduced_user.Scopes = user.Scopes

  c.JSON(http.StatusOK, reduced_user)
}

func (ap *AuthApp) UpdateUser_Scopes(c *gin.Context) {
  var reduced_user ReducedUser

  if (c.BindJSON(&reduced_user) != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  user, err := ap.User_utils.GetUserByID(reduced_user.ID)
  if (err != nil) {
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  user.Role = reduced_user.Role
  user.Scopes = reduced_user.Scopes
  err = ap.User_utils.UpdateUser(&user)
  if (err != nil) {
    c.AbortWithStatus(http.StatusInternalServerError)
    return
  }

  c.Status(http.StatusOK)
}
