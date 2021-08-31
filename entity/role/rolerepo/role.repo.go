package rolerepo

import (
	"gopher/entity/role/rolemodel"
	"gopher/internal/core"
	"gopher/internal/core/dberror"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/helper"
	"reflect"

	"gorm.io/gorm/clause"
)

// RoleRepo for injecting engine
type RoleRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideRoleRepo is used in wire and initiate the Cols
func ProvideRoleRepo(engine *core.Engine) RoleRepo {
	return RoleRepo{
		Engine: engine,
		Cols:   helper.TagExtractor(reflect.TypeOf(rolemodel.Role{}), rolemodel.Table),
	}
}

// FindByID finds the role via its id
func (p *RoleRepo) FindByID(params param.Param) (role rolemodel.Role, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("roles.id = ? ", params.ID).First(&role).Error

	err = dberror.DbError(p.Engine, err, "E1000068", role, rolemodel.Table, terms.Info)
	return
}

// FindByIDTx finds the role via its id and lock row (for update)
func (p *RoleRepo) FindByIDTx(params param.Param) (role rolemodel.Role, err error) {

	err = params.GetDB(p.Engine.DB).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("roles.id = ? ", params.ID).First(&role).Error

	err = dberror.DbError(p.Engine, err, "E1000069", role, rolemodel.Table, terms.Info)
	return
}

// FindAll it return list of all role
func (p *RoleRepo) FindAll(params param.Param) (roles []rolemodel.Role, err error) {

	err = params.GetDB(p.Engine.DB).Find(&roles).Error

	err = dberror.DbError(p.Engine, err, "E1000070", rolemodel.Role{}, rolemodel.Table, terms.Info)
	return
}

// List returns an array of roles
func (p *RoleRepo) List(params param.Param) (roles []rolemodel.Role, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000071"); err != nil {
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000072"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&roles).Error

	err = dberror.DbError(p.Engine, err, "E1000073", rolemodel.Role{}, rolemodel.Table, terms.List)

	return
}

// Count of roles, mainly calls with List
func (p *RoleRepo) Count(params param.Param) (count int64, err error) {

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000074"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).Table(rolemodel.Table).
		Where(whereStr).
		Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000075", rolemodel.Role{}, rolemodel.Table, terms.List)
	return
}

// Create a role
func (p *RoleRepo) Create(role rolemodel.Role, params param.Param) (u rolemodel.Role, err error) {

	if err = params.GetDB(p.Engine.DB).Create(&role).Scan(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000076", role, rolemodel.Table, terms.Created)
	}
	return
}

// Save the role, in case it is not exist create it
func (p *RoleRepo) Save(role rolemodel.Role, params param.Param) (u rolemodel.Role, err error) {

	if err = params.GetDB(p.Engine.DB).Save(&role).Find(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000077", role, rolemodel.Table, terms.Saved)
		return
	}

	return
}

// Delete the role
func (p *RoleRepo) Delete(role rolemodel.Role, params param.Param) (err error) {

	if err = params.GetDB(p.Engine.DB).Unscoped().Delete(&role).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000078", role, rolemodel.Table, terms.Deleted)
	}
	return
}
