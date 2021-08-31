package rolemodel

import (
	"gopher/internal/model"
)

// Role model
type Role struct {
	model.Common
	Name        string `gorm:"not null;unique" bind:"required" json:"name,omitempty"`
	Resources   string `gorm:"type:text" bind:"required" json:"resources,omitempty"`
	Description string `json:"description,omitempty"`
}

const (
	// Table is used inside the repo layer
	Table = "roles"
)

func (Role) TableName() string {
	return Table
}
