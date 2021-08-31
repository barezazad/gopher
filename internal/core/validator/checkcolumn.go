package validator

import (
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/pkg/generr"
	"gopher/pkg/helper"
	"gopher/pkg/logparser"
	"strings"
)

// CheckColumns will check columns for security
func CheckColumns(engine *core.Engine, cols []string, requestedCols, code string) (string, error) {
	var err error

	if requestedCols == "*" || requestedCols == "" {
		return strings.Join(cols, ","), nil
	}

	variates := strings.Split(requestedCols, ",")
	for _, v := range variates {
		if ok, _ := helper.Includes(cols, v); !ok {
			err = logparser.AddInvalidParam(err, v, terms.VisNotValid, v)
			err = engine.ErrorLog.TickCustom(err, code, generr.ValidationFailedErr, "", terms.ValidationFailed)
		}
	}

	return requestedCols, err

}
