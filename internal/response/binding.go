package response

import (
	"encoding/json"
	"gopher/internal/core/rsa"
	"gopher/internal/core/terms"
	"gopher/internal/model"

	"gopher/pkg/generr"
	"gopher/pkg/helper"
	"gopher/pkg/logparser"
)

// NotBind use special custom_error for reduced it
func (r *Response) NotBind(err error, code string) {
	err = logparser.Take(err, code).
		Message(terms.ErrorInBindingV, err.Error()).
		Custom(generr.BindingErr).Build()

	r.Error(err).JSON()
}

// Bind is used to make it more easear for binding items
func (r *Response) Bind(st interface{}, code string) (err error) {
	if err = r.Context.ShouldBindJSON(&st); err != nil {
		r.NotBind(err, code)
		return
	}

	return
}

// GetID returns the ID
func (r *Response) GetID(idIn, code string) (id uint, err error) {
	tmpID, err := helper.StrToUint64(idIn)
	if err != nil {
		err = logparser.Take(err, code).
			Message(terms.InvalidV, "ID").
			Custom(generr.ValidationFailedErr).Build()
		r.Error(err).JSON()
	}
	id = uint(tmpID)
	return
}

// BindErrorCipher use special custom_error for reduced it
func (r *Response) BindErrorCipher(err error, code string) {
	err = logparser.Take(err, code).
		Message(terms.ErrorInBinding).
		Custom(generr.BindingErr).Build()
	r.Error(err).JSON()
}

// BindCipher is used to make it more easear for binding cipher data and decode it
func (r *Response) BindCipher(st interface{}, code string) (err error) {

	// define cipher model
	var cipher model.Encryption

	// bind cipher data
	if err = r.Context.ShouldBindJSON(&cipher); err != nil {
		r.BindErrorCipher(err, code)
		return
	}

	// decrypt cipher data
	var decryptPayload []byte
	if decryptPayload, err = rsa.Decrypt(r.Engine, cipher.Data); err != nil {
		r.BindErrorCipher(err, code)
		return
	}

	// parse decrypted data
	if err = json.Unmarshal(decryptPayload, &st); err != nil {
		r.BindErrorCipher(err, code)
		return
	}

	return
}
