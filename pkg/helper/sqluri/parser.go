package sqluri

import (
	"fmt"
	"gopher/pkg/helper"
	"regexp"
	"strings"

	"gopher/pkg/logparser"
)

var swap = map[string]string{
	"[eq]":   " = ",
	"[ne]":   " != ",
	"[gt]":   " > ",
	"[lt]":   " < ",
	"[gte]":  " >= ",
	"[lte]":  " <= ",
	"[like]": " LIKE ",
	"[and]":  " AND ",
	"[or]":   " OR ",
}

var ops = []string{"eq", "ne", "gt", "lt", "gte", "lte", "like"}

//Parser will break the filter into the sub-query
func Parser(str string, cols []string) (string, error) {
	regCol := regexp.MustCompile(`\w+[.\w+]*`)
	arr := regCol.FindAllString(str, -1)

	regAfterDot := regexp.MustCompile(`\w+$`)
	var reducedCols []string
	for _, v := range cols {
		if strings.Contains(v, ".") {
			reducedCols = append(reducedCols, regAfterDot.FindString(v))
		}
		if strings.Contains(v, "as") {
			splitString := strings.Split(v, " ")
			reducedCols = append(reducedCols, splitString[0])
		}
	}

	cols = append(cols, reducedCols...)

	if len(arr) == 0 {
		return "", fmt.Errorf("filter is not valid")
	}

	pre := arr[0]
	for _, v := range arr {
		if ok, _ := helper.Includes(ops, v); ok {
			if ok, err := helper.Includes(cols, pre); !ok || err != nil {
				if err != nil {
					return "", err
				}
				err := fmt.Errorf("col '%s' not exist", pre)
				err = logparser.AddInvalidParam(err, pre,
					"column %v not not exist", pre)
				return "", err
			}
		}
		pre = v
	}

	for k, v := range swap {
		str = strings.Replace(str, k, v, -1)
	}

	return str, nil
}
