package server

import (
	"gopher/entity/access/accessenum"
	"gopher/entity/access/accessmiddleware"
	"gopher/internal/core"
	"gopher/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

func Route(rg gin.RouterGroup, engine *core.Engine) {

	authAPI := initAuthAPI(engine)
	userAPI := initUserAPI(engine)
	roleAPI := initRoleAPI(engine)
	settingAPI := initSettingAPI(engine)
	activityAPI := initActivityAPI(engine)
	documentAPI := initDocumentAPI(engine)
	cityAPI := initCityAPI(engine, documentAPI.Service)
	giftAPI := initGiftAPI(engine, documentAPI.Service)

	rg.POST("/login", authAPI.Login)
	rg.POST("/requestResetpassword", authAPI.RequestResetPassword)

	// guard middleware
	rg.Use(middleware.AuthGuard(engine))

	rg.GET("/profile", authAPI.Profile)
	rg.POST("/logout", authAPI.Logout)
	rg.POST("/resetpassword", authAPI.ResetPassword)
	rg.PUT("/updateProfile", authAPI.UpdateProfile)
	rg.PUT("/updateLang", authAPI.UpdateLang)

	// access function
	access := accessmiddleware.NewAccessMiddleware(engine)

	rg.GET("/activities", access.Check(accessenum.SuperAccess), activityAPI.ListAll)
	rg.GET("/activities/self", access.Check(accessenum.ActivitySelf), activityAPI.ListSelf)

	rg.GET("/documents/download/:docType/:docName", access.Check(accessenum.DocumentDownload), documentAPI.DownloadDocs)
	rg.DELETE("/documents/delete/:docType/:documentID", access.Check(accessenum.DocumentDelete), documentAPI.Delete)

	rg.GET("/settings", access.Check(accessenum.SettingRead), settingAPI.List)
	rg.GET("/settings/:settingID", access.Check(accessenum.SettingRead), settingAPI.FindByID)
	rg.PUT("/settings/:settingID", access.Check(accessenum.SettingWrite), settingAPI.Update)

	rg.GET("/roles", access.Check(accessenum.RoleRead), roleAPI.List)
	rg.GET("/resources", access.Check(accessenum.RoleRead), roleAPI.Resources)
	rg.GET("/all/roles", access.Check(accessenum.RoleAll), roleAPI.FindAll)
	rg.GET("/roles/:roleID", access.Check(accessenum.RoleRead), roleAPI.FindByID)
	rg.POST("/roles", access.Check(accessenum.RoleWrite), roleAPI.Create)
	rg.PUT("/roles/:roleID", access.Check(accessenum.RoleWrite), roleAPI.Update)
	rg.DELETE("roles/:roleID", access.Check(accessenum.RoleWrite), roleAPI.Delete)

	rg.GET("/users", access.Check(accessenum.UserRead), userAPI.List)
	rg.GET("/all/users", access.Check(accessenum.UserAll), userAPI.FindAll)
	rg.GET("/users/:userID", access.Check(accessenum.UserRead), userAPI.FindByID)
	rg.POST("/users", access.Check(accessenum.UserWrite), userAPI.Create)
	rg.PUT("/users/:userID", access.Check(accessenum.UserWrite), userAPI.Update)
	rg.DELETE("/users/:userID", access.Check(accessenum.UserWrite), userAPI.Delete)
	rg.GET("/excel/users", access.Check(accessenum.UserExcel), userAPI.Excel)

	rg.GET("/cities", access.Check(accessenum.CityRead), cityAPI.List)
	rg.GET("/all/cities", access.Check(accessenum.CityAll), cityAPI.FindAll)
	rg.GET("/cities/:cityID", access.Check(accessenum.CityRead), cityAPI.FindByID)
	rg.POST("/cities", access.Check(accessenum.CityWrite), cityAPI.Create)
	rg.PUT("/cities/:cityID", access.Check(accessenum.CityWrite), cityAPI.Update)
	rg.DELETE("cities/:cityID", access.Check(accessenum.CityWrite), cityAPI.Delete)

	rg.GET("/gifts", access.Check(accessenum.GiftRead), giftAPI.List)
	rg.GET("/gifts/:giftID", access.Check(accessenum.GiftRead), giftAPI.FindByID)
	rg.POST("/gifts", access.Check(accessenum.GiftWrite), giftAPI.Create)
	rg.PUT("/gifts/:giftID", access.Check(accessenum.GiftWrite), giftAPI.Update)
	rg.DELETE("gifts/:giftID", access.Check(accessenum.GiftWrite), giftAPI.Delete)

}
