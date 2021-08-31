package helper

import (
	"fmt"
	"reflect"
	"regexp"
)

type extractor struct {
	columns []string
	regex   *regexp.Regexp
	table   string
}

// TagExtractor extract the name of table and field from json and table tag
func TagExtractor(t reflect.Type, table string) []string {
	ext := extractor{
		columns: []string{},
		regex:   regexp.MustCompile(`\w+`),
		table:   table,
	}

	ext.getTag(t)

	return ext.columns
}

func (p *extractor) getTag(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		externalTable := field.Tag.Get("table")

		if field.Type.Kind() == reflect.Struct && field.Type.Name() == "Model" {
			p.columns = append(p.columns, p.table+"."+"id")
			p.columns = append(p.columns, p.table+"."+"created_at")
			p.columns = append(p.columns, p.table+"."+"updated_at")
			p.columns = append(p.columns, p.table+"."+"deleted_at")
		} else {
			// below code is for recursive
			if field.Type.Kind() == reflect.Struct && externalTable == "" {
				p.getTag(field.Type)
				continue
			}
		}

		column := field.Tag.Get("json")
		if column == "" {
			continue
		}
		column = p.regex.FindString(column)

		switch {
		case externalTable == "-":
			continue
		case externalTable != "":
			column = externalTable
		default:
			column = fmt.Sprintf("%v.%v", p.table, column)
		}

		p.columns = append(p.columns, column)
	}

}
