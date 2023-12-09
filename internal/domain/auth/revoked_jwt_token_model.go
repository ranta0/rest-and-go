package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RevokedToken represents a revoked token in the database
type RevokedJWTToken struct {
	gorm.Model `json:"-"`
	PublicID   string `json:"id"    form:"id"    validate:"" gorm:"unique;not null"`
	Token      string `json:"token" form:"token" validate:"" gorm:"unique_index"`
}

func (u *RevokedJWTToken) BeforeCreate(tx *gorm.DB) error {
	u.PublicID = uuid.New().String()
	return nil
}
