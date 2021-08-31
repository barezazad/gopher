package settingmodel

import "gopher/internal/model"

// Setting model
type Setting struct {
	model.Common
	Property    string `gorm:"not null;unique" bind:"required" json:"property,omitempty"`
	Value       string `gorm:"type:text" bind:"required" json:"value,omitempty"`
	Type        string `bind:"required" json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

const (
	// Table is used inside the repo layer
	Table = "settings"
)

func (Setting) TableName() string {
	return Table
}
