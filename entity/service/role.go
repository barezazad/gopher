package service

import (
	"encoding/json"
	"gopher/entity/role/rolemodel"
	"gopher/entity/role/rolerepo"
	"gopher/internal/core"
	"gopher/internal/core/validator"
	"gopher/internal/param"
)

type RoleService struct {
	Repo   rolerepo.RoleRepo
	Engine *core.Engine
}

// ProvideRoleService for role is used in wire
func ProvideRoleService(p rolerepo.RoleRepo) RoleService {
	return RoleService{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting role by it's id
func (p *RoleService) FindByID(params param.Param) (role rolemodel.Role, err error) {

	if role, err = p.Repo.FindByID(params); err != nil {
		return
	}
	return
}

// FindAll it return list of all role
func (p *RoleService) FindAll(params param.Param) (roles []rolemodel.Role, err error) {

	if cacheRoles, ok := p.getInCache(); ok {
		roles = cacheRoles
		return
	}

	if roles, err = p.Repo.FindAll(params); err != nil {
		return
	}

	// save all roles in cache
	p.Engine.Cache.Set(rolemodel.Table, roles, p.Engine.Environments.Redis.CacheApiTTL)
	return
}

// List of roles, it support pagination and search and return back count
func (p *RoleService) List(params param.Param) (roles []rolemodel.Role, count int64, err error) {

	if roles, err = p.Repo.List(params); err != nil {
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		return
	}
	return
}

// Create a role
func (p *RoleService) Create(role rolemodel.Role, params param.Param) (createdRole rolemodel.Role, err error) {

	// validate to create role
	if err = validator.BindTagExtractor(p.Engine, role, "E1000115", rolemodel.Table, core.Create); err != nil {
		return
	}

	if createdRole, err = p.Repo.Create(role, params); err != nil {
		return
	}

	// remove roles in cache
	p.Engine.Cache.Delete(rolemodel.Table)

	return
}

// Save role
func (p *RoleService) Save(role rolemodel.Role, params param.Param) (updatedRole, roleBefore rolemodel.Role, err error) {

	// validate to update role
	if err = validator.BindTagExtractor(p.Engine, role, "E1000116", rolemodel.Table, core.Save); err != nil {
		return
	}

	if roleBefore, err = p.FindByID(params); err != nil {
		return
	}

	if updatedRole, err = p.Repo.Save(role, params); err != nil {
		return
	}

	// remove roles in cache
	p.Engine.Cache.Delete(rolemodel.Table)

	return
}

// Delete role
func (p *RoleService) Delete(params param.Param) (role rolemodel.Role, err error) {

	if role, err = p.FindByID(params); err != nil {
		return
	}

	if err = p.Repo.Delete(role, params); err != nil {
		return
	}

	// remove roles in cache
	p.Engine.Cache.Delete(rolemodel.Table)

	return
}

// check value in cache and get it , if exist
func (p *RoleService) getInCache() (roles []rolemodel.Role, ok bool) {

	result, err := p.Engine.Cache.Get(rolemodel.Table)
	if err != nil {
		return
	}

	if result != "" {
		if err = json.Unmarshal([]byte(result), &roles); err != nil {
			return
		}
		ok = true
	}

	return
}
