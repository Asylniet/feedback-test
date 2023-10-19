package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type OrganizationRouteController struct {
	orgController controllers.OrganizationController
}

func NewOrganizationRouteController(orgController controllers.OrganizationController) OrganizationRouteController {
	return OrganizationRouteController{orgController}
}

func (rc *OrganizationRouteController) OrgRoute(rg *gin.RouterGroup) {
	router := rg.Group("organization", middleware.DeserializeUser("admin", "manager"))
	router2 := rg.Group("organization", middleware.DeserializeUser("any"))

	router.GET("/", rc.orgController.GetOrganizations)
	router.GET("/:id", rc.orgController.GetOrganization)
	router.GET("/:id/sub", rc.orgController.GetSubOrganizations)
	router.POST("/", rc.orgController.AddOrganization)
	router.PUT("/:id", rc.orgController.EditOrganization)
	router.DELETE("/:id", rc.orgController.DeleteOrganization)

	// New endpoint for rating sub-organizations
	router2.Use(middleware.AddFeedPoint(10)).POST("/:id/rate", rc.orgController.RateOrganization)
	router2.POST("/checkRate/:id", rc.orgController.CheckRate)
}
