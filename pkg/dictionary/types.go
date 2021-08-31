package dictionary

// R is used for parameters which we want to translated
type Translate string

// Lang enums
const (
	En = "en"
	Ku = "ku"
	Ar = "ar"
)

// Langs represents all accepted languages
var Langs = []string{
	En,
	Ku,
	Ar,
}

// Term is list of languages
type Term struct {
	En string `json:"en"`
	Ku string `json:"ku"`
	Ar string `json:"ar"`
}

// thisTerms used for holding language identifier as a string and Term Struct as value
var thisTerms map[string]Term
var translateInBackend bool
