package models

import (
	"gorm.io/gorm"
)

type Bonus struct {
	gorm.Model
	AchievementId uint        `json:"achievement_id" `
	Achievement   Achievement `json:"achievement" gorm:"foreignkey:AchievementId"`
	Rating        Rating      `gorm:"foreignkey:RatingId"`
	RatingId      uint        `json:"rating_id"`
}
