package param

import (
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/pkg/generr"
	"strings"
)

func (p *Param) parseWhere(cols []string) (whereStr string, err error) {
	var whereArr []string
	var resultFilter string

	if resultFilter, err = p.parseFilter(cols); err != nil {
		return
	}

	if resultFilter != "" {
		whereArr = append(whereArr, resultFilter)
	}

	if p.PreCondition != "" {
		whereArr = append(whereArr, p.PreCondition)
	}

	if len(whereArr) > 0 {
		whereStr = strings.Join(whereArr[:], " AND ")
	}

	if p.ForceCondition != "" {
		if whereStr != "" {
			whereStr = " ( " + whereStr + " ) AND ( " + p.ForceCondition + " ) "
		} else {
			whereStr = p.ForceCondition
		}
	}

	return
}

// ParseWhere combine preConditions and filter with each other
func (p *Param) ParseWhere(engine *core.Engine, cols []string, code string) (whereStr string, err error) {

	whereStr, err = p.parseWhere(cols)
	if err != nil {
		err = engine.ErrorLog.TickCustom(err, code, generr.ValidationFailedErr, "", terms.ValidationFailed)
	}
	return
}

// ParseWhereDelete is used when the table has deleted_at column
func (p *Param) ParseWhereDelete(cols []string) (whereStr string, err error) {
	if whereStr, err = p.parseWhere(cols); err != nil {
		return
	}

	if !p.ShowDeletedRows {
		whereStr = " deleted_at is NULL AND " + whereStr
	} else {
		whereStr = " deleted_at is not NULL AND " + whereStr
	}

	return
}
