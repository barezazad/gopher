package citymodel

import (
	"gopher/entity/document/documentmodel"
	"gopher/internal/model"
	"mime/multipart"
)

// City model
type City struct {
	model.Common
	Name        string                   `gorm:"not null;unique" bind:"required" json:"name,omitempty" form:"name"`
	Description string                   `json:"description,omitempty" form:"description"`
	Attachments []*multipart.FileHeader  `gorm:"-" table:"-" json:"attachments" form:"attachments"`
	Documents   []documentmodel.Document `gorm:"-" table:"-" json:"documents" `
}

const (
	// Table is used inside the repo layer
	Table = "cities"
)

func (City) TableName() string {
	return Table
}
