package middleware

import (
	"net/http"
	"time"

	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
)

func CheckSubscriptionExpireAt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("currentUser").(models.User)
		if user.RoleID != 1 {

			if user.Organization.SubExpireAt.IsZero() {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Your organization need to have subscription"})
				return
			} else if time.Now().After(user.Organization.SubExpireAt) {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Subscription expired. Currently you can not use this."})
				return
			}
		} else {
			ctx.Next()
		}

	}
}
