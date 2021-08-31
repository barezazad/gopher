package activitymodel

import "gopher/internal/model"

// Activity model
type Activity struct {
	model.Common
	Event    string `gorm:"index:event_idx" json:"event"`
	UserID   uint   `json:"user_id"`
	Username string `gorm:"index:username_idx" json:"username"`
	IP       string `json:"ip"`
	URI      string `gorm:"type:text" json:"uri"`
	Before   string `gorm:"type:text" json:"before"`
	After    string `gorm:"type:text" json:"after"`
}

const (
	// Table is used inside the repo layer
	Table = "activities"
)

func (Activity) TableName() string {
	return Table
}
