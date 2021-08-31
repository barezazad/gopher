package model

import (
	"time"

	"gorm.io/gorm"
)

// Common columns for all models
type Common struct {
	ID        uint           `gorm:"primary_key,unique" json:"id,omitempty" bind:"update:required"`
	CreatedAt *time.Time     `gorm:"<-:create;type:timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
	CreatedBy uint           `gorm:"<-:create;" json:"created_by,omitempty" bind:"create:required"`
	UpdatedBy uint           `gorm:"<-:update;" json:"updated_by,omitempty" bind:"update:required"`
}

// BasicModel columns for all models
type BasicModel struct {
	ID        uint           `gorm:"primary_key,unique" json:"id,omitempty" bind:"update:required"`
	CreatedAt *time.Time     `gorm:"<-:create;type:timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}
