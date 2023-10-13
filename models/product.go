package models

import (
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ProductID               int            `json:"product_id" gorm:"primary_key"`
	ProductName             string         `json:"product_name" gorm:"not null"`
	ProductDescription      string         `json:"product_description"`
	ProductImages           pq.StringArray `json:"product_images" gorm:"type:text[]"`
	ProductPrice            float64        `json:"product_price" gorm:"not null"`
	CompressedProductImages pq.StringArray `json:"compressed_product_images" gorm:"type:text[]"`
	CreatedAt               time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt               time.Time      `json:"updated_at"`
}
