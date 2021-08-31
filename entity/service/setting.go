package service

import (
	"gopher/entity/setting/settingmodel"
	"gopher/entity/setting/settingrepo"
	"gopher/internal/core"
	"gopher/internal/core/validator"
	"gopher/internal/param"
)

type SettingService struct {
	Repo   settingrepo.SettingRepo
	Engine *core.Engine
}

// ProvideSettingService for setting is used in wire
func ProvideSettingService(p settingrepo.SettingRepo) SettingService {
	return SettingService{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting setting by it's id
func (p *SettingService) FindByID(params param.Param) (setting settingmodel.Setting, err error) {

	if setting, err = p.Repo.FindByID(params); err != nil {
		return
	}
	return
}

// FindByProperty finds the setting via its property
func (p *SettingService) FindByProperty(params param.Param) (setting settingmodel.Setting, err error) {

	if setting, err = p.Repo.FindByProperty(params); err != nil {
		return
	}
	return
}

// FindAll it return list of all setting
func (p *SettingService) FindAll(params param.Param) (settings []settingmodel.Setting, err error) {

	if settings, err = p.Repo.FindAll(params); err != nil {
		return
	}
	return
}

// List of settings, it support pagination and search and return back count
func (p *SettingService) List(params param.Param) (settings []settingmodel.Setting, count int64, err error) {

	if settings, err = p.Repo.List(params); err != nil {
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		return
	}
	return
}

// Create a setting
func (p *SettingService) Create(setting settingmodel.Setting, params param.Param) (createdSetting settingmodel.Setting, err error) {

	// validate to create setting
	if err = validator.BindTagExtractor(p.Engine, setting, "E1000117", settingmodel.Table, core.Create); err != nil {
		return
	}

	if createdSetting, err = p.Repo.Create(setting, params); err != nil {
		return
	}

	return
}

// Save setting
func (p *SettingService) Save(setting settingmodel.Setting, params param.Param) (updatedSetting, settingBefore settingmodel.Setting, err error) {

	// validate to update setting
	if err = validator.BindTagExtractor(p.Engine, setting, "E1000118", settingmodel.Table, core.Save); err != nil {
		return
	}

	if settingBefore, err = p.FindByID(params); err != nil {
		return
	}

	if updatedSetting, err = p.Repo.Save(setting, params); err != nil {
		return
	}

	return
}

// Delete setting
func (p *SettingService) Delete(params param.Param) (setting settingmodel.Setting, err error) {

	if setting, err = p.FindByID(params); err != nil {
		return
	}

	if err = p.Repo.Delete(setting, params); err != nil {
		return
	}

	return
}
