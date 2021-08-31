package migrate

import (
	"gopher/entity/activity/activitymodel"
	"gopher/entity/city/citymodel"
	"gopher/entity/document/documentmodel"
	"gopher/entity/gift/giftmodel"
	"gopher/entity/role/rolemodel"
	"gopher/entity/setting/settingmodel"
	"gopher/entity/user/usermodel"
	"gopher/internal/core"
)

// Run the database for creating tables
func Run(engine *core.Engine) {
	engine.ActivityDB.AutoMigrate(&activitymodel.Activity{})
	engine.DB.AutoMigrate(&settingmodel.Setting{})
	engine.DB.AutoMigrate(&rolemodel.Role{})
	engine.DB.AutoMigrate(&usermodel.User{})
	engine.DB.Exec("ALTER TABLE users DROP COLUMN `role`;")
	engine.DB.Exec("ALTER TABLE users DROP COLUMN `str_resources`;")
	engine.DB.Exec("ALTER TABLE users ADD CONSTRAINT `fk_users_roles` FOREIGN KEY (role_id) " +
		"REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE RESTRICT;")
	engine.DB.AutoMigrate(&documentmodel.Document{})
	engine.DB.AutoMigrate(&citymodel.City{})
	engine.DB.AutoMigrate(&giftmodel.Gift{})
}
