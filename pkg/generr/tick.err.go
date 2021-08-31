package generr

import (
	"gopher/internal/core/terms"
	"gopher/pkg/logparser"
)

// Tick is combining AddCode and LogError in services to reduce the code
func (log *Log) Tick(err error, code string, message string, data ...interface{}) error {
	if code != "" {
		err = logparser.AddCode(err, code)
	}
	log.LogError(err, message, data...)
	return err
}

// TickCustom is combining AddCode and LogError in services to reduce the code and specify the
// error type
func (log *Log) TickCustom(err error, code string, custom logparser.CustomError,
	data interface{}, message string, params ...interface{}) error {
	// err = logparser.SetCustom(err, custom)
	// return Tick(log, err, code, message, data)

	limErr := logparser.Take(err, code).Custom(custom).
		Message(message, params...).Build()

	log.LogError(err, limErr.Error(), data)
	return limErr
}

//TickCustomMessage for custom message , get same status error,just change message
func TickCustomMessage(err error, code string, message string, params ...interface{}) error {
	limErr := logparser.Take(err, code).Message(message, params...).Build()
	return limErr
}

// TickValidate is automatically add validation error custom to the error
func (log *Log) TickValidate(err error, code string,
	entity, action string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(ValidationFailedErr).
		Message(terms.ValidationForVFailedInV, entity, action).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickBadRequest is automatically add badrequest error custom to the error
func (log *Log) TickBadRequest(err error, code string,
	entity string, action string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(BadRequestErr).
		Message(terms.BadRequestForVInV, entity, action).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickRecordNotFound is automatically add record not found error custom to the error
func (log *Log) TickRecordNotFound(err error, code string, table interface{}, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(NotFoundErr).
		Message(terms.RecordNotFoundForInVTable, table).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickBinding is automatically add binding error custom to the error
func (log *Log) TickBinding(err error, code string,
	entity string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(BindingErr).
		Message(terms.ErrorInBindingV, entity).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickUnauthorized is automatically add unauthorized error custom to the error
func (log *Log) TickUnauthorized(err error, code string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(UnauthorizedErr).
		Message(terms.Unauthorized).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickRouteNotFound is automatically add RouteNotFound error custom to the error
func (log *Log) TickRouteNotFound(err error, code string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(RouteNotFoundErr).
		Message(terms.RouteNotFound).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickForeign is automatically add forgin key error custom to the error
func (log *Log) TickForeign(err error, code string,
	value interface{}, action string, data ...interface{}) error {

	var limErr error
	switch action {
	case terms.Created:
		limErr = logparser.Take(err, code).Custom(ForeignErr).
			Message(terms.CheckThisVField, value).Build()
	default:
		limErr = logparser.Take(err, code).Custom(ForeignErr).
			Message(terms.ItHasRelationWithSomeElementSoItIsNotV, action).Build()
	}

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickDuplicate is automatically add duplicate error custom to the error
func (log *Log) TickDuplicate(err error, code string,
	value interface{}, entity string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(DuplicateErr).
		Message(terms.VisAlreadyExistInTableV, value, entity).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickInternalServer is automatically add InternalServer error custom to the error
func (log *Log) TickInternalServer(err error, code string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(InternalServerErr).
		Message(terms.InternalServerError).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}

// TickForbidden is automatically add Forbidden error custom to the error
func (log *Log) TickForbidden(err error, code string, data ...interface{}) error {

	limErr := logparser.Take(err, code).Custom(ForbiddenErr).
		Message(terms.Forbidden).Build()

	log.LogError(err, limErr.Error(), data...)
	return limErr
}
