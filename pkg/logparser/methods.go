package logparser

import (
	"fmt"
)

// WithCode is used for carrying the code of error
type WithCode struct {
	Err  error
	Code string
}

func (p *WithCode) Error() string { return fmt.Sprint(p.Err) }

// AddCode add custom code to the error, it is useful for tracing an error
// if error pass from different method or function you can gave it multiple error
// and in case error happened easily we can trace which function's passed that error
// similar to panic
func AddCode(err error, code string) error {
	return &WithCode{
		Err:  fmt.Errorf("#%v, %w", code, err),
		Code: code,
	}
}

// WithMessage keeps the message of the error, each error can have one message
type WithMessage struct {
	Err    error
	Msg    string
	Params []interface{}
}

func (p *WithMessage) Error() string { return fmt.Sprint(p.Err) }

// AddMessage add custom message to error, params can be used inside the translator function
func AddMessage(err error, msg string, params ...interface{}) error {
	return &WithMessage{
		Err:    err,
		Msg:    msg,
		Params: params,
	}
}

// WithType is add type and title to the error
type WithType struct {
	Err   error
	Type  string
	Title string
}

func (p *WithType) Error() string { return fmt.Sprint(p.Err) }

// AddType used for adding type to the error
func AddType(err error, errType string, title string) error {
	return &WithType{
		Err:   err,
		Type:  errType,
		Title: title,
	}
}

// WithPath attach path to the error
type WithPath struct {
	Err  error
	Path string
}

func (p *WithPath) Error() string { return fmt.Sprint(p.Err) }

// AddPath is used for adding path to the error, useful for REST API
func AddPath(err error, path string) error {
	return &WithPath{
		Err:  err,
		Path: path,
	}
}

// WithStatus attach status to the error
type WithStatus struct {
	Err    error
	Status int
}

func (p *WithStatus) Error() string { return fmt.Sprint(p.Err) }

// AddStatus can be used for adding HTTP status code like 404 or etc
func AddStatus(err error, status int) error {
	return &WithStatus{
		Err:    err,
		Status: status,
	}
}

// WithDomain attach domain to the error
type WithDomain struct {
	Err    error
	Domain string
}

func (p *WithDomain) Error() string { return fmt.Sprint(p.Err) }

// AddDomain separate different domains of application. In case you don't need it just ignore it
func AddDomain(err error, domain string) error {
	return &WithDomain{
		Err:    err,
		Domain: domain,
	}
}

// WithInvalidParam holds invalid parameters
type WithInvalidParam struct {
	Err    error
	Field  string
	Reason string
	Params []interface{}
}

func (p *WithInvalidParam) Error() string { return fmt.Sprint(p.Err) }

// AddInvalidParam is used for specify a field that has an error
func AddInvalidParam(err error, field, reason string, params ...interface{}) error {
	var gErr error
	if err == nil {
		gErr = fmt.Errorf(fmt.Sprintf(reason, params...))
	} else {
		gErr = err
	}

	return &WithInvalidParam{
		Err:    gErr,
		Field:  field,
		Reason: reason,
		Params: params,
	}
}

// WithCustom is used for holding the uniqError for filling the type and title based on local
// customization
type WithCustom struct {
	Err    error
	Custom CustomError
}

func (p *WithCustom) Error() string { return fmt.Sprint(p.Err) }

// SetCustom is used for adding a custom error to the error and it reduce size of the code
func SetCustom(err error, custom CustomError) error {
	return &WithCustom{
		Err:    err,
		Custom: custom,
	}
}
