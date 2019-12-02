package models

import "time"

// base model
type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	UID       int        `json:"uid"`
	CreatedAt time.Time  `gorm:"column:created_at;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}
