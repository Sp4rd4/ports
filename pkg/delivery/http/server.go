// http describes http_server that handles user data that is stored in database.
package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	jsoniter "github.com/json-iterator/go"
	"github.com/sp4rd4/ports/pkg/domain"
	"github.com/sp4rd4/ports/pkg/service"
	"go.uber.org/zap"
)

const errorTag = "http"

var json = jsoniter.ConfigFastest

type PortService interface {
	Get(id string) (*domain.Port, error)
}

type PortController struct {
	service PortService
	logger  *zap.Logger
}

func New(service PortService, logger *zap.Logger) *PortController {
	return &PortController{
		service: service,
		logger:  logger,
	}
}

func (pc *PortController) Get(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context())
	portID := chi.URLParam(r, "portID")
	rLog := pc.logger.With(zap.String("reqId", reqID), zap.String("portId", portID))
	port, err := pc.service.Get(portID)
	if err == nil {
		err = renderData(w, http.StatusOK, port)
	} else {
		err = renderError(err, w, rLog)
	}
	if err != nil {
		rLog.Error(fmt.Errorf("[%v] get: %w", errorTag, err).Error())
	}
}

type message struct {
	M string `json:"message"`
}

func renderError(err error, w http.ResponseWriter, logger *zap.Logger) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		err = renderData(w, http.StatusNotFound, message{M: http.StatusText(http.StatusNotFound)})
	case errors.Is(err, service.ErrPortMissingID):
		err = renderData(w, http.StatusBadRequest, message{M: http.StatusText(http.StatusBadRequest)})
	default:
		logger.Error(fmt.Errorf("[%v] render error: %w", errorTag, err).Error())
		err = renderData(w, http.StatusInternalServerError, message{M: http.StatusText(http.StatusInternalServerError)})
	}
	return err
}

func renderData(w http.ResponseWriter, code int, data interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	resp, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("[%v] marshal: %w", errorTag, err)
	}
	_, err = w.Write(resp)
	return fmt.Errorf("[%v] render: %w", errorTag, err)
}
