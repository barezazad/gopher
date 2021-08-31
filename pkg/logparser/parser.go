package logparser

import (
	"errors"
	"gopher/pkg/dictionary"
	"net/http"
)

// Parse convert chained error to the Final format for send in JSON format
func Parse(err error, lang string) (int, error) {
	var final Final
	var status int

	for err != nil {
		switch e := err.(type) {
		case interface{ Unwrap() error }:
			err = errors.Unwrap(err)
		case *WithMessage:
			if final.Message == "" {
				final.Message = dictionary.Render(e.Msg, lang, e.Params...)
			}
			err = e.Err
		case *WithCode:
			if final.Code == "" {
				final.Code = e.Code
			}
			err = e.Err
		case *WithType:
			final.Type = e.Type
			final.Title = dictionary.Render(e.Title, lang)
			err = e.Err
		case *WithPath:
			final.Path += appendText(final.Path, e.Path)
			err = e.Err
		case *WithStatus:
			final.Status = e.Status
			status = e.Status
			err = e.Err
		case *WithDomain:
			final.Domain = e.Domain
			err = e.Err
		case *WithInvalidParam:
			field := Field{
				Field:        e.Field,
				Reason:       dictionary.Render(e.Reason, lang, e.Params...),
				ReasonParams: e.Params,
			}
			final.InvalidParams = append(final.InvalidParams, field)
			err = e.Err
		case *WithCustom:
			err = e.Err
		case error:
			final.OriginalError += e.Error()
			err = errors.Unwrap(err)
		default:
			return http.StatusInternalServerError, &final
		}
	}
	return status, &final
}

// GetCustom extract custom error from error's interface
func GetCustom(err error) (customError CustomError) {
	for err != nil {
		switch e := err.(type) {
		case interface{ Unwrap() error }:
			err = errors.Unwrap(err)
		case *WithCustom:
			return e.Custom
		case error:
			if errCast, ok := e.(*WithMessage); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*WithCode); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*WithType); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*WithPath); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*WithStatus); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*WithDomain); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*WithInvalidParam); ok {
				err = errCast.Err
				continue
			}
			return
		default:
			return
		}
	}

	return
}

// ApplyCustom add custom errors to the error's interface
func ApplyCustom(err error, theme ErrorTheme) error {
	err = AddStatus(err, theme.Status)
	return err

}

func appendText(str string, txt string) (result string) {
	if str == "" {
		result = txt
	} else {
		result = str + ", " + txt
	}
	return
}
