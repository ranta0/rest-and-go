package domain

import (
	"gorm.io/gorm"
)

// All models implement the gorm function to set PublicId
type ModelInterface interface {
	BeforeCreate(tx *gorm.DB)
}
