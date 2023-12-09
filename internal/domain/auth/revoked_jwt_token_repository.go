package auth

import (
	"gorm.io/gorm"
)

type RevokedJWTTokenRepository struct {
	db *gorm.DB
}

func NewRevokedJWTTokenRepository(db *gorm.DB) *RevokedJWTTokenRepository {
	return &RevokedJWTTokenRepository{db: db}
}

func (r *RevokedJWTTokenRepository) RevokeToken(token string) error {
	return r.db.Create(&RevokedJWTToken{Token: token}).Error
}

// IsTokenRevoked checks if a token is in the database, indicating it has been revoked
func (r *RevokedJWTTokenRepository) IsTokenRevoked(token string) bool {
	var revokedToken RevokedJWTToken
	err := r.db.Where("token = ?", token).First(&revokedToken).Error
	return err == nil
}
