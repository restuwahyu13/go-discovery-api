package main

import (
	"compress/zlib"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/oxequa/grace"
	"github.com/unrolled/secure"

	"github.com/restuwahyu13/discovery-api/configs"
	"github.com/restuwahyu13/discovery-api/handlers"
	"github.com/restuwahyu13/discovery-api/helpers"
	"github.com/restuwahyu13/discovery-api/middlewares"
	"github.com/restuwahyu13/discovery-api/packages"
	"github.com/restuwahyu13/discovery-api/services"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		env    *configs.Environtment = new(configs.Environtment)
		router *chi.Mux              = chi.NewRouter()
	)

	if err := packages.ViperRead(".env", &env); err != nil {
		packages.Logrus("fatal", err)
		return
	}

	consul, err := packages.NewConsul()
	if err != nil {
		packages.Logrus("fatal", err)
		return
	}

	if env.ENV != "production" {
		router.Use(middleware.Logger)
	}

	router.Use(middlewares.AuthToken)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.CleanPath)
	router.Use(middleware.NoCache)
	router.Use(middleware.GetHead)
	router.Use(middleware.Compress(zlib.BestCompression))
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Authorization", "Origin", "X-DISCOVERY-TOKEN"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             900,
	}))
	router.Use(secure.New(secure.Options{
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		STSSeconds:           900,
	}).Handler)

	service := services.NewServiceDiscovery(&services.ServiceDiscovery{Env: env, Consul: consul})
	handler := handlers.NewHandlerDiscovery(&handlers.HandlerDiscovery{Service: service})

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", func(rw http.ResponseWriter, r *http.Request) {
			helpers.ApiResponse(rw, &helpers.Response{StatCode: http.StatusOK, StatMsg: "Ping!"})
		})

		r.Group(func(rg chi.Router) {
			rg.Post("/agent/service/register", handler.Register)
			rg.Delete("/agent/service/deregister/{serviceId}", handler.Deregister)
			rg.Get("/agent/service", handler.ListDiscovery)
			rg.Get("/agent/service/{serviceId}", handler.DetailDiscovery)
		})

		r.Group(func(rg chi.Router) {
			rg.Post("/agent/health/register", handler.CheckRegister)
			rg.Delete("/agent/health/deregister/{healthId}", handler.CheckDeregister)
			rg.Delete("/agent/health/{serviceName}", handler.HealthCheck)
		})
	})

	err = packages.Graceful(func() *packages.GracefulConfig {
		return &packages.GracefulConfig{Handler: router, Port: env.PORT}
	})

	recover := grace.Recover(&err)
	recover.Stack()

	if err != nil {
		packages.Logrus("fatal", err)
		os.Exit(1)
	}
}
