package config

type Base struct {
	App    App
	Server Server
}

type App struct {
	Name string `env:"APP_NAME"`
	Env  Env    `env:"APP_ENV"`
}

type Server struct {
	HTTP HTTP
}

type HTTP struct {
	Host string `env:"HTTP_HOST"`
	Port string `default:"8080"             env:"HTTP_PORT"`
}
