package response

import (
	"errors"
	"gopher/internal/core"
	"gopher/internal/param"
	"gopher/pkg/dictionary"
	"gopher/pkg/generr"
	"gopher/pkg/logparser"

	"github.com/gin-gonic/gin"
)

// Result is a standard output for success and failed requests
type Result struct {
	Message string                 `json:"message,omitempty"`
	Data    interface{}            `json:"data,omitempty"`
	Error   error                  `json:"error,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`

	// CustomError corerr.CustomError     `json:"custom_error,omitempty"`
}

// Response holding method related to response
type Response struct {
	Result  Result
	status  int
	Engine  *core.Engine
	Context *gin.Context
	abort   bool
	params  param.Param
	Entity  string
}

// New initiate the Response object
func New(engine *core.Engine, context *gin.Context, entity string) *Response {
	return &Response{
		Engine:  engine,
		Context: context,
		Entity:  entity,
	}
}

// NewParam initiate the Response object and params
func NewParam(engine *core.Engine, context *gin.Context, entity string) (*Response, param.Param) {
	params := param.Get(context, engine, entity)
	return &Response{
		Engine:  engine,
		Context: context,
		params:  params,
	}, params
}

// Error is used for add error to the result
func (r *Response) Error(err interface{}, data ...interface{}) *Response {
	if errCast, ok := err.(string); ok {
		r.Result.Error = errors.New(errCast)
	}
	if errCast, ok := err.(error); ok {
		r.Result.Error = errCast
	}
	r.Result.Data = data
	return r
}

// Status is used for add error to the result
func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}

// Message is used for add error to the result
func (r *Response) Message(message string) *Response {
	r.Result.Message = message
	return r
}

// MessageT get a message and params then translate it
func (r *Response) MessageT(message string, params ...interface{}) *Response {
	r.Result.Message = dictionary.Render(message, r.GetLanguage(), params...)
	return r
}

// Abort prepare response to abort instead of return in last step (JSON)
func (r *Response) Abort() *Response {
	r.abort = true
	return r
}

// JSON write output as json
func (r *Response) JSON(data ...interface{}) {
	var parsedError error
	if r.Result.Error != nil {
		r.Result.Error = logparser.AddPath(r.Result.Error, r.Context.Request.RequestURI)

		customError := logparser.GetCustom(r.Result.Error)
		lang := r.GetLanguage()
		r.Result.Error = logparser.ApplyCustom(r.Result.Error, generr.UniqErrorMap[customError])

		//tra := translator(lang)
		r.status, parsedError = logparser.Parse(r.Result.Error, lang)
	}

	// if data is one element don't put it in array
	var finalData interface{}
	if data != nil {
		finalData = data
		if len(data) == 1 {
			finalData = data[0]
		}
	}

	if r.abort {
		r.Context.AbortWithStatusJSON(r.status, &Result{
			Message: r.Result.Message,
			Error:   parsedError,
			Data:    finalData,
		})
	} else {
		r.Context.JSON(r.status, &Result{
			Message: r.Result.Message,
			Error:   parsedError,
			Data:    finalData,
			// CustomError: r.Result.Error,
		})
	}
}

func (r *Response) GetLanguage() string {
	var langLevel string

	langLevel = r.Engine.Environments.DefaultLanguage

	lang, ok := r.Context.Get("LANGUAGE")
	if ok {
		langLevel = lang.(string)
	}

	switch langLevel {
	case dictionary.En:
		return dictionary.En
	case dictionary.Ku:
		return dictionary.Ku
	case dictionary.Ar:
		return dictionary.Ar
	}

	return dictionary.En
}
