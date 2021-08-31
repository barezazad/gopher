package table

import (
	"gopher/entity/service"
	"gopher/entity/user/usermodel"
	"gopher/entity/user/userrepo"
	"gopher/internal/core"
	"gopher/internal/model"
	"gopher/internal/param"
)

// InsertUsers for add required users
func InsertUsers(engine *core.Engine) {
	engine.DB.Exec("DELETE FROM users WHERE username = ?", engine.Environments.AdminUsername)
	userRepo := userrepo.ProvideUserRepo(engine)
	userService := service.ProvideUserService(userRepo)
	users := []usermodel.User{
		{
			Common: model.Common{
				ID:        1,
				CreatedBy: 1,
			},
			RoleID:   1,
			Name:     engine.Environments.AdminUsername,
			Username: engine.Environments.AdminUsername,
			Password: engine.Environments.AdminPassword,
			Lang:     engine.Environments.DefaultLanguage,
			Email:    "barez.azad@gmail.com",
			Status:   "active",
		},
	}

	for _, v := range users {
		var params param.Param
		params.ID = v.ID
		if _, err := userService.FindByID(params); err == nil {
			if _, _, err := userService.Save(v, params); err != nil {
				engine.ServerLog.Fatal("error in saving users", err)
			}
		} else {
			if _, err := userService.Create(v, params); err != nil {
				engine.ServerLog.Fatal("error in creating users", err)
			}
		}
	}

}
