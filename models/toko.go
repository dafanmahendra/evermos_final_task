package models

import "time"

type Toko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IdUser    uint      `json:"id_user"` // Foreign Key ke User
	NamaToko  string    `json:"nama_toko"`
	UrlToko   string    `json:"url_toko"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Toko) TableName() string {
	return "toko"
}