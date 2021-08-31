package documentmodel

import (
	"gopher/internal/model"
)

// Document model
type Document struct {
	model.Common
	Name     string `gorm:"not null" bind:"required" json:"name,omitempty"`
	FileType string `bind:"required" json:"file_type,omitempty"`
	DocId    uint   `gorm:"not null" json:"doc_id,omitempty"`
	DocType  string `gorm:"type:enum('cities','gifts')" bind:"required,oneof=documenttypes" json:"doc_type,omitempty"`
}

const (
	// Table is used inside the repo layer
	Table = "documents"
)

func (Document) TableName() string {
	return Table
}
