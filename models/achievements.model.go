package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Achievement struct {
	gorm.Model
	Photo          string        `json:"photo"`
	Title          string        `json:"title"`
	OrganizationId *uuid.UUID    `json:"organization_id"`
	Organization   *Organization `json:"organization"`
	CategoryId     *uint         `json:"category_id"`
	Category       *Category     `json:"category"`
	RoleId         *uint         `json:"role_id"`
	Role           *Role         `json:"role"`
	Points         uint          `json:"points"`
}
type AchievementInput struct {
	Title          string     `json:"title"`
	CategoryId     *uint      `json:"category_id"`
	OrganizationId *uuid.UUID `json:"organization_id"`
	RoleId         *uint      `json:"role_id"`
	Points         uint       `json:"points"`
}
