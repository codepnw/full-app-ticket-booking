package routes

import (
	"github.com/codepnw/ticket-api/cmd/config"
	"github.com/codepnw/ticket-api/handlers"
	userRepository "github.com/codepnw/ticket-api/repositories/user"
	userService "github.com/codepnw/ticket-api/services/user"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func UserRoutes(cfg config.EnvConfig, db *sqlx.DB, r *gin.Engine, version string) {
	repo := userRepository.NewUserRepository(db)
	service := userService.NewUserService(cfg, repo)
	handler := handlers.NewUserHandler(service)

	g := r.Group(version + "/users")

	g.POST("/", handler.SignupUser)
	g.POST("/admin", handler.SignupAdmin)
	g.GET("/:user_id", handler.GetProfile)
	g.POST("/signin", handler.SignIn)
	g.POST("/refresh", handler.RefreshPassport)
}
