package routes

import (
	"github.com/enzhas/feedback_back/controllers"
	"github.com/enzhas/feedback_back/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("auth")

	router.POST("/register", rc.authController.SignUpSender)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", rc.authController.LogoutUser)

	router.POST("/registerReceiver", middleware.DeserializeUser("manager", "admin"), rc.authController.SignUpReceiver)
	router.POST("/registerManager", middleware.DeserializeUser("admin"), rc.authController.SignUpManager)
	router.POST("/registerAdmin", middleware.DeserializeUser("admin"), rc.authController.SignUpAdmin)

	router.POST("/forgotpassword", rc.authController.ForgotPassword)
	router.PATCH("/resetpassword/:resetToken", rc.authController.ResetPassword)

	router.GET("/verifyemail/:verificationCode", rc.authController.VerifyEmail)
}
