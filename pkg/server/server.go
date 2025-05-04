package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"Cloud/pkg/logger"
)

const readyProbPath = "/_health/ready"

var checkPathIsProb = func(path string) bool {
	return readyProbPath == path
}

type BalanceStrategy interface {
	NextBackend(backends []*Backend) (*Backend, error)
}

type Backend struct {
	URL   *url.URL
	Alive bool
	Proxy *httputil.ReverseProxy
}

type Server struct {
	*http.Server

	mux       *http.ServeMux
	ready     bool
	balancing bool
	addr      string
	logger    logger.ILogger
	strategy  BalanceStrategy

	backends    []*Backend
	middlewares []func(http.Handler) http.Handler
}

func New(options ...option) *Server {
	mux := http.NewServeMux()

	server := &Server{
		mux:    mux,
		ready:  false,
		Server: &http.Server{Handler: mux},
	}

	server.WithOptions(options...)

	mux.HandleFunc(readyProbPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if server.ready {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"not ready yet"}`))
	})

	return server
}

func (s *Server) Start(ctx context.Context) error {
	var handler http.Handler = s.mux

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		handler = s.middlewares[i](handler)
	}

	if s.balancing {
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			for {
				select {
				case <-ticker.C:
					s.HealthCheck()
				}
			}
		}()

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if checkPathIsProb(r.URL.Path) {
				s.mux.ServeHTTP(w, r)
				return
			}

			attempts := 0
			for attempts < 3 {
				backend, err := s.strategy.NextBackend(s.backends)
				if err != nil {
					s.logger.Error("no available backends")
					w.WriteHeader(http.StatusServiceUnavailable)
					return
				}

				s.logger.Info("forwarding request", "backend", backend.URL.String())
				backend.Proxy.ServeHTTP(w, r)

				if !isProxyError(w) {
					return
				}

				attempts++
				time.Sleep(10 * time.Millisecond)
			}

			w.WriteHeader(http.StatusBadGateway)
		})
	}

	s.Server.Handler = handler
	s.Server.Addr = s.addr

	errs := make(chan error, 1)

	go func() {
		s.logger.Info("start http server", "addr", s.addr)
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.WithError(err).Error("start server")
			errs <- err
		}
	}()

	const waitServerStart = time.Second
	timer := time.NewTimer(waitServerStart)
	defer timer.Stop()

	select {
	case err := <-errs:
		return err
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return errors.New("start server. context deadline exceed")
	}
}

func isProxyError(w http.ResponseWriter) bool {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return w.Header().Get("X-Proxy-Error") == "true"
}

func (s *Server) WithOptions(options ...option) *Server {
	for _, option := range options {
		option(s)
	}
	return s
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (s *Server) HealthCheck() {
	for _, b := range s.backends {
		alive := isBackendAlive(b.URL)
		b.Alive = alive
		s.logger.Info("health check", "backend", b.URL, "alive", alive)
	}
}

// Shutdown the server
func (s *Server) Shutdown() error {
	return s.Server.Shutdown(context.Background())
}

// SetReady set flag whether or not ready to accept client requests
func (s *Server) SetReady(ready bool) {
	s.ready = ready
}
