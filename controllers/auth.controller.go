package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/enzhas/feedback_back/initializers"
	"github.com/enzhas/feedback_back/models"
	"github.com/enzhas/feedback_back/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// SignUpSender Sign Up User
func (ac *AuthController) SignUpSender(ctx *gin.Context) {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	verified := config.GinMode == "debug"
	(*AuthController).SignUpUser(ac, ctx, 2, verified, "local")
}

func (ac *AuthController) SignUpReceiver(ctx *gin.Context) {
	currentRole := ctx.MustGet("currentRole").(models.Role)
	(*AuthController).SignUpUser(ac, ctx, 3, true, currentRole.Name)
}

func (ac *AuthController) SignUpManager(ctx *gin.Context) {
	(*AuthController).SignUpUser(ac, ctx, 4, true, "admin")
}

func (ac *AuthController) SignUpAdmin(ctx *gin.Context) {
	(*AuthController).SignUpUser(ac, ctx, 1, true, "admin")
}

func (ac *AuthController) SendEmail(ctx *gin.Context, userID uuid.UUID) {
	var newUser models.User
	if result := ac.DB.Where("id = ?", userID).First(&newUser); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something went wrong("})
		return
	}
	config, _ := initializers.LoadConfig()
	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verificationCode
	ac.DB.Save(newUser)

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verify/" + verificationCode,
		FirstName: firstName,
		Subject:   "Ссылка верификации",
	}

	if err := utils.SendEmail(&newUser, &emailData, "verificationCode.html"); err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context, roleId uint, verified bool, provider string) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if !utils.IsValidEmail(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Mail should be @kbtu.kz"})
		return
	}

	if !utils.IsValidPassword(payload.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords should contain:\nUppercase letters: A-Z\nLowercase letters: a-z\nNumbers: 0-9\n"})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var organization models.Organization
	organizationID := payload.OrganizationID
	if organizationID == nil {
		mail := strings.Split(payload.Email, "@")[1]
		if result := ac.DB.Where("email LIKE ?", mail).First(&organization); result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Could not register user under some organization"})
			return
		}
		organizationID = &organization.ID
	}

	newUser := models.User{
		Name:           payload.Name,
		Surname:        payload.Surname,
		Email:          strings.ToLower(payload.Email),
		Password:       hashedPassword,
		RoleID:         roleId,
		Verified:       verified, //false for sending mail
		Photo:          payload.Photo,
		Provider:       provider,
		OrganizationID: organizationID,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	if newUser.Verified == false {
		ac.SendEmail(ctx, newUser.ID)

		message := "We sent an email with a verification code to " + newUser.Email
		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
	} else {
		ReturnUsers(ctx, ac.DB, strconv.Itoa(int(newUser.RoleID)), (*newUser.OrganizationID).String())
	}
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Enter correct data"})
		return
	}

	var user *models.User
	result := ac.DB.Preload("Role").Where("email = ?", strings.ToLower(payload.Email)).First(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, _ := initializers.LoadConfig()

	// Generate Tokens
	accessToken, err := utils.GenerateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.GenerateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{
		"role":              user.Role.Name,
		"access_token":      accessToken,
		"refresh_token":     refreshToken,
		"access_token_age":  config.AccessTokenMaxAge * 60,
		"refresh_token_age": config.RefreshTokenMaxAge * 60,
	}})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	var refreshToken = ""
	authorizationHeader := ctx.Request.Header.Get("Authorization")

	if fields := strings.Fields(authorizationHeader); len(fields) != 0 && fields[0] == "Bearer" {
		refreshToken = fields[1]
	}

	if refreshToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
		return
	}
	config, _ := initializers.LoadConfig()

	sub, err := utils.ValidateToken(refreshToken, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	accessToken, err := utils.GenerateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "token": gin.H{
		"access_token":     accessToken,
		"access_token_age": config.AccessTokenMaxAge * 60,
	}})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var payload *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
		return
	}

	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)
	user.PasswordResetToken = passwordResetToken
	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	ac.DB.Save(&user)

	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       "" + "/api/auth/resetpassword/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	if err := utils.SendEmail(&user, &emailData, "resetPassword.html"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(payload.Password)

	passwordResetToken := utils.Encode(resetToken)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	updatedUser.Password = hashedPassword
	updatedUser.PasswordResetToken = ""
	ac.DB.Save(&updatedUser)
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("token", "", -1, "/", "", true, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)
	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "verification_code = ?", verificationCode)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if updatedUser.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	updatedUser.VerificationCode = ""
	updatedUser.Verified = true
	ac.DB.Save(&updatedUser)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}
