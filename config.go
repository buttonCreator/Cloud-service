package main

import (
	"Cloud/handler"
	"Cloud/pkg/config"
	"Cloud/pkg/server"
)

type Config struct {
	CommonConfig         CommonConfig
	HTTPAPI              handler.HTTPAPIConfig
	Repository           PostgresRepositoryConfig
	ConfigBackendServers ConfigBackendServers
}

type CommonConfig struct {
	Env        config.Env `env:"ENV"`
	Repository PostgresRepositoryConfig
}

type PostgresRepositoryConfig struct {
	ConnString string `env:"POSTGRESQL_CONNECTION_STRING" required:"true"`
}

type ConfigBackendServers struct {
	LoadBalancerStrategy server.BalancerStrategy `env:"LOAD_BALANCER_STRATEGY"`
	BackendURLs          string                  `env:"BACKEND_URLS"`
}
