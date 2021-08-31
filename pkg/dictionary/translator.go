package dictionary

import (
	"fmt"
)

// SafeTranslate doesn't add !!! around word in case of not exist for translate
func SafeTranslate(str string, lang string, params ...interface{}) (string, bool) {
	if !translateInBackend {
		return str, true
	}

	term, ok := thisTerms[str]
	if ok {
		var pattern string

		switch lang {
		case En:
			pattern = term.En
		case Ku:
			pattern = term.Ku
		case Ar:
			pattern = term.Ar
		default:
			pattern = str
		}

		// if type of param is dictionary.Translate then translate it
		for i, v := range params {
			switch v.(type) {
			case Translate:
				term := v.(Translate)
				params[i] = Render(string(term), lang)
			}
		}

		if params != nil {
			if !(params[0] == nil || params[0] == "") {
				pattern = fmt.Sprintf(pattern, params...)
			}
		}

		return pattern, true

	}

	return "", false

}

// Render the requested term
func Render(str string, lang string, params ...interface{}) string {
	if !translateInBackend {
		return str
	}

	pattern, ok := SafeTranslate(str, lang, params...)
	if ok {
		return pattern
	}

	return str
}
