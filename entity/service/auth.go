package service

import (
	"fmt"
	"gopher/entity/auth/authenum"
	"gopher/entity/auth/authmodel"
	"gopher/entity/user/userenum"
	"gopher/entity/user/usermodel"
	"gopher/entity/user/userrepo"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/model"
	"gopher/internal/param"
	"gopher/pkg/generr"
	"gopher/pkg/helper/password"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// AuthService defining auth service
type AuthService struct {
	Engine *core.Engine
}

// ProvideAuthService for auth is used in wire
func ProvideAuthService(engine *core.Engine) AuthService {
	return AuthService{
		Engine: engine,
	}
}

// Login User
func (p *AuthService) Login(auth authmodel.Auth, params param.Param) (user usermodel.User, err error) {

	// validate to login model
	if err = validator.BindTagExtractor(p.Engine, auth, "E1000085", authenum.Login, core.Login); err != nil {
		return
	}

	jwtKey := p.Engine.Environments.JWT.SecretKey
	passwordSalt := p.Engine.Environments.JWT.PasswordSalt
	jwtExpiration := p.Engine.Environments.JWT.Expiration

	params.Search = auth.Username
	userServ := ProvideUserService(userrepo.ProvideUserRepo(p.Engine))
	if user, err = userServ.FindByUsername(params); err != nil {
		err = p.Engine.ErrorLog.TickCustom(err, "E1000086", generr.UnauthorizedErr, user, terms.UsernameOrPasswordIsWrong)
		return
	}

	if password.Verify(auth.Password, user.Password, passwordSalt) {

		now := time.Now().In(p.Engine.TZ)

		expirationTime := now.Add(time.Duration(jwtExpiration) * time.Second)
		claims := &model.JWTClaims{
			Username: auth.Username,
			ID:       user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		user.Password = ""
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		if user.Token, err = token.SignedString([]byte(jwtKey)); err != nil {
			err = p.Engine.ErrorLog.TickInternalServer(err, "E1000087", user)
			p.Engine.ServerLog.CheckError(err, "E1000088", terms.CanNotInGenerateToken, user)
			return
		}

		// clear user info in cache
		authCacheKey := fmt.Sprintf("%v-%v", terms.Auth, user.ID)
		p.Engine.Cache.Delete(authCacheKey)

	} else {
		err = p.Engine.ErrorLog.TickCustom(err, "E1000089", generr.UnauthorizedErr,
			auth, terms.UsernameOrPasswordIsWrong)
	}

	return
}

// Logout erase resources from the cache
func (p *AuthService) Logout(params param.Param) {
	authCacheKey := fmt.Sprintf("%v-%v", terms.Auth, params.UserID)
	p.Engine.Cache.Delete(authCacheKey)
}

// Profile return user's information
func (p *AuthService) Profile(params param.Param) (user usermodel.User, err error) {
	userServ := ProvideUserService(userrepo.ProvideUserRepo(p.Engine))

	if user, err = userServ.FindByID(params); err != nil {
		return
	}
	return
}

// TemporaryToken generate instant token
func (p *AuthService) TemporaryToken(params param.Param) (tmpKey string, err error) {
	jwtKey := p.Engine.Environments.JWT.SecretKey
	now := time.Now().In(p.Engine.TZ)

	expirationTime := now.Add(core.TemporaryTokenDuration * time.Second)
	claims := &model.JWTClaims{
		Username: params.Username,
		ID:       params.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tmpKey, err = token.SignedString([]byte(jwtKey)); err != nil {
		p.Engine.ServerLog.CheckError(err, "E1000090", terms.CanNotInGenerateToken, claims)
		return
	}

	return
}

// ResetPasswordToken generate instant token for reset password
func (p *AuthService) ResetPasswordToken(hour int, params param.Param) (tmpKey string, err error) {

	now := time.Now().In(p.Engine.TZ)
	jwtKey := p.Engine.Environments.JWT.SecretKey

	expirationTime := now.Add(time.Duration(hour) * time.Hour)
	claims := &model.JWTClaims{
		Username: params.Username,
		ID:       params.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tmpKey, err = token.SignedString([]byte(jwtKey)); err != nil {
		err = p.Engine.ErrorLog.TickInternalServer(err, "E1000091")
		p.Engine.ServerLog.CheckError(err, "E1000092", terms.CanNotInGenerateToken)
		return
	}

	resetPasswordKey := fmt.Sprintf("%v-%v", terms.ResetPassword, params.UserID)
	hourToSec := hour * 3600
	p.Engine.Cache.Set(resetPasswordKey, tmpKey, hourToSec)
	return
}

// UpdateProfile for user
func (p *AuthService) UpdateProfile(profile usermodel.User, params param.Param) (oldProfile, updatedProfile usermodel.User, err error) {

	var inProfile usermodel.User
	// find user by id
	if oldProfile, err = p.Profile(params); err != nil {
		return
	}

	// check if have change to password
	passwordSalt := p.Engine.Environments.JWT.PasswordSalt
	if profile.Password != "" {
		if !password.Verify(profile.OldPassword, oldProfile.Password, passwordSalt) {
			err = p.Engine.ErrorLog.TickCustom(err, "E1000093", generr.BadRequestErr, "", terms.OldPasswordIsWrong)
			return
		}

		// validate to update profile
		if err = validator.BindTagExtractor(p.Engine, profile, "E1000094", usermodel.Table, core.Save); err != nil {
			return
		}

		if profile.Password, err = password.Hash(profile.Password, passwordSalt); err != nil {
			err = p.Engine.ErrorLog.TickBadRequest(err, "E1000095", userenum.UpdateUser, core.Update)
			return
		}
	} else {
		profile.Password = oldProfile.Password
	}

	oldProfile.Password = ""
	inProfile = oldProfile
	inProfile.Password = profile.Password
	inProfile.Lang = profile.Lang
	inProfile.Email = profile.Email
	inProfile.Name = profile.Name
	inProfile.Phone = profile.Phone

	userServ := ProvideUserService(userrepo.ProvideUserRepo(p.Engine))
	if updatedProfile, err = userServ.Repo.Save(inProfile, params); err != nil {
		return
	}
	updatedProfile.Password = ""

	return
}

// ResetPassword user
func (p *AuthService) ResetPassword(resetpassword authmodel.ResetPassword, token string, params param.Param) (oldUser, updatedUser usermodel.User, err error) {

	var user usermodel.User

	// find user by id
	if oldUser, err = p.Profile(params); err != nil {
		return
	}
	oldUser.Password = ""

	// check in cache to know this token is exist
	resetPasswordKey := fmt.Sprintf("%v-%v", terms.ResetPassword, params.UserID)
	cacheToken, _ := p.Engine.Cache.Get(resetPasswordKey)
	if !strings.Contains(cacheToken, token) {
		err = p.Engine.ErrorLog.TickCustom(err, "E1000096", generr.BadRequestErr, "", terms.UserDataInvalidTryToSendRequestAgain)
		return
	}

	user = oldUser
	user.Password = resetpassword.NewPassword

	userServ := ProvideUserService(userrepo.ProvideUserRepo(p.Engine))
	if updatedUser, _, err = userServ.Save(user, params); err != nil {
		return
	}

	// remove token in cache
	p.Engine.Cache.Delete(resetPasswordKey)
	return
}

// find user by email
func (p *AuthService) FindByEmail(params param.Param) (user usermodel.User, err error) {

	userServ := ProvideUserService(userrepo.ProvideUserRepo(p.Engine))
	if user, err = userServ.FindByEmail(params); err != nil {
		return
	}
	return
}

// UpdateLang for user
func (p *AuthService) UpdateLang(profile usermodel.User, params param.Param) (oldProfile, updatedProfile usermodel.User, err error) {

	var inProfile usermodel.User
	// find user by id
	if oldProfile, err = p.Profile(params); err != nil {
		return
	}

	inProfile = oldProfile
	inProfile.Lang = profile.Lang

	userServ := ProvideUserService(userrepo.ProvideUserRepo(p.Engine))
	if updatedProfile, err = userServ.Repo.Save(inProfile, params); err != nil {
		return
	}
	updatedProfile.Password = ""

	return
}
