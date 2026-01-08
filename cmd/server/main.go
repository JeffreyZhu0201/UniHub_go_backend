package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"unihub/internal/config"
	"unihub/internal/db"
	"unihub/internal/model"
	"unihub/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gormDB, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	// Auto-migrate schemas on startup.
	if err := model.AutoMigrate(gormDB); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()

	router.Register(engine, cfg, gormDB)

	if err := engine.Run(cfg.Server.Port); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
