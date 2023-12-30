package user

import (
	"github.com/google/uuid"
	"github.com/ranta0/rest-and-go/domain/role"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	PublicID   string      `json:"id"              form:"id"       validate:""         gorm:"unique;not null"`
	Username   string      `json:"username"        form:"username" validate:"required" gorm:"unique_index;not null"`
	Password   string      `json:"-"               form:"-"        validate:""         gorm:"not null"`
	Name       string      `json:"name"            form:"name"     validate:""`
	Age        int         `json:"age"             form:"age"      validate:""`
	Status     string      `json:"status"          form:"status"   validate:""`
	Roles      []role.Role `json:"roles,omitempty" form:"roles"    validate:""         gorm:"many2many:user_roles;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.PublicID = uuid.New().String()
	return nil
}
