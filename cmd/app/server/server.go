package server

import (
	"fmt"
	"gopher/internal/core"
	"gopher/internal/core/middleware"
	"gopher/internal/core/terms"
	"gopher/internal/response"
	"gopher/pkg/generr"
	"gopher/pkg/logparser"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start(engine *core.Engine) *gin.Engine {

	var r *gin.Engine

	// TODO to be check [mos] [Default/New]
	if engine.Environments.GinMode == "debug" {
		r = gin.Default()
	} else {
		r = gin.New()
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "127.0.0.1"
		},
		//MaxAge: 12 * time.Hour,
	}))
	r.Use(middleware.APILogger(engine))

	// No Route "Not Found"
	notFoundRoute(r, engine)

	rg := r.Group("/api/gopher/v1")
	{
		Route(*rg, engine)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", engine.Environments.Addr, engine.Environments.Server.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//logger.Info("starting server: ", engine.Environments.Addr, engine.Environments.Port)
	fmt.Printf("starting server: %v:%v\n", engine.Environments.Addr, engine.Environments.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

	return r
}

func notFoundRoute(r *gin.Engine, engine *core.Engine) {
	r.NoRoute(func(c *gin.Context) {
		err := logparser.New("route not found", "E1000000").Custom(generr.RouteNotFoundErr).
			Message(terms.RouteNotFound).Build()
		response.New(engine, c, "server").Error(err).JSON()
	})
}
