package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type CategoryRouteController struct {
	categoryController controllers.CategoryController
}

func NewCategoryRouteController(categoryController controllers.CategoryController) CategoryRouteController {
	return CategoryRouteController{categoryController}
}

func (rc *CategoryRouteController) CategoryRoute(rg *gin.RouterGroup) {
	router := rg.Group("organization/category")

	router.GET("/", middleware.DeserializeUser("admin", "manager"), rc.categoryController.GetAllCategories)
	router.GET("/:id", rc.categoryController.GetCategory)
	router.GET("/search", rc.categoryController.SearchCategories)
	router.POST("/", middleware.DeserializeUser("admin", "manager"), rc.categoryController.CreateCategory)
	router.PUT("/:id", middleware.DeserializeUser("admin", "manager"), rc.categoryController.UpdateCategory)
	router.DELETE("/:id", middleware.DeserializeUser("admin", "manager"), rc.categoryController.DeleteCategory)
}
