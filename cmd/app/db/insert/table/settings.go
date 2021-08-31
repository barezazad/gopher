package table

import (
	"gopher/entity/service"
	"gopher/entity/setting/settingmodel"
	"gopher/entity/setting/settingrepo"
	"gopher/internal/core"
	"gopher/internal/model"
	"gopher/internal/param"

	"gopher/pkg/dictionary"
)

// InsertSettings for add required settings
func InsertSettings(engine *core.Engine) {

	settingRepo := settingrepo.ProvideSettingRepo(engine)
	settingService := service.ProvideSettingService(settingRepo)
	settings := []settingmodel.Setting{
		{
			Common: model.Common{
				ID: 1,
			},
			Property:    engine.Environments.DefaultLanguage,
			Value:       string(dictionary.En),
			Type:        "string",
			Description: "in case of user JWT not specified this value has been used",
		},
	}

	for _, v := range settings {
		var params param.Param
		params.ID = v.ID
		if _, err := settingService.FindByID(params); err != nil {
			if _, _, err := settingService.Save(v, params); err != nil {
				engine.ServerLog.Fatal("error in saving settings", err)
			}
		} else {
			if _, _, err := settingService.Save(v, params); err != nil {
				engine.ServerLog.Fatal("error in creating settings", err)
			}
		}
	}

}
