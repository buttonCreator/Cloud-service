package config

import (
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

func Parse(cfg any, configPath ...string) error {
	if err := parseConfig(cfg, configPath...); err != nil {
		return fmt.Errorf("loading envirement: %w", err)
	}

	return nil
}

func parseConfig(cfg any, configPath ...string) error {
	_ = godotenv.Load()

	configLoader := configor.New(&configor.Config{
		ErrorOnUnmatchedKeys: true,
		Silent:               true,
		Environment:          "",
		ENVPrefix:            "-",
		Verbose:              false,
		Debug:                false,
		AutoReload:           false,
		AutoReloadInterval:   0,
		AutoReloadCallback:   nil,
	})

	if err := configLoader.Load(cfg, configPath...); err != nil {
		return fmt.Errorf("loading envirement: %w", err)
	}

	return nil
}
