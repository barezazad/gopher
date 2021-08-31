package authmodel

// Auth model
type Auth struct {
	Username string `bind:"required" json:"username"`
	Password string `bind:"required" json:"password"`
}

// ResetPasswordRequest model
type ResetPasswordRequest struct {
	Email string `bind:"required" json:"email"`
}

// ResetPassword model
type ResetPassword struct {
	NewPassword string `bind:"required" json:"new_password"`
}
