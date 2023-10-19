package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type AchievementsRouteController struct {
	achievementsController controllers.AchievementsController
}

func NewAchievementsController(achievementsController controllers.AchievementsController) AchievementsRouteController {
	return AchievementsRouteController{achievementsController}
}
func (ac *AchievementsRouteController) AchievementsRoute(rg *gin.RouterGroup) {
	router := rg.Group("achievements", middleware.DeserializeUser("admin", "manager"))

	router.GET("/", ac.achievementsController.GetAllAchievements)
	router.POST("/", ac.achievementsController.CreateAchievement)
	router.PUT("/:id", ac.achievementsController.UpdateAchievement)
	router.DELETE("/:id", ac.achievementsController.DeleteAchievement)
}
