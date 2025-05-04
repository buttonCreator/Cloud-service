package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"Cloud/pkg/logger"
	"Cloud/pkg/server/middleware"
)

type option func(s *Server)

// WithDefaultMiddlewares add default middlewares
func WithDefaultMiddlewares() option {
	return func(s *Server) {
		WithMiddlewares(
			middleware.AccessLog(s.logger, checkPathIsProb),
			middleware.Recover(s.logger),
		)(s)
	}
}

// WithMiddlewares add specific middlewares
func WithMiddlewares(middlewares ...func(http.Handler) http.Handler) option {
	return func(s *Server) {
		s.middlewares = append(s.middlewares, middlewares...)
	}
}

// WithLogger set logger
func WithLogger(log logger.ILogger) option {
	return func(s *Server) {
		s.logger = log
	}
}

// WithListenAt set listen address
func WithListenAt(addr string) option {
	return func(s *Server) {
		s.addr = addr
	}
}

// WithLoadBalancer set load balancer with backends URLs and balance strategy
func WithLoadBalancer(backendURLs []string, strategyName BalancerStrategy) option {
	return func(s *Server) {
		creator, exists := strategyCreators[strategyName]
		if !exists {
			s.logger.Error("unknown strategy", "strategy", strategyName)
			return
		}

		strategy := creator()

		var backends []*Backend
		for _, urlStr := range backendURLs {
			u, err := url.Parse(urlStr)
			if err != nil {
				s.logger.WithError(err).Error("parse backend URL")
				continue
			}

			proxy := httputil.NewSingleHostReverseProxy(u)
			proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				s.logger.WithError(err).Error("proxy error")
				for _, b := range s.backends {
					if b.URL.String() == u.String() {
						b.Alive = false
						break
					}
				}
				w.WriteHeader(http.StatusBadGateway)
			}

			backends = append(backends, &Backend{
				URL:   u,
				Alive: true,
				Proxy: proxy,
			})
		}

		s.backends = backends
		s.strategy = strategy
		s.balancing = len(backends) > 0
	}
}
