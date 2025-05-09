package main

import (
	"github.com/wahid-wex/e-commerce-backend/api"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/cache"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/db/migrations"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {

	cfg := config.GetConfig()

	logger := logging.NewLogger(cfg)

	err := cache.InitRedis(cfg)

	defer cache.CloseRedis()

	if err != nil {
		println("init redis failed.", err.Error())
		logger.Fatal(logging.Redis, logging.Startup, err.Error(), nil)
	}

	err = db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}
	migrations.Up()

	api.InitServer(cfg)
}
