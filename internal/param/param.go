package param

import (
	"fmt"
	"gopher/internal/core"

	"gorm.io/gorm"
)

// Param for describing request's parameter
type Param struct {
	Pagination
	ID              uint     // it use for send id to service, repo ...
	Search          string   // for send string, ex: find user by username
	Filter          string   // in where condition to filter list
	PreCondition    string   // in where to set condition in api or service
	ForceCondition  string   // it will be force cond, ex (PreCondition and Filter ) And (ForceCondition)
	UserID          uint     // it incloud userid
	Username        string   // it incloud username
	Lang            string   // it incloud lang
	ShowDeletedRows bool     // to show deleted rows in table, default false
	Tx              *gorm.DB // to handle transaction for repo
}

// Pagination is a struct, contains the fields which affected the front-end pagination
type Pagination struct {
	Select string
	Order  string
	Limit  int
	Offset int
}

// New return an initiate of the param with default limit
func New() Param {
	var param Param
	param.Limit = core.DefaultLimit
	param.ShowDeletedRows = core.ShowDeletedRows
	param.Order = "id"

	return param
}

// IsElementDeleted is used for checking delete an element
func IsElementDeleted(table string, col string, id interface{}) Param {
	var param Param
	param.Limit = 1
	param.Select = "*"
	param.Order = fmt.Sprintf("%v.id asc", table)
	param.PreCondition = fmt.Sprintf("%v.%v = %v", table, col, id)
	param.PreCondition += fmt.Sprintf(" AND %v.deleted_at IS NULL ", table)

	return param
}

func (p *Param) GetDB(db *gorm.DB) *gorm.DB {
	if p.Tx != nil {
		return p.Tx
	}
	return db
}
