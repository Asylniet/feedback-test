package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/enzhas/feedback_back/initializers"
	"github.com/enzhas/feedback_back/models"
	"github.com/enzhas/feedback_back/utils"
	"github.com/gin-gonic/gin"
)

func DeserializeUser(roles ...interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config, _ := initializers.LoadConfig()
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPrivateKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User
		result := initializers.DB.Preload("Organization").First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}
		var role models.Role
		initializers.DB.First(&role, "id = ?", user.RoleID)
		for _, Role := range roles {
			if role.Name == Role || Role == "any" {
				ctx.Set("currentUser", user)
				ctx.Set("currentRole", role)
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User must have roles: " + fmt.Sprintf("%v", roles),
		})

	}
}
