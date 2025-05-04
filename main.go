package main

import (
	"log"
	"net/http"
	"strings"

	"Cloud/handler"
	"Cloud/pkg/config"
	"Cloud/pkg/logger"
	"Cloud/pkg/runner"
	"Cloud/pkg/server"
	"Cloud/repository"
	"Cloud/usecase"
)

func main() {
	var cfg Config
	if err := config.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	rootLogger, err := logger.New(logger.WithEnv(cfg.CommonConfig.Env))
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(cfg.Repository.ConnString, rootLogger)
	application := usecase.New(rootLogger, repo)
	httpServer := buildHTTPServer(cfg, rootLogger)

	httpAPI := handler.New(rootLogger, cfg.HTTPAPI, application)
	httpAPI.SetupMux(httpServer.Handler.(*http.ServeMux))

	r := runner.New(rootLogger, httpServer, repo)
	if err = r.RunUtilsSignalExit(); err != nil {
		log.Fatal("run services", "error", err)
	}
}

func buildHTTPServer(cfg Config, l logger.ILogger) *server.Server {
	if cfg.CommonConfig.Env == config.EnvCommon {
		return server.New(
			server.WithListenAt(cfg.HTTPAPI.ListenAddress),
			server.WithLogger(l),
			server.WithDefaultMiddlewares(),
		)
	}

	backendsURLs := strings.Split(cfg.ConfigBackendServers.BackendURLs, ",")

	return server.New(
		server.WithListenAt(cfg.HTTPAPI.ListenAddress),
		server.WithLogger(l),
		server.WithLoadBalancer(backendsURLs, cfg.ConfigBackendServers.LoadBalancerStrategy),
		server.WithDefaultMiddlewares(),
	)
}
