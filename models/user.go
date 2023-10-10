package models

import "time"

type User struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Name      string     `json:"user_name" gorm:"not null"`
	Mobile    string     `json:"user_mobile" gorm:"not null"`
	Latitude  float64    `json:"user_latitude" gorm:"not null"`
	Longitude float64    `json:"user_longitude" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time  `json:"updated_at"`
}
