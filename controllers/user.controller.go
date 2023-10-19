package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/enzhas/feedback_back/models"
	"github.com/enzhas/feedback_back/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if err := uc.DB.Preload("Role").Preload("Organization").Where("id", id).Find(&user); err.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return
	}
	response := models.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Surname:      user.Surname,
		Email:        user.Email,
		Photo:        user.Photo,
		Role:         user.Role.Name,
		Provider:     user.Provider,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Organization: user.Organization,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": response}})
}

func ReturnUsers(ctx *gin.Context, db *gorm.DB, role, organization string) {
	var users []models.User
	if len(role) > 0 && len(organization) > 0 {
		if err := db.Preload("Role").Preload("Organization").Where("role_id = ? AND organization_id = ?", role, organization).Find(&users); err.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
			return
		}
	} else if len(role) > 0 {
		if err := db.Preload("Role").Preload("Organization").Where("role_id = ?", role).Find(&users); err.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
			return
		}
	} else if len(organization) > 0 {
		if err := db.Preload("Role").Preload("Organization").Where("organization_id = ?", organization).Find(&users); err.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
			return
		}
	} else {
		if err := db.Preload("Role").Preload("Organization").Find(&users); err.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
			return
		}
	}

	response := make([]models.UserResponse, len(users))

	for i, item := range users {
		response[i] = models.UserResponse{
			ID:           item.ID,
			Name:         item.Name,
			Surname:      item.Surname,
			Email:        item.Email,
			Photo:        item.Photo,
			Role:         item.Role.Name,
			Provider:     item.Provider,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
			Organization: item.Organization,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"users": response}})
}

func (uc *UserController) GetUsers(ctx *gin.Context) {
	role := ctx.Query("role_id")
	organization := ctx.Query("organization_id")
	ReturnUsers(ctx, uc.DB, role, organization)
}

func (uc *UserController) EditUser(ctx *gin.Context) {
	id := ctx.Param("id")
	type UserInput struct {
		Name           string `json:"name" binding:"required"`
		Surname        string `json:"surname" binding:"required"`
		Photo          string `json:"photo" binding:"required"`
		Email          string `json:"email"`
		OrganizationID *uint  `json:"organization_id"`
		RoleID         uint   `json:"role_id"`
	}
	var userInput UserInput
	var user models.User

	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if result := uc.DB.Where("id = ?", id).First(&user); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "No user found"})
		return
	}

	user.Name = userInput.Name
	user.Surname = userInput.Surname
	if len(userInput.Email) > 0 {
		user.Email = userInput.Email
	}
	user.Photo = userInput.Photo

	if result := uc.DB.Save(&user); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}

	roleID := strconv.Itoa(int(userInput.RoleID))
	var orgID string
	if userInput.OrganizationID != nil {
		orgID = strconv.Itoa(int(*userInput.OrganizationID))
	}
	ReturnUsers(ctx, uc.DB, roleID, orgID)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if result := uc.DB.Where("id = ?", id).First(&user); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Could not delete user"})
		return
	}

	if result := uc.DB.Where("id = ?", id).Delete(&user); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Could not delete user"})
		return
	}

	roleID := strconv.Itoa(int(user.RoleID))
	var orgID string
	if user.OrganizationID != nil {
		orgID = (*user.OrganizationID).String()
	}

	ReturnUsers(ctx, uc.DB, roleID, orgID)
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	currentRole := ctx.MustGet("currentRole").(models.Role)
	// var currentRole models.Role
	// uc.DB.First(&currentRole, currentUser.RoleID)

	userResponse := &models.UserResponse{
		ID:           currentUser.ID,
		Name:         currentUser.Name,
		Surname:      currentUser.Surname,
		Email:        currentUser.Email,
		Photo:        currentUser.Photo,
		Role:         currentRole.Name,
		Organization: currentUser.Organization,
		Provider:     currentUser.Provider,
		CreatedAt:    currentUser.CreatedAt,
		UpdatedAt:    currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

func (uc *UserController) Estimate(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var rating models.Rating

	if len(currentUser.Name) < 1 {
		if result := uc.DB.Where("provider = guest").Find(&currentUser); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong"})
			return
		}
	}

	if err := ctx.ShouldBindJSON(&rating); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if rating.Rating > 5 || rating.Rating <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid data, rating should be less than 6 and more than 0"})
		return
	}
	rating.SenderID = currentUser.ID

	result := uc.DB.Create(&rating)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid rating data"})
		return
	}

	updRating := models.TotalRating{
		ReceiverID: rating.ReceiverID,
	}

	if uc.DB.Model(&updRating).Where("id = ?", rating.ReceiverID).Updates(&updRating).RowsAffected == 0 {
		uc.DB.Create(&updRating)
	}
	// uc.DB.First(&updRating, "id = ?", rating.ReceiverID)
	updRating.Count++
	updRating.Rating += rating.Rating
	uc.DB.Save(&updRating)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (uc *UserController) AddSubject(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var subject models.Subject

	if err := ctx.ShouldBindJSON(&subject); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if (subject.StartTime < 0 && subject.StartTime > 24) || (subject.EndTime < 0 && subject.EndTime > 24) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid data, time should be between 0 and 24"})
		return
	}

	subject.UserID = currentUser.ID

	result := uc.DB.Create(&subject)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid rating data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (uc *UserController) FastSignUpReceiver(c *gin.Context) {
	// Get the uploaded file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Failed to retrieve the uploaded file"})
		return
	}

	// Create a new file on the server to store the uploaded file
	newUuid := uuid.New()
	extension := filepath.Ext(file.Filename)
	newFileName := newUuid.String() + "." + extension
	f, err := os.Create(newFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create the file on the server"})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to open the uploaded file"})
		return
	}

	// Copy the contents of the uploaded file to the newly created file on the server
	_, err = io.Copy(f, src)
	if err := src.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save the file on the server"})
		return
	}
	if err := f.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save the file on the server"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to save the file on the server"})
		return
	}
	receivers, err := excelize.OpenFile(newFileName)
	sheet := receivers.GetSheetName(0)
	for i := 2; ; i++ {
		name, err := receivers.GetCellValue(sheet, fmt.Sprintf("A%d", i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("problem on line %d with name", i)})
			return
		}
		surname, err := receivers.GetCellValue(sheet, fmt.Sprintf("B%d", i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("problem on line %d with surname", i)})
			return
		}
		email, err := receivers.GetCellValue(sheet, fmt.Sprintf("C%d", i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("problem on line %d with email", i)})
			return
		}
		password, err := receivers.GetCellValue(sheet, fmt.Sprintf("D%d", i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("problem on line %d with password", i)})
			return
		}
		if len(name) == 0 {
			break
		}
		// (*AuthController).FastSignUpReceiver(&AuthController{}, c, name, surname, email, password)
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": fmt.Sprintf("problem on line %d %s", i, err.Error())})
			return
		}
		if !utils.IsValidEmail(email) {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": fmt.Sprintf("Mail should is not valid on line %d with email", i)})
			return
		}

		if !utils.IsValidPassword(password) {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": fmt.Sprintf("Invalid on line %d", i)})
			return
		}
		newReceiver := models.User{
			Name:     name,
			Surname:  surname,
			Email:    email,
			Password: hashedPassword,
			RoleID:   3,
			Verified: true,
			Photo:    "default.png",
			Provider: "manager",
		}

		if uc.DB.Model(&newReceiver).Where("email = ?", email).Updates(&newReceiver).RowsAffected == 0 {
			uc.DB.Create(&newReceiver)
		}
	}
	// Return a success message
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "File uploaded successfully"})
	if err := receivers.Close(); err != nil {
		return
	}

	if err := os.Remove(newFileName); err != nil {
		return
	}
}

func (uc *UserController) AddSchedule(ctx *gin.Context) {
	// currentUser := ctx.MustGet("currentUser").(models.User)
	var schedule models.Schedule

	if err := ctx.ShouldBindJSON(&schedule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := uc.DB.Create(&schedule)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (uc *UserController) Schedule(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var schedule []models.Schedule

	result := uc.DB.Find(&schedule, "sender_id = ?", currentUser.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid data"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"schedule": schedule}})
}
