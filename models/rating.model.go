package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	Rating     float64
	Comment    string        `gorm:"type:text" json:"comment"`
	SenderID   uuid.UUID     `json:"sender_id"`
	ReceiverID uuid.UUID     `json:"receiver_id"`
	Sender     User          `gorm:"foreignkey:SenderID"`
	Receiver   *Organization `gorm:"foreignkey:ReceiverID"`
}

type TotalRating struct {
	ReceiverID uuid.UUID     `json:"receiver_id"`
	Receiver   *Organization `gorm:"foreignkey:ReceiverID"`
	Rating     float64       `gorm:"default:0"`
	Count      uint64        `gorm:"default:0"`
	Sum        float64       `gorm:"default:0"`
}

// need to change
type RatingInput struct {
	Rating  float64 `json:"rating" binding:"required,min=1,max=5"`
	Comment string  `json:"comment"`
}
