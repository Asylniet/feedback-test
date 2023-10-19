package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subject struct {
	gorm.Model
	UserID    uuid.UUID
	User      User
	Year      string `json:"year"`
	Semester  string `json:"semester"`
	Day       string `json:"day" gorm:"not null"`
	StartTime uint8  `json:"start_time" gorm:"not null"`
	EndTime   uint8  `json:"end_time"  gorm:"not null"`
}

type Attendance struct {
	gorm.Model
	SenderID  uuid.UUID `json:"sender_id"`
	SubjectID uint
	Subject   Subject
}
