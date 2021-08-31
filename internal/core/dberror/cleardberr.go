package dberror

import (
	"errors"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/pkg/generr"
	"gopher/pkg/helper"
	"gopher/pkg/logparser"
	"strings"

	"gorm.io/gorm"
)

// dbError is an internal method for generate proper database error
func DbError(engine *core.Engine, err error, code string, data interface{}, entity, action string) error {

	switch {
	// nil
	case err == nil:
		err = nil

		// record not found
	case errors.Is(err, gorm.ErrRecordNotFound):
		err = engine.ErrorLog.TickRecordNotFound(err, code, entity, data)

		// duplicate error
	case strings.Contains(strings.ToUpper(err.Error()), "DUPLICATE"):

		field := helper.RegexFindBetweenTwoPattern(`for key`, ``, err.Error())
		value := helper.RegexFindBetweenTwoPattern(`Duplicate entry`, `for key`, err.Error())

		err = logparser.AddInvalidParam(err, field, terms.VisAlreadyExistInTableV, value, entity)
		err = engine.ErrorLog.TickCustom(err, code, generr.DuplicateErr, data, terms.VisUniqueInTableV, field, entity)

		// unknown column error
	case strings.Contains(strings.ToUpper(err.Error()), "UNKNOWN COLUMN"):
		err = engine.ErrorLog.TickValidate(err, code, entity, action, data)

		// foregin key error
	case strings.Contains(strings.ToUpper(err.Error()), "FOREIGN"):
		field := helper.RegexFindBetweenTwoPattern(`FOREIGN KEY \(`, `\) REFERENCES`, err.Error())
		err = engine.ErrorLog.TickForeign(err, code, field, action, data)

		// default
	default:
		err = engine.ErrorLog.TickBadRequest(err, code, entity, action, data)
	}

	return err
}
