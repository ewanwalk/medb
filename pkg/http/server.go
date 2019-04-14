package http

import (
	"context"
	"encoder-backend/pkg/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type Server struct {
	*http.Server
}

func New() *Server {

	bindTo := config.AppBindTo

	if env := os.Getenv(config.EnvBind); len(env) != 0 {
		bindTo = env
	}

	s := &Server{
		Server: &http.Server{
			Addr:           bindTo,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 8 << 11,
		},
	}

	s.Handler = newRouter()

	go func() {
		err := s.ListenAndServe()
		if err != http.ErrServerClosed {
			logrus.WithError(err).Fatalf("http.server: failed to serve")
		}
	}()

	return s
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		logrus.WithError(err).Warn("http.server.close: failed to gracefully shutdown")
	}
}
