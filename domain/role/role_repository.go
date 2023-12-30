package role

import (
	"gorm.io/gorm"
)

// RoleRepository handles database operations for the Role entity using gorm.
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository with the provided gorm.DB instance.
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// Create adds a new role to the database.
func (ur *RoleRepository) Create(role *Role) error {
	return ur.db.Create(role).Error
}

// GetByID retrieves a role from the database by ID.
func (ur *RoleRepository) GetByID(id string) (*Role, error) {
	var role Role
	query := ur.db.Where("public_id = ?", id)
	err := query.First(&role).Error

	return &role, err
}

// Update updates a role in the database.
func (ur *RoleRepository) Update(id string, role *Role) error {
	return ur.db.Model(&Role{}).Where("public_id = ?", id).Updates(&role).Error
}

// Delete deletes a role from the database.
func (ur *RoleRepository) Delete(id string) error {
	return ur.db.Where("public_id = ?", id).Delete(&Role{}).Error
}

// GetAll retrieves all roles from the database
func (ur *RoleRepository) GetAll(limit, offset int) ([]Role, error) {
	var roles []Role
	query := ur.db.Model(&Role{}).Limit(limit).Offset(offset)
	err := query.Find(&roles).Error

	return roles, err
}

// CountAll get the coutn of all values from the database
func (ur *RoleRepository) CountAll() (int, error) {
	var count int64

	err := ur.db.Model(&Role{}).Count(&count).Error
	return int(count), err
}
