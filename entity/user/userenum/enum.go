package userenum

const (
	Entity     = "users"
	CreateUser = "user-create"
	UpdateUser = "user-update"
	DeleteUser = "user-delete"
	ListUser   = "user-list"
	ViewUser   = "user-view"
	AllUser    = "user-all"
	ExcelUser  = "user-excel"
)

const (
	Active    = "active"
	Inactive  = "inactive"
	Terminate = "terminate"
)

// user status enum
var UserStatus = []string{
	Active,
	Inactive,
	Terminate,
}
