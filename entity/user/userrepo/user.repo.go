package userrepo

import (
	"gopher/entity/user/usermodel"
	"gopher/internal/core"
	"gopher/internal/core/dberror"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/helper"
	"reflect"

	"gorm.io/gorm/clause"
)

// UserRepo for injecting engine
type UserRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideUserRepo is used in wire and initiate the Cols
func ProvideUserRepo(engine *core.Engine) UserRepo {
	return UserRepo{
		Engine: engine,
		Cols:   helper.TagExtractor(reflect.TypeOf(usermodel.User{}), usermodel.Table),
	}
}

// FindByID finds the user via its id
func (p *UserRepo) FindByID(params param.Param) (user usermodel.User, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000149"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Select(colsStr).
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? ", params.ID).First(&user).Error

	err = dberror.DbError(p.Engine, err, "E1000150", user, usermodel.Table, terms.Info)
	return
}

// FindByIDTx finds the user via its id and lock row (for update)
func (p *UserRepo) FindByIDTx(params param.Param) (user usermodel.User, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000151"); err != nil {
		return
	}

	// it use for share
	// Clauses(clause.Locking{Strength: "SHARE",Table: clause.Table{Name: clause.CurrentTable},})
	err = params.GetDB(p.Engine.DB).Clauses(clause.Locking{Strength: "UPDATE"}).
		Select(colsStr).
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where("users.id = ? ", params.ID).First(&user).Error

	err = dberror.DbError(p.Engine, err, "E1000152", user, usermodel.Table, terms.Info)

	return
}

// FindAll it return list of all user
func (p *UserRepo) FindAll(params param.Param) (users []usermodel.User, err error) {

	err = params.GetDB(p.Engine.DB).Find(&users).Error

	for i := range users {
		users[i].Password = ""
	}

	err = dberror.DbError(p.Engine, err, "E1000153", usermodel.User{}, usermodel.Table, terms.Info)
	return
}

// FindByUsername finds the user via its username
func (p *UserRepo) FindByUsername(params param.Param) (user usermodel.User, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000151"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).Select(colsStr).
		Joins("inner join roles on roles.id = users.role_id").
		Where("users.username = ?", params.Search).
		First(&user).Error

	err = dberror.DbError(p.Engine, err, "E1000154", user, usermodel.Table, terms.Info)
	return
}

// FindByEmail finds the user via its email
func (p *UserRepo) FindByEmail(params param.Param) (user usermodel.User, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("users.email = ?", params.Search).
		First(&user).Error

	err = dberror.DbError(p.Engine, err, "E1000155", user, usermodel.Table, terms.Info)
	return
}

// List returns an array of users
func (p *UserRepo) List(params param.Param) (users []usermodel.User, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000156"); err != nil {
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000157"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Select(colsStr).
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&users).Error

	err = dberror.DbError(p.Engine, err, "E1000158", usermodel.User{}, usermodel.Table, terms.List)

	for i := range users {
		users[i].Password = ""
	}

	return
}

// Count of users, mainly calls with List
func (p *UserRepo) Count(params param.Param) (count int64, err error) {

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000159"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).Table(usermodel.Table).
		Joins("INNER JOIN roles ON roles.id = users.role_id").
		Where(whereStr).
		Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000160", usermodel.User{}, usermodel.Table, terms.List)
	return
}

// Create a user
func (p *UserRepo) Create(user usermodel.User, params param.Param) (u usermodel.User, err error) {

	if err = params.GetDB(p.Engine.DB).Create(&user).Scan(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000161", user, usermodel.Table, terms.Created)
	}
	return
}

// Save the user, in case it is not exist create it
func (p *UserRepo) Save(user usermodel.User, params param.Param) (u usermodel.User, err error) {

	if err = params.GetDB(p.Engine.DB).Save(&user).Find(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000162", user, usermodel.Table, terms.Saved)
		return
	}

	return
}

// Delete the user
func (p *UserRepo) Delete(user usermodel.User, params param.Param) (err error) {

	if err = params.GetDB(p.Engine.DB).Delete(&user).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000163", user, usermodel.Table, terms.Deleted)
	}
	return
}
