package giftmodel

import (
	"gopher/entity/document/documentmodel"
	"gopher/internal/model"
)

// Gift model
type Gift struct {
	model.Common
	Name        string                   `gorm:"not null;unique" bind:"required" json:"name,omitempty" form:"name"`
	Attachments []string                 `gorm:"-" table:"-" json:"attachments" form:"attachments"`
	Documents   []documentmodel.Document `gorm:"-" table:"-" json:"documents" `
}

const (
	// Table is used inside the repo layer
	Table = "gifts"
)

func (Gift) TableName() string {
	return Table
}
