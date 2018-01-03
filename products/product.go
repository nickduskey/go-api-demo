package products

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Product represents a product
type Product struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
}

// GetSingle fetches a product from the db
func (p *Product) GetSingle(db *gorm.DB) error {
	return db.First(&p).Error
}

// Update updates a product by id in the db
func (p *Product) Update(db *gorm.DB) error {
	return db.Save(&p).Error
}

// Delete deletes a product
func (p *Product) Delete(db *gorm.DB) error {
	return db.Delete(&p).Error
}

// Create creates a product
func (p *Product) Create(db *gorm.DB) error {
	return db.Create(&p).Error
}

// Get fetches products from the db
func Get(db *gorm.DB, offset, limit int) ([]Product, error) {
	var p []Product
	if err := db.Offset(offset).Limit(limit).Find(&p).Error; err != nil {
		return p, err
	}
	return p, nil
}
