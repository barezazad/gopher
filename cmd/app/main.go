package main

import (
	"gopher/cmd/app/db/insert"
	"gopher/cmd/app/db/migrate"
	"gopher/cmd/app/server"
	"gopher/entity/activity/activitymodel"
	"gopher/entity/activity/activityrepo"
	"gopher/entity/service"
	settings "gopher/entity/setting"
	"gopher/internal/core"
)

func main() {
	engine := core.Initialize()

	// migrate tables
	if engine.Environments.AutoMigrate {
		// migrate tables
		migrate.Run(&engine)
		// seed database with test data
		insert.Insert(&engine)
	}

	// to generate new error code
	// generateerrcode.GenerateNewErrCode()

	core.SyncDirectories(&engine)
	settings.LoadSetting(&engine)

	engine.ActivityCh = make(chan activitymodel.Activity, 1)
	activityRepo := activityrepo.ProvideActivityRepo(&engine)
	activityService := service.ProvideActivityService(activityRepo)
	go activityService.ActivityWatcher()

	server.Start(&engine)
}
