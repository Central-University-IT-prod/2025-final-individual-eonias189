package rest

import (
	"advertising/pkg/logger"
	"advertising/pkg/middlewares"
	api "advertising/pkg/ogen/advertising-service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
	"go.uber.org/zap"
)

type Server struct {
	srv *http.Server
}

func errorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	if err != nil {
		logger.FromCtx(ctx).Debug("handling error", zap.Error(err))
	}
	switch code {
	case http.StatusBadRequest:
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
	}
	w.WriteHeader(code)
}

func NewServer(handler api.Handler, staticHandler http.Handler, l *zap.Logger) (*Server, error) {
	ogenHandler, err := api.NewServer(handler, api.WithErrorHandler(errorHandler))
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/static/{name}", staticHandler)
	mux.Handle("/", ogenHandler)

	httpHandler := middlewares.Apply(
		mux,
		middlewares.Recover(),
		middlewares.LoggerProvider(l),
		middlewares.Logging(),
		middlewares.Cors(),
	)

	return &Server{
		srv: &http.Server{
			Handler: httpHandler,
		},
	}, nil
}

func (s *Server) Start(ctx context.Context, port int) error {
	s.srv.Addr = fmt.Sprintf(":%d", port)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
