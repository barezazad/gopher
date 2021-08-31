package generr

import (
	"gopher/internal/core/terms"
	"gopher/pkg/logparser"
	"net/http"
)

// CustomError model
const (
	Nil logparser.CustomError = iota
	UnknownErr
	UnauthorizedErr
	NotFoundErr
	RouteNotFoundErr
	ValidationFailedErr
	ForeignErr
	DuplicateErr
	InternalServerErr
	BindingErr
	ForbiddenErr
	PreDataInsertedErr //428
	BadRequestErr
)

// UniqErrorMap is used for categorized errors and connect error with error page also primary fill
// the status code and entity and title
var UniqErrorMap logparser.CustomErrorMap

func init() {
	UniqErrorMap = make(map[logparser.CustomError]logparser.ErrorTheme)

	UniqErrorMap[UnauthorizedErr] = logparser.ErrorTheme{
		Type:   "#Unauthorized",
		Title:  terms.Unauthorized,
		Status: http.StatusUnauthorized,
	}

	UniqErrorMap[ValidationFailedErr] = logparser.ErrorTheme{
		Type:   "#VALIDATION_FAILED",
		Title:  terms.ValidationFailed,
		Status: http.StatusUnprocessableEntity,
	}

	UniqErrorMap[NotFoundErr] = logparser.ErrorTheme{
		Type:   "#NOT_FOUND",
		Title:  terms.RecordNotFound,
		Status: http.StatusNotFound,
	}

	UniqErrorMap[RouteNotFoundErr] = logparser.ErrorTheme{
		Type:   "#NOT_FOUND",
		Title:  terms.RouteNotFound,
		Status: http.StatusNotFound,
	}

	UniqErrorMap[ForeignErr] = logparser.ErrorTheme{
		Type:   "#FOREIGN_KEY",
		Title:  terms.ForeignKeyError,
		Status: http.StatusConflict,
	}

	UniqErrorMap[InternalServerErr] = logparser.ErrorTheme{
		Type:   "#INTERNAL_SERVER_ERROR",
		Title:  terms.InternalServerError,
		Status: http.StatusInternalServerError,
	}

	UniqErrorMap[DuplicateErr] = logparser.ErrorTheme{
		Type:   "#DUPLICATE_ERROR",
		Title:  terms.DuplicateHappened,
		Status: http.StatusConflict,
	}

	UniqErrorMap[BindingErr] = logparser.ErrorTheme{
		Type:   "#NOT_BIND",
		Title:  terms.BindFailed,
		Status: http.StatusUnprocessableEntity,
	}

	UniqErrorMap[ForbiddenErr] = logparser.ErrorTheme{
		Type:   "#FORBIDDEN",
		Title:  terms.Forbidden,
		Status: http.StatusForbidden,
	}

	UniqErrorMap[BadRequestErr] = logparser.ErrorTheme{
		Type:   "#BAD_REQUEST",
		Title:  terms.BadRequest,
		Status: http.StatusBadRequest,
	}
}
