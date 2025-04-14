package handlers

import (
	"advertising/advertising-service/internal/models"
	"advertising/advertising-service/internal/repo"
	"advertising/pkg/logger"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type StaticUsecase interface {
	LoadStatic(ctx context.Context, name string) (models.Static, error)
}

type StaticHandler struct {
	su StaticUsecase
}

func NewStaticHandler(su repo.StaticRepo) *StaticHandler {
	return &StaticHandler{
		su: su,
	}
}

func (sh *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	if name == "" {
		http.NotFound(w, r)
		return
	}

	static, err := sh.su.LoadStatic(r.Context(), name)
	if err != nil {
		if errors.Is(err, models.ErrStaticNotFound) {
			http.NotFound(w, r)
			return
		}
		logger.FromCtx(r.Context()).Error("get static", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", static.ContentType)
	w.Header().Add("Content-length", fmt.Sprint(static.Size))
	io.Copy(w, static.Data)
}
