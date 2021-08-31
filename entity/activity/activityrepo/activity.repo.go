package activityrepo

import (
	"gopher/entity/activity/activitymodel"
	"gopher/internal/core"
	"gopher/internal/core/dberror"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/helper"
	"reflect"
)

// ActivityRepo for injecting engine
type ActivityRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideActivityRepo is used in wire and initiate the Cols
func ProvideActivityRepo(engine *core.Engine) ActivityRepo {
	return ActivityRepo{
		Engine: engine,
		Cols:   helper.TagExtractor(reflect.TypeOf(activitymodel.Activity{}), activitymodel.Table),
	}
}

// FindByID finds the activity via its id
func (p *ActivityRepo) FindByID(params param.Param) (activity activitymodel.Activity, err error) {

	err = params.GetDB(p.Engine.ActivityDB).
		Where("activities.id = ? ", params.ID).Find(&activity).Error

	err = dberror.DbError(p.Engine, err, "E1000001", activity, activitymodel.Table, terms.Info)
	return
}

// List returns an array of activities
func (p *ActivityRepo) List(params param.Param) (activities []activitymodel.Activity, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000002"); err != nil {
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000003"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.ActivityDB).
		Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&activities).Error

	err = dberror.DbError(p.Engine, err, "E1000004", activitymodel.Activity{}, activitymodel.Table, terms.List)

	return
}

// Count of activities, mainly calls with List
func (p *ActivityRepo) Count(params param.Param) (count int64, err error) {

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000005"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.ActivityDB).
		Table(activitymodel.Table).
		Where(whereStr).
		Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000006", activitymodel.Activity{}, activitymodel.Table, terms.List)
	return
}

// Create a activity
func (p *ActivityRepo) Create(activity activitymodel.Activity, params param.Param) (u activitymodel.Activity, err error) {

	if err = params.GetDB(p.Engine.ActivityDB).Create(&activity).Scan(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000007", activity, activitymodel.Table, terms.Created)
	}
	return
}

// CreateBatch a activity
func (p *ActivityRepo) CreateBatch(activities []activitymodel.Activity, params param.Param) (u activitymodel.Activity, err error) {

	err = params.GetDB(p.Engine.ActivityDB).Create(&activities).Error
	return
}

// Save the activity, in case it is not exist create it
func (p *ActivityRepo) Save(activity activitymodel.Activity, params param.Param) (u activitymodel.Activity, err error) {

	if err = params.GetDB(p.Engine.ActivityDB).Save(&activity).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000008", activity, activitymodel.Table, terms.Saved)
		return
	}

	params.GetDB(p.Engine.ActivityDB).Where("id = ?", activity.ID).Find(&u)
	return
}

// Delete the activity
func (p *ActivityRepo) Delete(activity activitymodel.Activity, params param.Param) (err error) {

	if err = params.GetDB(p.Engine.ActivityDB).Unscoped().Delete(&activity).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000009", activity, activitymodel.Table, terms.Deleted)
	}
	return
}
