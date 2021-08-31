package authapi

import (
	"gopher/entity/auth/authenum"
	"gopher/entity/auth/authmodel"
	"gopher/entity/service"
	"gopher/entity/user/userenum"
	"gopher/entity/user/usermodel"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"gopher/pkg/dictionary"
	"gopher/pkg/helper/email"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthAPI for injecting auth service
type AuthAPI struct {
	Service service.AuthService
	Engine  *core.Engine
}

// ProvideAuthAPI for auth used in wire
func ProvideAuthAPI(p service.AuthService) AuthAPI {
	return AuthAPI{Service: p, Engine: p.Engine}
}

// Login auth
func (p *AuthAPI) Login(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, authenum.Login)
	var auth authmodel.Auth

	if err := resp.Bind(&auth, "E1000010"); err != nil {
		return
	}

	user, err := p.Service.Login(auth, params)
	if err != nil {
		resp.Error(err).JSON()
		resp.Record(authenum.LoginFailed, auth.Username, len(auth.Password))
		return
	}

	tmpUser := user
	tmpUser.Token = ""

	for _, v := range strings.Split(user.StrResources, ",") {
		user.Resources = append(user.Resources, strings.TrimSpace(v))
	}
	if user.Resources != nil {
		user.StrResources = ""
	}

	resp.RecordCreate(authenum.Login, user)
	resp.Status(http.StatusOK).
		Message(terms.UserLoggedInSuccessfully).
		JSON(user)
}

// Profile returns the user's information
func (p *AuthAPI) Profile(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, authenum.ViewProfile)
	var user usermodel.User
	var err error

	params.ID = params.UserID
	if user, err = p.Service.Profile(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	user.Password = ""

	resp.Record(authenum.ViewProfile, user)
	resp.Status(http.StatusOK).
		MessageT(authenum.ViewProfile).
		JSON(user)
}

// UpdateProfile to update profile
func (p *AuthAPI) UpdateProfile(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, authenum.UpdateProfile)
	var profile, profileBefore, profileUpdated usermodel.User
	var err error

	if err = resp.Bind(&profile, "E1000011"); err != nil {
		return
	}

	params.ID = params.UserID
	profile.ID = params.UserID
	profile.UpdatedBy = params.UserID
	if profileBefore, profileUpdated, err = p.Service.UpdateProfile(profile, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(userenum.UpdateUser, profileBefore, profileUpdated)
	resp.Status(http.StatusOK).
		MessageT(terms.VUpdatedSuccessfully, dictionary.Translate(terms.Profile)).
		JSON(profileUpdated)
}

// Logout will erase the resources from Cache
func (p *AuthAPI) Logout(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, authenum.ViewProfile)

	p.Service.Logout(params)

	resp.Record(authenum.Logout, params.Username)
	resp.Status(http.StatusOK).
		Message(authenum.Logout).
		JSON()
}

// RequestResetPassword to request reset password by user
func (p *AuthAPI) RequestResetPassword(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, authenum.ResetPasswordRequest)
	var err error
	var request authmodel.ResetPasswordRequest
	var user usermodel.User

	if err := resp.Bind(&request, "E1000012"); err != nil {
		return
	}

	params.Search = request.Email
	// find user by email
	if user, err = p.Service.FindByEmail(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	// generate token
	params.UserID = user.ID
	params.Username = user.Username
	token, err := p.Service.ResetPasswordToken(1, params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	// send email link
	link := p.Engine.Environments.Setting.ResetPasswordUrl + "?tkn=" + token
	email.ResetPasswordNotification(p.Engine, link, request.Email)
	user.Password = ""

	resp.Record(authenum.ResetPasswordRequest, nil, user)
	resp.Status(http.StatusOK).
		MessageT(terms.RequestToResetPasswordSuccessfull).
		JSON("")

}

// ResetPassword for user
func (p *AuthAPI) ResetPassword(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, authenum.ResetPassword)
	var err error
	var resetPassword authmodel.ResetPassword
	var oldUser, updatedUser usermodel.User
	var token string

	if err := resp.Bind(&resetPassword, "E1000013"); err != nil {
		return
	}

	if tkn, ok := c.Get("TOKEN"); ok {
		token = tkn.(string)
	}

	params.ID = params.UserID
	if oldUser, updatedUser, err = p.Service.ResetPassword(resetPassword, token, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(authenum.ResetPassword, oldUser, updatedUser)
	resp.Status(http.StatusOK).
		MessageT(terms.ResetPasswordSuccessfully).
		JSON(updatedUser)
}

// UpdateLang to update language
func (p *AuthAPI) UpdateLang(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, authenum.Updatelang)
	var profile, profileBefore, profileUpdated usermodel.User
	var err error

	if err = resp.Bind(&profile, "E1000014"); err != nil {
		return
	}

	params.ID = params.UserID
	profile.ID = params.UserID
	profile.UpdatedBy = params.UserID
	if profileBefore, profileUpdated, err = p.Service.UpdateLang(profile, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(userenum.UpdateUser, profileBefore, profileUpdated)
	resp.Status(http.StatusOK).
		MessageT(terms.VUpdatedSuccessfully, dictionary.Translate(terms.Language)).
		JSON(profileUpdated)
}
