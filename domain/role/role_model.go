package role

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model `json:"-"`
	PublicID   string `json:"id"   form:"id"   validate:""         gorm:"unique;not null"`
	Name       string `json:"name" form:"name" validate:"required" gorm:"unique_index;not null"`
}

func (u *Role) BeforeCreate(tx *gorm.DB) error {
	u.PublicID = uuid.New().String()
	return nil
}
