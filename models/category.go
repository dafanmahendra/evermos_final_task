package models

import "time"

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"type:varchar(255)" json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}