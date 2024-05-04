package main

import (
  "context"
  "os"

  "github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/gin-contrib/sessions/mongo/mongodriver"
  "github.com/utrack/gin-csrf"
  "go.mongodb.org/mongo-driver/mongo"

  "github.com/starnight/risk_assessment/auth"
  "github.com/starnight/risk_assessment/config"
  "github.com/starnight/risk_assessment/middleware"
  "github.com/starnight/risk_assessment/risk_assessment"
  "github.com/starnight/risk_assessment/database"
)

type Apps struct {
  AuthApp IAuthApp
  ScopesApp IScopesApp
  AssetsApp IAssetsApp
}

func getSecretString() string {
  return "secret123"
}

func setupRouter(apps *Apps, store sessions.Store) *gin.Engine {
  r := gin.Default()

  r.Use(sessions.Sessions("sessionid", store))
  r.Use(csrf.Middleware(csrf.Options{
    Secret: getSecretString(),
    ErrorFunc: middleware.CSRFError,
  }))

  r.StaticFile("/", "./assets/index.html")
  r.Static("/assets", "./assets")

  public := r.Group("/")
  PublicAuthRoutes(public, apps.AuthApp)

  private := r.Group("/")
  private.Use(middleware.AuthenticationRequired)
  PrivateAuthRoutes(private, apps.AuthApp)
  ScopesRoutes(private, apps.ScopesApp)
  AssetsRoutes(private, apps.AssetsApp)

  privilege := r.Group("/")
  privilege.Use(middleware.AuthenticationRequired)
  privilege.Use(middleware.AuthorizationRequired)
  PrivilegeAuthRoutes(privilege, apps.AuthApp)

  return r
}

func prepareDb() *mongo.Client {
  return database.ConnectDB(database.GetDBStr(""))
}

func prepareSessionStore(db_client *mongo.Client) (sessions.Store) {
  if db_client == nil {
    return cookie.NewStore([]byte("secret"))
  }

  col := db_client.Database(config.DB_NAME).Collection("sessions")
  return mongodriver.NewStore(col, 3600, true, []byte("secret"))
}

func getPort() string {
  port := ":8080"
  /*
   * Azure Function will pass the forwarding port with environment variable
   * "FUNCTIONS_CUSTOMHANDLER_PORT"
   */
  if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
      port = ":" + val
  }
  return port
}

func main() {
  db_client := prepareDb()
  defer func() {
  if err := db_client.Disconnect(context.TODO()); err != nil {
      panic(err)
    }
  }()

  session_store := prepareSessionStore(db_client)
  csrf_utils := middleware.CsrfUtils{}

  user_utils := auth.UserUtils{ DB_Client: db_client }
  auth_ap := AuthApp{
    User_utils: &user_utils,
    Csrf_utils: &csrf_utils,
  }

  scope_utils := risk_assessment.ScopeUtils{ DB_Client: db_client }
  scopes_ap := ScopesApp{
    User_utils: &user_utils,
    Csrf_utils: &csrf_utils,
    Scope_utils: &scope_utils,
  }

  asset_utils := risk_assessment.AssetUtils{ DB_Client: db_client }
  assets_ap := AssetsApp{
    User_utils: &user_utils,
    Csrf_utils: &csrf_utils,
    Asset_utils: &asset_utils,
  }

  apps := Apps{ AuthApp: &auth_ap, ScopesApp: &scopes_ap, AssetsApp: &assets_ap }
  r := setupRouter(&apps, session_store)
  r.Run(getPort())
}
