package param

import "gopher/pkg/helper/sqluri"

// parseFilter call parser for convert urlQuery to SQL query
func (p *Param) parseFilter(cols []string) (result string, err error) {
	if p.Filter == "" {
		return
	}

	if result, err = sqluri.Parser(p.Filter, cols); err != nil {
		return
	}

	return
}
