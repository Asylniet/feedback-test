package middleware

import (
	"github.com/enzhas/feedback_back/initializers"
	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
)

func AddFeedPoint(feedPoints uint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if ctx.Writer.Status() >= 200 && ctx.Writer.Status() <= 300 {

			user := ctx.MustGet("currentUser").(models.User)
			user.FeedPoints += feedPoints
			initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&user)
		}
	}
}
