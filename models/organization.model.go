package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Organization struct {
	gorm.Model
	ID               uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name             string        `json:"name" gorm:"not null"`
	Rating           int           `json:"rating"`
	ParentID         *uuid.UUID    `json:"parent_id"`
	RootID           *uuid.UUID    `json:"root_id"`
	CategoryID       uint          `json:"category_id"`
	Parent           *Organization `gorm:"foreignkey:ParentID" json:"parent"`
	Root             *Organization `gorm:"foreignkey:RootID" json:"root"`
	Category         Category      `gorm:"foreignkey:CategoryID" json:"category" `
	ShortName        string        `json:"shortName" gorm:"not null"`
	Level            int8          `json:"level" gorm:"not null"`
	Photo            string        `json:"photo" gorm:"not null"`
	Email            string        `json:"email"`
	CoolDownDuration uint64        `json:"cool_down_duration"`
	SubExpireAt      time.Time     `json:"subscription_expire_at"`
}

type Category struct {
	gorm.Model
	Name               string        `json:"name"`
	RootOrganizationID *uuid.UUID    `json:"root_organization_id"`
	RootOrganization   *Organization `gorm:"foreignkey:RootOrganizationID" json:"root_organization"`
}

type OrganizationInput struct {
	Name             string     `json:"name" binding:"required"`
	ParentID         *uuid.UUID `json:"parent_id"`
	CategoryID       uint       `json:"category_id" binding:"required"`
	ShortName        string     `json:"shortName" binding:"required"`
	Email            string     `json:"email"`
	CoolDownDuration uint64     `json:"cool_down_duration"`
	SubExpireAt      string     `json:"subscription_expire_at"`
}

type OrganizationDTO struct {
	ID               uint             `json:"id"`
	Name             string           `json:"name"`
	Parent           *OrganizationDTO `json:"parent,omitempty"`
	ShortName        string           `json:"shortName"`
	Category         CategoryDTO      `json:"category"`
	Level            int8             `json:"level"`
	Rating           int8             `json:"rating"`
	Photo            int8             `json:"photo"`
	Email            string           `json:"email"`
	CoolDownDuration uint64           `json:"cool_down_duration"`
	SubExpireAt      time.Time        `json:"subscription_expire_at"`
}

type CategoryDTO struct {
	ID                 uint          `json:"id"`
	Name               string        `json:"name"`
	RootOrganizationID *uuid.UUID    `json:"root_organization_id"`
	RootOrganization   *Organization `json:"root_organization"`
}
