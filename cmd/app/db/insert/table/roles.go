package table

import (
	"gopher/entity/access/accessenum"
	"gopher/entity/role/rolemodel"
	"gopher/entity/role/rolerepo"
	"gopher/entity/service"
	"gopher/internal/core"
	"gopher/internal/model"
	"gopher/internal/param"
	"strings"
)

// InsertRoles for add required roles
func InsertRoles(engine *core.Engine) {
	engine.DB.Exec("UPDATE roles SET deleted_at = null WHERE id IN (1,2,3)")
	roleRepo := rolerepo.ProvideRoleRepo(engine)
	roleService := service.ProvideRoleService(roleRepo)
	roles := []rolemodel.Role{
		{
			Common: model.Common{
				ID:        1,
				CreatedBy: 1,
			},
			Name: "Admin",
			Resources: strings.Join([]string{
				accessenum.SuperAccess,
				accessenum.ReadDeleted,
				accessenum.ActivitySelf, accessenum.ActivityAll,
				accessenum.SettingRead, accessenum.SettingWrite,
				accessenum.UserWrite, accessenum.UserRead, accessenum.UserAll, accessenum.UserExcel,
				accessenum.RoleRead, accessenum.RoleWrite, accessenum.RoleAll,
				accessenum.CityRead, accessenum.CityWrite, accessenum.CityAll,
				accessenum.DocumentDownload, accessenum.DocumentDelete,
				accessenum.GiftRead, accessenum.GiftWrite,
			}, ","),
			Description: "admin has all privileges - do not edit",
		},
		{
			Common: model.Common{
				ID:        2,
				CreatedBy: 1,
			},
			Name: "Reader",
			Resources: strings.Join([]string{
				accessenum.SettingRead,
				accessenum.UserRead, accessenum.UserExcel,
				accessenum.RoleRead,
			}, ","),
			Description: "Reader can see all part without changes",
		},
	}

	for _, v := range roles {
		var params param.Param
		params.ID = v.ID
		if _, err := roleService.FindByID(params); err == nil {
			if _, _, err := roleService.Save(v, params); err != nil {
				engine.ServerLog.Fatal("error in saving roles", err)
			}
		} else {
			if _, err := roleService.Create(v, params); err != nil {
				engine.ServerLog.Fatal("error in creating roles", err)
			}
		}
	}

}
