package service

import (
	"gopher/entity/user/usermodel"
	"gopher/entity/user/userrepo"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/generr"
	"gopher/pkg/helper"
	"gopher/pkg/helper/password"
)

type UserService struct {
	Repo   userrepo.UserRepo
	Engine *core.Engine
}

// ProvideUserService for user is used in wire
func ProvideUserService(p userrepo.UserRepo) UserService {
	return UserService{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting user by it's id
func (p *UserService) FindByID(params param.Param) (user usermodel.User, err error) {

	if user, err = p.Repo.FindByID(params); err != nil {
		return
	}
	return
}

// FindAll it return list of all user
func (p *UserService) FindAll(params param.Param) (users []usermodel.User, err error) {

	if users, err = p.Repo.FindAll(params); err != nil {
		return
	}
	return
}

// FindByUsername find user with username, used for auth
func (p *UserService) FindByUsername(params param.Param) (user usermodel.User, err error) {

	if user, err = p.Repo.FindByUsername(params); err != nil {
		return
	}
	return
}

// FindByEmail find user with email, used for auth
func (p *UserService) FindByEmail(params param.Param) (user usermodel.User, err error) {

	if user, err = p.Repo.FindByEmail(params); err != nil {
		err = p.Engine.ErrorLog.TickCustom(err, "E1000119", generr.NotFoundErr, "", terms.ProvidedEmailIsWrong)
		return
	}
	return
}

// List of users, it support pagination and search and return back count
func (p *UserService) List(params param.Param) (users []usermodel.User, count int64, err error) {

	if users, err = p.Repo.List(params); err != nil {
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		return
	}
	return
}

// Create a user
func (p *UserService) Create(user usermodel.User, params param.Param) (createdUser usermodel.User, err error) {

	// validate to create user
	if err = validator.BindTagExtractor(p.Engine, user, "E1000120", usermodel.Table, core.Create); err != nil {
		return
	}

	// hashing password
	if user.Password, err = password.Hash(user.Password, p.Engine.Environments.JWT.PasswordSalt); err != nil {
		err = p.Engine.ErrorLog.TickBadRequest(err, "E1000121", usermodel.Table, core.Create, user)
		p.Engine.ServerLog.CheckError(err, "E1000122", terms.HashingPasswordFailed, user)
	}

	if createdUser, err = p.Repo.Create(user, params); err != nil {
		return
	}

	// set empty for password
	createdUser.Password = ""

	return
}

// Save user
func (p *UserService) Save(user usermodel.User, params param.Param) (updatedUser, userBefore usermodel.User, err error) {

	// validate to update user
	if err = validator.BindTagExtractor(p.Engine, user, "E1000123", usermodel.Table, core.Save); err != nil {
		return
	}

	if userBefore, err = p.FindByID(params); err != nil {
		return
	}

	if user.Password != "" {
		if user.Password, err = password.Hash(user.Password, p.Engine.Environments.JWT.PasswordSalt); err != nil {
			err = p.Engine.ErrorLog.TickBadRequest(err, "E1000124", usermodel.Table, core.Save, user)
			p.Engine.ServerLog.CheckError(err, "E1000125", terms.HashingPasswordFailed, user)
		}
	} else {
		user.Password = userBefore.Password
	}

	if updatedUser, err = p.Repo.Save(user, params); err != nil {
		return
	}

	// set empty for password
	updatedUser.Password = ""

	// to delete resource in cache
	p.Engine.Cache.Delete(helper.UintToStr(user.ID))

	return
}

// Delete user
func (p *UserService) Delete(params param.Param) (user usermodel.User, err error) {
	if user, err = p.FindByID(params); err != nil {
		return
	}

	if err = p.Repo.Delete(user, params); err != nil {
		return
	}

	return
}
