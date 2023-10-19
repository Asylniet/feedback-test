package controllers

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BaseController struct {
	DB *gorm.DB
}

func NewBaseController(DB *gorm.DB) BaseController {
	return BaseController{DB}
}

func (bc *BaseController) UploadFile(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	id := ctx.Query("id")
	userID := ctx.Query("user_id")
	achievementId := ctx.Query("achievement_id")
	var path string
	var err error

	if id != "" && (currentUser.RoleID == 1 || currentUser.RoleID == 4) {
		var organization models.Organization
		if result := bc.DB.Where("id = ?", id).First(&organization); result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Not found"})
			return
		}

		if organization.Photo != "" {
			oldPath := filepath.Join(organization.Photo)
			if _, err = os.Stat(oldPath); err == nil {
				if err = os.Remove(oldPath); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Could not update photo"})
					return
				}
			}
		}

		path, err = uploadFile(ctx, "organizations")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization.Photo = path
		bc.DB.Save(&organization)

	} else if userID != "" {
		if currentUser.Photo != "" {
			oldPath := filepath.Join(currentUser.Photo)
			if _, err := os.Stat(oldPath); err == nil {
				if err = os.Remove(oldPath); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Could not update photo"})
					return
				}
			}
		}

		path, err = uploadFile(ctx, "users")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		currentUser.Photo = path
		bc.DB.Save(&currentUser)
	} else if achievementId != "" && (currentUser.RoleID == 1 || currentUser.RoleID == 4) {
		var achievement models.Achievement

		if result := bc.DB.Where("id = ?", achievementId).First(&achievement); result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "not found an achievement"})
			return
		}

		if achievement.Photo != "" {
			oldPath := filepath.Join(achievement.Photo)
			if _, err = os.Stat(oldPath); err == nil {
				if err = os.Remove(oldPath); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "could not update photo"})
					return
				}
			}
		}

		path, err = uploadFile(ctx, "achievements")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		achievement.Photo = path
		bc.DB.Save(&achievement)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "You don't have enough permission to change organization avatar"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"photo": path}})
}

func uploadFile(ctx *gin.Context, part string) (string, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return "", err
	}

	randomFileName := randomFileName()
	ext := filepath.Ext(file.Filename)

	path := filepath.Join("files", part, randomFileName+ext)

	// Create directory if doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return "", err
	}

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		return "", err
	}

	return path, nil
}

func randomFileName() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(52)]
	}
	return string(b)
}

func (bc *BaseController) GetUploadedFile(ctx *gin.Context) {
	part := ctx.Param("part")
	fileName := ctx.Param("filename")

	path := filepath.Join("files", part, fileName)

	// Check exists or not
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "error": "File not found"})
		return
	}

	ctx.File(path)
}
