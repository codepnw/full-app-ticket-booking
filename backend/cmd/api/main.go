package main

import (
	"github.com/codepnw/ticket-api/cmd/config"
	"github.com/codepnw/ticket-api/cmd/database"
	"github.com/codepnw/ticket-api/cmd/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	envFile = "dev.env"
	version = "v1"
)

func main() {
	envConfig := config.NewEnvConfig(envFile)

	db, err := database.InitDatabase(envConfig)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(cors.Default())
	routes.EventRoutes(db, r, version)

	r.Run(":" + envConfig.ServerPort)
}
