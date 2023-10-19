package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	ID             uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name           string        `json:"name" gorm:"not null"`
	Description    string        `json:"description"`
	OrganizationID *uuid.UUID    `json:"organization_id" gorm:"not null"`
	Organization   *Organization `gorm:"foreignkey:OrganizationID" json:"organization"`
	Sender         *User         `gorm:"foreignkey:SenderID" json:"sender"`
	SenderID       uuid.UUID     `json:"sender_id" gorm:"not null"`
	State          string        `json:"state" gorm:"not null"`
	Votes          []Vote        `json:"votes"`
}

type TodoInput struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	OrganizationID *uuid.UUID `json:"organization_id"`
	SenderID       uuid.UUID  `json:"sender_id"`
	State          string     `json:"state"`
}

type TodoEditInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
}

type VoteInput struct {
	TodoID uuid.UUID `json:"todo_id"`
}

type Vote struct {
	gorm.Model
	SenderID uuid.UUID `json:"sender_id"`
	Sender   *User     `gorm:"foreignkey:SenderID" json:"sender"`
	TodoID   uuid.UUID `json:"todo_id"`
	Todo     *Todo     `gorm:"foreignkey:TodoID" json:"todo"`
}

type TodoChanges struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ManagerID    uuid.UUID `json:"manager_id" gorm:"not null"`
	Manager      *User     `gorm:"foreignKey:ManagerID" json:"manager"`
	ToDoID       uuid.UUID `json:"todo_id" gorm:"not null"`
	ChangedField string    `json:"changed_field"`
	OldValue     string    `json:"old_value"`
	NewValue     string    `json:"new_value"`
}
