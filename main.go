package main

import (
	"log"
	"net/http"

	"github.com/enzhas/feedback_back/migrate"

	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/initializers"
	"github.com/enzhas/feedback_back/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// nodemon --exec go run main.go

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	OrgController      controllers.OrganizationController
	OrgRouteController routes.OrganizationRouteController

	TodoController      controllers.TodoController
	TodoRouteController routes.TodoRouteController

	CategoryController      controllers.CategoryController
	CategoryRouteController routes.CategoryRouteController

	BaseRouteController *routes.BaseRouteController

	AchievementsController      controllers.AchievementsController
	AchievementsRouteController routes.AchievementsRouteController

	BonusController      controllers.BonusController
	BonusRouteController routes.BonusRouteController
)

func init() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	OrgController = controllers.NewOrganizationController(initializers.DB)
	OrgRouteController = routes.NewOrganizationRouteController(OrgController)

	TodoController = controllers.NewTodoController(initializers.DB)
	TodoRouteController = routes.NewTodoRouteController(TodoController)

	CategoryController = controllers.NewCategoryController(initializers.DB)
	CategoryRouteController = routes.NewCategoryRouteController(CategoryController)

	BaseController := controllers.NewBaseController(initializers.DB)
	BaseRouteController = routes.NewBaseRouteController(&BaseController)

	AchievementsController = controllers.NewAchievementsController(initializers.DB)
	AchievementsRouteController = routes.NewAchievementsController(AchievementsController)

	BonusController = controllers.NewBonusController(initializers.DB)
	BonusRouteController = routes.NewBonusController(BonusController)

	gin.SetMode(config.GinMode)
	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	migrate.Migrate()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "https://feedback-asylniet.vercel.app", config.ClientOrigin}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Cookie", "Access-Control-Allow-Origin"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Feedback! " + config.ClientOrigin
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	OrgRouteController.OrgRoute(router)
	TodoRouteController.TodoRoute(router)
	CategoryRouteController.CategoryRoute(router)
	BaseRouteController.BaseRoute(router)
	AchievementsRouteController.AchievementsRoute(router)
	BonusRouteController.BonusesRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
