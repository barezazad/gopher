package settingapi

import (
	"gopher/entity/service"
	"gopher/entity/setting/settingenum"
	"gopher/entity/setting/settingmodel"
	"gopher/internal/core"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"net/http"

	"gopher/pkg/dictionary"

	"github.com/gin-gonic/gin"
)

// SettingAPI for injecting setting service
type SettingAPI struct {
	Service service.SettingService
	Engine  *core.Engine
}

// ProvideSettingAPI for setting is used in wire
func ProvideSettingAPI(c service.SettingService) SettingAPI {
	return SettingAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a setting by it's id
func (p *SettingAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, settingmodel.Table)
	var setting settingmodel.Setting
	var err error

	if params.ID, err = resp.GetID(c.Param("settingID"), "E1000126"); err != nil {
		return
	}

	if setting, err = p.Service.FindByID(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(settingenum.ViewSetting)
	resp.Status(http.StatusOK).
		MessageT(terms.VInfo, terms.Setting).
		JSON(setting)
}

// FindAll it return list of all setting
func (p *SettingAPI) FindAll(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, settingmodel.Table)
	var settings []settingmodel.Setting
	var err error

	if settings, err = p.Service.FindAll(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(settingenum.AllSetting)
	resp.Status(http.StatusOK).
		MessageT(terms.AllV, terms.Settings).
		JSON(settings)
}

// List of settings
func (p *SettingAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, settingmodel.Table)
	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(settingenum.ListSetting)
	resp.Status(http.StatusOK).
		MessageT(terms.ListOfV, terms.Settings).
		JSON(data)
}

// Create setting
func (p *SettingAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, settingmodel.Table)
	var setting, createdSetting settingmodel.Setting
	var err error

	if err = resp.Bind(&setting, "E1000127"); err != nil {
		return
	}

	setting.CreatedBy = params.UserID
	if createdSetting, err = p.Service.Create(setting, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(settingenum.CreateSetting, createdSetting)
	resp.Status(http.StatusOK).
		MessageT(terms.VCreatedSuccessfully, dictionary.Translate(terms.Setting)).
		JSON(createdSetting)
}

// Update setting
func (p *SettingAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, settingmodel.Table)
	var setting, settingBefore, settingUpdated settingmodel.Setting
	var err error

	if params.ID, err = resp.GetID(c.Param("settingID"), "E1000128"); err != nil {
		return
	}

	if err = resp.Bind(&setting, "E1000129"); err != nil {
		return
	}

	setting.ID = params.ID
	setting.UpdatedBy = params.UserID
	if settingUpdated, settingBefore, err = p.Service.Save(setting, params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(settingenum.UpdateSetting, settingBefore, settingUpdated)
	resp.Status(http.StatusOK).
		MessageT(terms.VUpdatedSuccessfully, dictionary.Translate(terms.Setting)).
		JSON(settingUpdated)
}

// Delete setting
func (p *SettingAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, settingmodel.Table)
	var setting settingmodel.Setting
	var err error

	if params.ID, err = resp.GetID(c.Param("settingID"), "E1000130"); err != nil {
		return
	}

	if setting, err = p.Service.Delete(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(settingenum.DeleteSetting, setting)
	resp.Status(http.StatusOK).
		MessageT(terms.VDeletedSuccessfully, terms.Setting).
		JSON()
}
