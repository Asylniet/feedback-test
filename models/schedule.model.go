package models

import "github.com/google/uuid"

type Schedule struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	SenderID  uuid.UUID `json:"sender_id"`
	SubjectID uint      `json:"subject_id"`
	Subject   Subject   `json:"-"`
}
