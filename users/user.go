package users

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User represents a user
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Username  string     `gorm:"type:varchar(100);unique_index;not null" json:"username"`
	Password  string     `json:"password"`
}

// GetSingle fetches a user from the db
func (u *User) GetSingle(db *gorm.DB) error {
	return db.First(&u).Error
}

// Update updates a user in the db
func (u *User) Update(db *gorm.DB) error {
	return db.Save(&u).Error
}

// Delete deletes a user
func (u *User) Delete(db *gorm.DB) error {
	return db.Delete(&u).Error
}

// Create inserts a user in the db
func (u *User) Create(db *gorm.DB) error {
	return db.Create(&u).Error
}

// Get retrieves users from the db
func Get(db *gorm.DB, offset, limit int) ([]User, error) {
	var u []User
	if err := db.Offset(offset).Limit(limit).Find(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}
