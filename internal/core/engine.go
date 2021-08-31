package core

import (
	"gopher/config"
	"gopher/entity/activity/activitymodel"
	"gopher/internal/core/cache"
	"gopher/internal/core/db"
	typesModel "gopher/internal/model"
	"gopher/pkg/dictionary"
	"gopher/pkg/generr"
	"log"
	"time"

	envParser "github.com/caarlos0/env/v6"
	"gorm.io/gorm"
)

type Engine struct {
	DB           *gorm.DB
	ActivityDB   *gorm.DB
	Cache        *cache.Client
	ActivityCh   chan activitymodel.Activity
	ServerLog    *generr.Log
	ErrorLog     *generr.Log
	Environments config.Environment
	Setting      map[string]typesModel.SettingMap
	TZ           *time.Location
}

func Initialize() Engine {
	var env config.Environment

	if err := envParser.Parse(&env); err != nil {
		log.Fatalln(err)
	}

	tz, err := time.LoadLocation(env.TimeZone)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	DataDB := db.ConnectDB(env.DataDSN, env.Database.ShowQueryLogs)
	activityDB := db.ConnectActivityDB(env.ActivityDSN)

	dictionary.Init(env.TermsPath, env.TranslateInBackend)

	serverLog, err := generr.InitLog(env.Server.LogFormat, env.Server.LogOutput, env.Server.LogLevel, env.Server.LogIndentation, true)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	errorLog, err := generr.InitLog(env.ErrorLog.LogFormat, env.ErrorLog.LogOutput, env.ErrorLog.LogLevel, env.ErrorLog.LogIndentation, true)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	redisClient, err := db.InitRedis(env.Redis.Address, env.Redis.Password, env.Redis.DB)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	cacheClient, err := cache.RedisCache(redisClient)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	engine := Engine{
		DB:           DataDB,
		Cache:        cacheClient,
		ActivityDB:   activityDB,
		ServerLog:    serverLog,
		ErrorLog:     errorLog,
		Environments: env,
		TZ:           tz,
	}

	//setting.LoadSetting(&engine)
	return engine

}
