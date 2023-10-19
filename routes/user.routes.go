package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")
	router.GET("/me", middleware.DeserializeUser("any"), uc.userController.GetMe)
	router.POST("/estimate", middleware.DeserializeUser("any"), middleware.CheckSubscriptionExpireAt(), uc.userController.Estimate)
	router.POST("/addSubject", middleware.DeserializeUser("any"), middleware.CheckSubscriptionExpireAt(), uc.userController.AddSubject)
	router.POST("/addSchedule", middleware.DeserializeUser("manager"), uc.userController.AddSchedule)
	router.POST("/upload", middleware.DeserializeUser("manager"), uc.userController.FastSignUpReceiver)
	router.GET("/mySchedule", middleware.DeserializeUser("sender"), uc.userController.Schedule)

	router.GET("/:id", middleware.DeserializeUser("admin"), uc.userController.GetUser)
	router.GET("/", middleware.DeserializeUser("admin", "manager"), uc.userController.GetUsers)
	router.PUT("/:id", middleware.DeserializeUser("admin", "manager"), uc.userController.EditUser)
	router.DELETE("/:id", middleware.DeserializeUser("admin", "manager"), uc.userController.DeleteUser)
}
