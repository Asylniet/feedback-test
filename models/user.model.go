package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                 uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name               string       `gorm:"type:varchar(255);not null" json:"name"`
	Surname            string       `gorm:"type:varchar(255);not null" json:"surname"`
	Email              string       `gorm:"uniqueIndex;not null" json:"email"`
	Password           string       `gorm:"not null" json:"password"`
	RoleID             uint         `json:"role_id"`
	Role               Role         `gorm:"not null" json:"role"`
	Provider           string       `gorm:"not null" json:"provider"`
	Photo              string       `gorm:"not null" json:"photo"`
	VerificationCode   string       `json:"verification_code"`
	PasswordResetToken string       `json:"password_reset_token"`
	PasswordResetAt    time.Time    `json:"password_reset_at"`
	Verified           bool         `gorm:"not null" json:"verified"`
	OrganizationID     *uuid.UUID   `json:"organization_id"`
	Organization       Organization `json:"organization"`
	FeedPoints         uint         `json:"feed_points" gorm:"default:0"`
	//Bonuses            []Bonus      `json:"bonuses"`
	// Sender             []*Rating `gorm:"foreignkey:SenderID"`
	// Receiver           []*Rating `gorm:"foreignkey:ReceiverID"`
}

type SignUpInput struct {
	Name            string     `json:"name" binding:"required"`
	Surname         string     `json:"surname" binding:"required"`
	Email           string     `json:"email" binding:"required"`
	Password        string     `json:"password" binding:"required,min=8"`
	PasswordConfirm string     `json:"passwordConfirm" binding:"required"`
	Photo           string     `json:"photo"`
	OrganizationID  *uuid.UUID `json:"organization_id"`
	RoleID          uint       `json:"role_id"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID           uuid.UUID    `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Surname      string       `json:"surname,omitempty"`
	Email        string       `json:"email,omitempty"`
	Role         string       `json:"role,omitempty"`
	Photo        string       `json:"photo,omitempty"`
	Provider     string       `json:"provider"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Organization Organization `json:"organization"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
