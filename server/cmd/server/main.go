// @title						Tofex API
// @version					1.0
// @description				REST API for Tofex restaurant/confectionery chain (staff admin + public guest ordering).
// @host						localhost:8080
// @BasePath					/
// @schemes					http https
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				JWT access token. Format: Bearer {token}
//
//go:generate swag init -g main.go -o ../../docs --parseDependency --parseInternal --dir .,../../internal/handlers,../../internal/apidocs
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/tofex/backend/docs"
	"github.com/tofex/backend/internal/config"
	"github.com/tofex/backend/internal/database"
	"github.com/tofex/backend/internal/handlers"
)

func main() {
	config.LoadDotEnv()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	db, err := database.Connect(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if err := database.SeedRBAC(db); err != nil {
		log.Fatalf("seed rbac: %v", err)
	}
	if err := database.EnsureBranchPermissions(db); err != nil {
		log.Fatalf("branch permissions: %v", err)
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("upload dir: %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.MaxMultipartMemory = int64(cfg.MaxUploadMB) << 20

	handlers.RegisterRoutes(r, db, cfg)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
