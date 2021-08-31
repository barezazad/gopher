package core

import (
	"gopher/entity/document/documentenum"
	"gopher/entity/user/userenum"
	"gopher/pkg/dictionary"
)

// Action enums
const (
	Create = "create"
	Delete = "delete"
	Update = "update"
	Save   = "save"
	Login  = "login"
	Active = "active"
)

var MustBeInTypes = map[string][]string{
	"languagetype":  dictionary.Langs,
	"userstatus":    userenum.UserStatus,
	"documenttypes": documentenum.DocumentTypes,
}
