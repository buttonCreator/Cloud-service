package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"Cloud/pkg/logger"
	"Cloud/pkg/types"
	"Cloud/usecase"
)

type HTTPAPIConfig struct {
	ListenAddress string `default:"0.0.0.0:8080"                      env:"HTTP_ADDR"`
	BaseURL       string `default:"localhost:8080"                    env:"BASE_URL"`
	OpenAPIURL    string `default:"http://localhost:8080/openapi.yml" env:"OPEN_API_URL"`
}

type (
	Application interface {
		Register(ctx context.Context, ID int) error
		UpdateUser(ctx context.Context, updateUser *types.User) error
		SomeRequest(ctx context.Context, ID int) error
	}
)

const (
	successStatus status = "success"
	errorStatus   status = "error"
)

type status string

type api struct {
	app    Application
	logger logger.ILogger
	cfg    HTTPAPIConfig
}

func (a *api) SetupMux(mux *http.ServeMux) {
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/index.html")
	})

	fs := http.FileServer(http.Dir("docs"))
	mux.Handle("/", fs)

	mux.HandleFunc("/v1/user/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			a.register(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/user/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			a.updateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/user/request", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			a.someRequest(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

type response struct {
	Msg    string `json:"message,omitempty"`
	Status status `json:"status"`
	Body   any    `json:"body,omitempty"`
}

func New(log logger.ILogger, cfg HTTPAPIConfig, app Application) *api {
	return &api{
		logger: log,
		cfg:    cfg,
		app:    app,
	}
}

// handleUseCaseError logs the error and sends an appropriate HTTP response.
func (a *api) handleUseCaseError(w http.ResponseWriter, ctx context.Context, err error) {
	if err == nil {
		return
	}

	var statusCode int
	var msg string

	switch {
	case errors.Is(err, usecase.ErrNoContent):
		a.logger.WithError(err).Warn("api not found error")
		statusCode = http.StatusNotFound
		msg = err.Error()

	case errors.Is(err, usecase.ErrDuplicate):
		a.logger.WithError(err).Warn("api duplication error")
		statusCode = http.StatusConflict
		msg = err.Error()

	case errors.Is(err, usecase.ErrValidationFailed):
		a.logger.WithError(err).Warn("api validation error")
		statusCode = http.StatusBadRequest
		msg = err.Error()

	case errors.Is(err, usecase.ErrLimitExceeded):
		a.logger.Ctx(ctx).WithError(err).Error("token limit exceeded")
		statusCode = http.StatusRequestEntityTooLarge
		msg = err.Error()

	default:
		a.logger.Ctx(ctx).WithError(err).Error("api error")
		statusCode = http.StatusInternalServerError
		msg = "internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response{
		Status: errorStatus,
		Msg:    msg,
	})
}
