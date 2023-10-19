package controllers

import (
	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type AchievementsController struct {
	DB *gorm.DB
}

func NewAchievementsController(DB *gorm.DB) AchievementsController {
	return AchievementsController{DB}
}
func (ac *AchievementsController) GetAllAchievements(ctx *gin.Context) {
	var achievements []models.Achievement

	if err := ac.DB.Find(&achievements).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "no achievements found"})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"achievements": achievements}})
}

func (ac *AchievementsController) CreateAchievement(ctx *gin.Context) {
	var achievementInput models.AchievementInput
	if err := ctx.ShouldBindJSON(&achievementInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
	}
	achievement := models.Achievement{
		Title:          achievementInput.Title,
		CategoryId:     achievementInput.CategoryId,
		OrganizationId: achievementInput.OrganizationId,
		Photo:          "",
		RoleId:         achievementInput.RoleId,
		Points:         achievementInput.Points,
	}

	if err := ac.DB.Create(&achievement).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to create achievement"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"achievement": achievement}})
}

//Should I update the photo? to-do

func (ac *AchievementsController) UpdateAchievement(ctx *gin.Context) {
	id := ctx.Param("id")
	var achievement models.Achievement
	if err := ac.DB.First(&achievement, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "achievement not found"})
	}

	var newAchievement models.AchievementInput
	if err := ctx.ShouldBindJSON(&newAchievement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	achievement.Title = newAchievement.Title
	achievement.Points = newAchievement.Points

	if err := ac.DB.Save(&achievement).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to update achievement"})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"achievement": achievement}})
}
func (ac *AchievementsController) DeleteAchievement(ctx *gin.Context) {
	id := ctx.Param("id")

	var achievement models.Category

	if err := ac.DB.First(&achievement, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "achievement not found"})
	}
	if err := ac.DB.Delete(&achievement).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "failed to load achievement"})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "achievement deleted successfully"})
}
