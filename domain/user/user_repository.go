package user

import (
	"gorm.io/gorm"
)

// UserRepository handles database operations for the User entity using gorm.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository with the provided gorm.DB instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create adds a new user to the database.
func (ur *UserRepository) Create(user *User) error {
	return ur.db.Create(user).Error
}

// GetByID retrieves a user from the database by ID.
func (ur *UserRepository) GetByID(id string) (*User, error) {
	var user User
	err := ur.db.Where("public_id = ?", id).First(&user).Error

	return &user, err
}

// Update updates a user in the database.
func (ur *UserRepository) Update(id string, user *User) error {
	return ur.db.Model(&User{}).Where("public_id = ?", id).Updates(&user).Error
}

// Delete deletes a user from the database.
func (ur *UserRepository) Delete(id string) error {
	return ur.db.Where("public_id = ?", id).Delete(&User{}).Error
}

// GetAll retrieves all users from the database
func (ur *UserRepository) GetAll(limit, offset int) ([]User, error) {
	var users []User
	err := ur.db.
		Model(&User{}).
		Limit(limit).
		Offset(offset).
		Find(&users).
		Error

	return users, err
}

// CountAll get the coutn of all values from the database
func (ur *UserRepository) CountAll() (int, error) {
	var count int64

	err := ur.db.Model(&User{}).Count(&count).Error
	return int(count), err
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	var user User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
