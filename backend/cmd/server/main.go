package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/config"
	"github.com/tofex/backend/internal/database"
	"github.com/tofex/backend/internal/handlers"
)

func main() {
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
