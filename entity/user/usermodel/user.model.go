package usermodel

import (
	"gopher/internal/model"
)

// User model
type User struct {
	model.Common
	RoleID       uint     `gorm:"index:role_id_idx" json:"role_id" bind:"required"`
	Name         string   `gorm:"not null" json:"name,omitempty" bind:"required,min=5"`
	Username     string   `gorm:"not null;unique" json:"username,omitempty" bind:"required"`
	Password     string   `gorm:"not null" json:"password,omitempty" bind:"required,min=8,max=20"`
	Lang         string   `gorm:"type:varchar(2);default:'en'" json:"lang,omitempty" bind:"oneof=languagetype"`
	Email        string   `gorm:"not null" json:"email,omitempty" bind:"required,email"`
	Phone        float64  `json:"phone,omitempty"`
	Status       string   `gorm:"default:'active';type:enum('active','inactive','terminate')" json:"status,omitempty" bind:"oneof=userstatus"`
	Token        string   `gorm:"-" json:"token,omitempty" table:"-"`
	StrResources string   `gorm:"->" json:"str_resources,omitempty" table:"roles.resources as str_resources"`
	Resources    []string `gorm:"-" json:"resources,omitempty" table:"-"`
	Role         string   `gorm:"->" json:"role,omitempty" table:"roles.name as role"`
	OldPassword  string   `gorm:"-"  json:"old_password,omitempty" table:"-"`
}

const (
	// Table is used inside the repo layer
	Table = "users"
)

func (p *User) TableName() string {
	return Table
}
