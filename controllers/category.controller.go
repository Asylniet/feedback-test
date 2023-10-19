package controllers

import (
	"fmt"
	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(DB *gorm.DB) CategoryController {
	return CategoryController{DB}
}

func (cc *CategoryController) CreateCategory(ctx *gin.Context) {
	var category models.Category
	var RootOrganization models.Organization
	fmt.Println(1)
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cc.DB.Where("id = ?", category.RootOrganizationID).First(&RootOrganization)
	if RootOrganization.RootID != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "you cannot create category for child organizations"})
		return
	}
	fmt.Println(3)

	if err := cc.DB.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to create category"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"category": category}})
}

func (cc *CategoryController) GetAllCategories(ctx *gin.Context) {
	user := ctx.MustGet("currentUser").(models.User)
	var categories []models.Category
	if err := cc.DB.Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No categories found"})
		return
	}
	var localCategories []models.Category
	var generalCategories []models.Category

	// manager
	if user.RoleID == 4 {
		for _, category := range categories {
			if category.RootOrganizationID == nil {
				generalCategories = append(generalCategories, category)
			}
		}
		cc.DB.Where("root_organization_id = ?", user.Organization.RootID).Find(&localCategories)
	} else if user.RoleID == 1 {
		for _, category := range categories {
			if category.RootOrganizationID == nil {
				generalCategories = append(generalCategories, category)
			}
		}
	}

	data := map[string]interface{}{
		"localCategories":   localCategories,
		"generalCategories": generalCategories,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (cc *CategoryController) GetCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var category models.Category
	if err := cc.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Category not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"category": category}})
}

func (cc *CategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var category models.Category
	if err := cc.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Category not found"})
		return
	}

	var newCategory models.Category
	if err := ctx.ShouldBindJSON(&newCategory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = newCategory.Name
	category.RootOrganizationID = newCategory.RootOrganizationID
	if err := cc.DB.Save(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to update category"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"category": category}})
}

func (cc *CategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var category models.Category
	if err := cc.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Category not found"})
		return
	}

	if err := cc.DB.Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to delete category"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Category deleted successfully"})
}
func (cc *CategoryController) SearchCategories(ctx *gin.Context) {
	// Get the search query parameter from the URL
	query := ctx.Query("query")

	var categories []models.Category

	// Perform a database query to search for categories
	if err := cc.DB.Where("name ILIKE ?", "%"+query+"%").Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"categories": categories}})
}
