package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type BonusRouteController struct {
	bonusController controllers.BonusController
}

func NewBonusController(bonusController controllers.BonusController) BonusRouteController {
	return BonusRouteController{bonusController}
}

func (bc *BonusRouteController) BonusesRoute(rg *gin.RouterGroup) {
	router := rg.Group("bonus", middleware.DeserializeUser("admin", "manager"))

	router.GET("/", bc.bonusController.GetAllBonuses)
	router.POST("/", bc.bonusController.CreateBonus)
	router.PUT("/:id", bc.bonusController.UpdateBonus)
	router.DELETE("/:id", bc.bonusController.DeleteBonus)
	router.GET("/:id", bc.bonusController.GetBonus)
}
