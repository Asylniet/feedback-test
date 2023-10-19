package controllers

import (
	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type BonusController struct {
	DB *gorm.DB
}

func NewBonusController(DB *gorm.DB) BonusController {
	return BonusController{DB}
}

func (bc *BonusController) GetAllBonuses(ctx *gin.Context) {
	var bonuses []models.Bonus

	if err := bc.DB.Find(&bonuses).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "no bonuses found"})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"bonuses": bonuses}})
}

func (bc *BonusController) CreateBonus(ctx *gin.Context) {
	var bonus models.Bonus
	if err := ctx.ShouldBindJSON(&bonus); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := bc.DB.Create(&bonus).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to create bonus"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"bonus": bonus}})
}

//Should I update the photo? to-do

func (bc *BonusController) UpdateBonus(ctx *gin.Context) {
	id := ctx.Param("id")
	var bonus models.Bonus
	if err := bc.DB.First(&bonus, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "bonus not found"})
	}

	var newBonus models.Bonus
	if err := ctx.ShouldBindJSON(&newBonus); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if newBonus.AchievementId != 0 && newBonus.RatingId != 0 {
		bonus.RatingId = newBonus.RatingId
		bonus.AchievementId = newBonus.AchievementId
	} else if newBonus.AchievementId == 0 {
		bonus.RatingId = newBonus.RatingId
	} else {
		bonus.AchievementId = newBonus.AchievementId
	}

	if err := bc.DB.Save(&bonus).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to update bonus"})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"bonus": bonus}})
}
func (bc *BonusController) DeleteBonus(ctx *gin.Context) {
	id := ctx.Param("id")

	var bonus models.Bonus

	if err := bc.DB.First(&bonus, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "bonus not found"})
	}
	if err := bc.DB.Delete(&bonus).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "failed to load bonus"})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "bonus deleted successfully"})
}

func (bc *BonusController) GetBonus(ctx *gin.Context) {
	id := ctx.Param("id")

	var bonus models.Bonus

	if err := bc.DB.First(&bonus, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "bonus not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"bonus": bonus}})
}
