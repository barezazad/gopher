// +build wireinject

package server

import (
	"gopher/entity/activity/activityapi"
	"gopher/entity/activity/activityrepo"
	"gopher/entity/auth/authapi"
	"gopher/entity/city/cityapi"
	"gopher/entity/city/cityrepo"
	"gopher/entity/document/documentapi"
	"gopher/entity/document/documentrepo"
	"gopher/entity/gift/giftapi"
	"gopher/entity/gift/giftrepo"
	"gopher/entity/role/roleapi"
	"gopher/entity/role/rolerepo"
	"gopher/entity/service"
	"gopher/entity/setting/settingapi"
	"gopher/entity/setting/settingrepo"
	"gopher/entity/user/userapi"
	"gopher/entity/user/userrepo"
	"gopher/internal/core"

	"github.com/google/wire"
)

func initSettingAPI(e *core.Engine) settingapi.SettingAPI {
	wire.Build(settingrepo.ProvideSettingRepo, service.ProvideSettingService,
		settingapi.ProvideSettingAPI)
	return settingapi.SettingAPI{}
}

func initRoleAPI(e *core.Engine) roleapi.RoleAPI {
	wire.Build(rolerepo.ProvideRoleRepo, service.ProvideRoleService,
		roleapi.ProvideRoleAPI)
	return roleapi.RoleAPI{}
}

func initUserAPI(engine *core.Engine) userapi.UserAPI {
	wire.Build(userrepo.ProvideUserRepo, service.ProvideUserService, userapi.ProvideUserAPI)
	return userapi.UserAPI{}
}

func initAuthAPI(e *core.Engine) authapi.AuthAPI {
	wire.Build(service.ProvideAuthService, authapi.ProvideAuthAPI)
	return authapi.AuthAPI{}
}

func initActivityAPI(engine *core.Engine) activityapi.ActivityAPI {
	wire.Build(activityrepo.ProvideActivityRepo, service.ProvideActivityService, activityapi.ProvideActivityAPI)
	return activityapi.ActivityAPI{}
}

func initDocumentAPI(engine *core.Engine) documentapi.DocumentAPI {
	wire.Build(documentrepo.ProvideDocumentRepo, service.ProvideDocumentService, documentapi.ProvideDocumentAPI)
	return documentapi.DocumentAPI{}
}

func initCityAPI(engine *core.Engine, documentService service.DocumentService) cityapi.CityAPI {
	wire.Build(cityrepo.ProvideCityRepo, service.ProvideCityService, cityapi.ProvideCityAPI)
	return cityapi.CityAPI{}
}

func initGiftAPI(engine *core.Engine, documentService service.DocumentService) giftapi.GiftAPI {
	wire.Build(giftrepo.ProvideGiftRepo, service.ProvideGiftService, giftapi.ProvideGiftAPI)
	return giftapi.GiftAPI{}
}
