package controllers

import (
	"errors"
	"github.com/enzhas/feedback_back/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type OrganizationController struct {
	DB *gorm.DB
}

func NewOrganizationController(DB *gorm.DB) OrganizationController {
	return OrganizationController{DB}
}

func (oc *OrganizationController) GetOrganization(ctx *gin.Context) {
	id := ctx.Param("id")
	var organization models.Organization
	if err := oc.DB.Preload("Category").Preload("Root").Preload("Parent").Where("id = ?", id).First(&organization); err.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"organization": organization}})
}

func (oc *OrganizationController) GetOrganizations(ctx *gin.Context) {
	var organizations []models.Organization

	// Preload the Category field while querying
	if err := oc.DB.Preload("Category").Where("level = ?", 0).Find(&organizations); err.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"organizations": organizations}})
}

func (oc *OrganizationController) GetSubOrganizations(ctx *gin.Context) {
	id := ctx.Param("id")
	level := ctx.Query("level")
	var organizations []models.Organization
	if len(level) > 0 {
		if err := oc.DB.Preload("Category").Preload("Root").Preload("Parent").Where("parent_id =? OR root_id = ? AND level = ?", id, id, level).Find(&organizations); err.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
			return
		}
	} else {
		if err := oc.DB.Preload("Category").Preload("Root").Preload("Parent").Where("parent_id =?", id).Find(&organizations); err.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Not found"})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"organizations": organizations}})
}

func (oc *OrganizationController) AddOrganization(ctx *gin.Context) {
	var organizationInput models.OrganizationInput
	if err := ctx.ShouldBindJSON(&organizationInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02 15:04:05"
	expireTime, err := time.Parse(layout, organizationInput.SubExpireAt)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organization := models.Organization{
		Name:             organizationInput.Name,
		ShortName:        organizationInput.ShortName,
		CategoryID:       organizationInput.CategoryID,
		Photo:            "",
		Email:            organizationInput.Email,
		CoolDownDuration: organizationInput.CoolDownDuration,
		SubExpireAt:      expireTime,
	}

	var parentOrg models.Organization
	if organizationInput.ParentID != nil {
		if result := oc.DB.First(&parentOrg, organizationInput.ParentID); result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "No parent organization found"})
			return
		}
		organization.ParentID = organizationInput.ParentID
		if parentOrg.RootID != nil {
			organization.RootID = parentOrg.RootID
		} else {
			organization.RootID = &parentOrg.ID
		}
		organization.Email = parentOrg.Email
		organization.Level = parentOrg.Level + 1
	} else {
		organization.Level = 0
	}

	if result := oc.DB.Create(&organization); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	if err := oc.DB.Model(&organization).Preload("Category").Preload("Parent").Preload("Root").First(&organization).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"organization": organization}})
}

func (oc *OrganizationController) EditOrganization(ctx *gin.Context) {
	id := ctx.Param("id")
	type OrganizationProp struct {
		Name         string `json:"name" binding:"required"`
		Category     string `json:"category" binding:"required"`
		ShortName    string `json:"shortName" binding:"required"`
		Email        string `json:"email"`
		SubExpiresAt string `json:"subscription_expire_at"`
	}
	var organizationProp OrganizationProp
	var organization models.Organization
	var organizations []models.Organization
	if err := ctx.ShouldBindJSON(&organizationProp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02 15:04:05"
	expireTime, err := time.Parse(layout, organizationProp.SubExpiresAt)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := oc.DB.Where("id = ?", id).First(&organization); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "No organization found"})
		return
	}

	organization.Name = organizationProp.Name
	//category? need to check if this category exists
	organization.ShortName = organizationProp.ShortName
	organization.Email = organizationProp.Email
	organization.SubExpireAt = expireTime

	if result := oc.DB.Save(&organization); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}

	if organization.ParentID != nil {
		if result := oc.DB.Preload("Root").Preload("Parent").Where("parent_id = ?", organization.ParentID).Find(&organizations); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "No Something went wrong"})
			return
		}
	} else {
		if result := oc.DB.Preload("Root").Preload("Parent").Where("level = ?", 0).Find(&organizations); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "No Something went wrong"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"organizations": organizations}})
}

func (oc *OrganizationController) DeleteOrganization(ctx *gin.Context) {
	var organization models.Organization
	var organizations []models.Organization
	id := ctx.Param("id")
	if result := oc.DB.Where("id = ?", id).Find(&organization); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}
	parentId := organization.ParentID
	if result := oc.DB.Where("id = ?", id).Delete(&organization); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Something went wrong"})
		return
	}
	if parentId != nil {
		if result := oc.DB.Preload("Root").Preload("Parent").Where("parent_id = ?", parentId).Find(&organizations); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "No Something went wrong"})
			return
		}
	} else {
		if result := oc.DB.Preload("Root").Preload("Parent").Where("level = ?", 0).Find(&organizations); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "No Something went wrong"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"organizations": organizations}})
}

func (oc *OrganizationController) RateOrganization(ctx *gin.Context) {
	rootId := ctx.Param("id")
	subOrgId := ctx.Query("sub_id")

	var input models.RatingInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := ctx.MustGet("currentUser").(models.User)

	if user.RoleID == 4 {
		if (*user.OrganizationID).String() == rootId {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "You can't rate your own organization"})
			return
		}
	}

	var id string
	if subOrgId != "" {
		id = subOrgId
	} else {
		id = rootId
	}

	var organization models.Organization
	if result := oc.DB.Where("id = ?", id).First(&organization); result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sub-organization ID"})
		return
	}

	if subOrgId != "" {
		if organization.RootID.String() != rootId {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Sub-organization is not a child of this main organization"})
			return
		}
	}

	str := oc.checkCoolDown(&user, &organization)
	if str != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": str})
		return
	}

	rating := &models.Rating{
		Rating:     input.Rating,
		Comment:    input.Comment,
		SenderID:   user.ID,
		ReceiverID: organization.ID,
	}

	if res := oc.DB.Create(rating); res.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error in creating a new rating"})
		return
	}

	if err := oc.updateOrganizationRating(&organization, input.Rating, 1); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rating submitted successfully"})
}

func (oc *OrganizationController) CheckRate(ctx *gin.Context) {
	id := ctx.Param("id")
	user := ctx.MustGet("currentUser").(models.User)
	var org models.Organization

	if err := oc.DB.Where("id = ?", id).First(&org).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Organization was not found"})
		return
	}

	str := oc.checkCoolDown(&user, &org)
	if str != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": str})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "As of now you can rate this organization"})
}

func (oc *OrganizationController) updateOrganizationRating(subOrganization *models.Organization, sum float64, count uint64) error {
	var totalRating models.TotalRating

	if err := oc.DB.Where("receiver_id = ?", subOrganization.ID).First(&totalRating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			totalRating = models.TotalRating{
				ReceiverID: subOrganization.ID,
				Count:      count,
				Sum:        sum,
				Rating:     sum / float64(count),
			}
			oc.DB.Create(&totalRating)
		} else {
			return err
		}
	} else {
		totalRating.Count += count
		totalRating.Sum += sum
		totalRating.Rating = totalRating.Sum / float64(totalRating.Count)
		oc.DB.Model(&models.TotalRating{}).Where("receiver_id = ?", totalRating.ReceiverID).Updates(&totalRating)
	}

	if subOrganization.Level == 0 {
		return nil
	}

	var organization models.Organization
	if err := oc.DB.First(&organization, subOrganization.ParentID).Error; err != nil {
		return err
	}

	if err := oc.updateOrganizationRating(&organization, sum, count); err != nil {
		return err
	}

	return nil
}

func (oc *OrganizationController) checkCoolDown(user *models.User, organization *models.Organization) string {
	var rating models.Rating

	if err := oc.DB.Where("sender_id = ? and receiver_id = ?", user.ID, organization.ID).Last(&rating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ""
		} else {
			return err.Error()
		}
	}

	coolDownDuration := time.Duration(organization.CoolDownDuration) * time.Second
	coolDownEndTime := rating.CreatedAt.Add(coolDownDuration)
	if time.Now().After(coolDownEndTime) {
		return ""
	}

	remainingTime := time.Until(coolDownEndTime)
	return "Try again in " + remainingTime.String()
}
