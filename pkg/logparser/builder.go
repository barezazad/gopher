package logparser

import (
	"errors"
	"gopher/pkg/dictionary"
)

//LimbErr is main type for controlling the package
type LimbErr struct {
	err error
}

// New return an initiate of the LimbErr
func New(errStr string, code ...string) *LimbErr {
	var limbErr LimbErr
	err := errors.New(errStr)

	if len(code) > 0 {
		limbErr.err = AddCode(err, code[0])
	} else {
		limbErr.err = err
	}

	return &limbErr

}

// Take initiate the
func Take(err error, code ...string) *LimbErr {
	var limbErr LimbErr
	if len(code) > 0 {
		limbErr.err = AddCode(err, code[0])
	} else {
		limbErr.err = err
	}

	return &limbErr
}

// Code assign a code to the error
func (p *LimbErr) Code(code string) *LimbErr {
	p.err = AddCode(p.err, code)
	return p
}

// Message append a message to the error
func (p *LimbErr) Message(message string, params ...interface{}) *LimbErr {
	p.err = AddMessage(p.err, message, getPhrase(params...)...)
	return p
}

// Custom is used when some value like status code and basic data needs to be appended to the error
func (p *LimbErr) Custom(custom CustomError) *LimbErr {
	p.err = SetCustom(p.err, custom)
	return p
}

// Domain for specify the position of the error in the projcet
func (p *LimbErr) Domain(domain string) *LimbErr {
	p.err = AddDomain(p.err, domain)
	return p
}

// Path demonstrate ApiURL of the request
func (p *LimbErr) Path(path string) *LimbErr {
	p.err = AddPath(p.err, path)
	return p
}

// Status can be used in HTTP Status code
func (p *LimbErr) Status(status int) *LimbErr {
	p.err = AddStatus(p.err, status)
	return p
}

// InvalidParam is used when want to pint to a field which caused the error
func (p *LimbErr) InvalidParam(field, reason string, params ...interface{}) *LimbErr {
	p.err = AddInvalidParam(p.err, field, reason, params...)
	return p
}

// Build return an initiate of the struct
func (p *LimbErr) Build() error {
	return p.err
}

func getPhrase(phrases ...interface{}) (result []interface{}) {
	// check is phrase will be translate or not
	for _, v := range phrases {
		switch i := v.(type) {
		case string:
			result = append(result, dictionary.Translate(i))
		default:
			result = append(result, i)
		}
	}

	return
}
