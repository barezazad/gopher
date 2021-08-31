package setting

import (
	"gopher/entity/setting/settingmodel"
	"gopher/entity/setting/settingrepo"
	"gopher/internal/core"
	types "gopher/internal/model"
	"gopher/internal/param"
)

// LoadSetting read settings from database and assign them to the engine.Setting
func LoadSetting(engine *core.Engine) {

	params := param.Param{
		Pagination: param.Pagination{
			Select: "*",
			Order:  "id asc",
			Limit:  1000,
			Offset: 0,
		},
	}

	settingRepo := settingrepo.ProvideSettingRepo(engine)
	var settings []settingmodel.Setting
	var err error
	if settings, err = settingRepo.List(params); err != nil {
		engine.ServerLog.Fatal(err, "failed in loading settings")
	}

	engine.Setting = make(map[string]types.SettingMap, len(settings))

	for _, v := range settings {
		settingVal := types.SettingMap{
			Value: v.Value,
			Type:  v.Type,
		}
		engine.Setting[v.Property] = settingVal
	}

}
