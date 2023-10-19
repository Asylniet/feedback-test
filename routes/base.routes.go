package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type BaseRouteController struct {
	baseController *controllers.BaseController
}

func NewBaseRouteController(baseController *controllers.BaseController) *BaseRouteController {
	return &BaseRouteController{baseController}
}

func (bc *BaseRouteController) BaseRoute(router *gin.RouterGroup) {
	router.POST("/upload", middleware.DeserializeUser("any"), bc.baseController.UploadFile)
	router.GET("/files/:part/:filename", bc.baseController.GetUploadedFile)
}
