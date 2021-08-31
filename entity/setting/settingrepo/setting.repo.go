package settingrepo

import (
	"gopher/entity/setting/settingmodel"
	"gopher/internal/core"
	"gopher/internal/core/dberror"
	"gopher/internal/core/terms"
	"gopher/internal/core/validator"
	"gopher/internal/param"
	"gopher/pkg/helper"
	"reflect"
)

// SettingRepo for injecting engine
type SettingRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideSettingRepo is used in wire and initiate the Cols
func ProvideSettingRepo(engine *core.Engine) SettingRepo {
	return SettingRepo{
		Engine: engine,
		Cols:   helper.TagExtractor(reflect.TypeOf(settingmodel.Setting{}), settingmodel.Table),
	}
}

// FindByID finds the setting via its id
func (p *SettingRepo) FindByID(params param.Param) (setting settingmodel.Setting, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("settings.id = ?", params.ID).Find(&setting).Error

	err = dberror.DbError(p.Engine, err, "E1000131", setting, settingmodel.Table, terms.Info)
	return
}

// FindByProperty finds the setting via its property
func (p *SettingRepo) FindByProperty(params param.Param) (setting settingmodel.Setting, err error) {

	err = params.GetDB(p.Engine.DB).
		Where("settings.property = ? ", params.Search).First(&setting).Error

	err = dberror.DbError(p.Engine, err, "E1000132", setting, settingmodel.Table, terms.Info)
	return
}

// FindAll it return list of all setting
func (p *SettingRepo) FindAll(params param.Param) (settings []settingmodel.Setting, err error) {

	err = params.GetDB(p.Engine.DB).Find(&settings).Error

	err = dberror.DbError(p.Engine, err, "E1000133", settingmodel.Setting{}, settingmodel.Table, terms.Info)
	return
}

// List returns an array of settings
func (p *SettingRepo) List(params param.Param) (settings []settingmodel.Setting, err error) {

	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Engine, p.Cols, params.Select, "E1000134"); err != nil {
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000135"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&settings).Error

	err = dberror.DbError(p.Engine, err, "E1000136", settingmodel.Setting{}, settingmodel.Table, terms.List)

	return
}

// Count of settings, mainly calls with List
func (p *SettingRepo) Count(params param.Param) (count int64, err error) {

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Engine, p.Cols, "E1000137"); err != nil {
		return
	}

	err = params.GetDB(p.Engine.DB).
		Table(settingmodel.Table).
		Where(whereStr).
		Count(&count).Error

	err = dberror.DbError(p.Engine, err, "E1000138", settingmodel.Setting{}, settingmodel.Table, terms.List)
	return
}

// Create a setting
func (p *SettingRepo) Create(setting settingmodel.Setting, params param.Param) (u settingmodel.Setting, err error) {

	if err = params.GetDB(p.Engine.DB).Create(&setting).Scan(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000139", setting, settingmodel.Table, terms.Created)
	}
	return
}

// Save the setting, in case it is not exist create it
func (p *SettingRepo) Save(setting settingmodel.Setting, params param.Param) (u settingmodel.Setting, err error) {

	if err = params.GetDB(p.Engine.DB).Save(&setting).Find(&u).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000140", setting, settingmodel.Table, terms.Saved)
		return
	}

	return
}

// Delete the setting
func (p *SettingRepo) Delete(setting settingmodel.Setting, params param.Param) (err error) {

	if err = params.GetDB(p.Engine.DB).Unscoped().Delete(&setting).Error; err != nil {
		err = dberror.DbError(p.Engine, err, "E1000141", setting, settingmodel.Table, terms.Deleted)
	}
	return
}
