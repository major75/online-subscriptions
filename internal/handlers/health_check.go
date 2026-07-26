package handlers

import (
	"net/http"

	"github.com/major75/online-subscriptions/internal/models"
	"github.com/major75/online-subscriptions/internal/utils"
	"github.com/major75/online-subscriptions/pkg/logger"
)

type HealthHandler struct {
	logger logger.Logger
}

func NewHealthHandler(l logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger: l,
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags system
// @Produce json
// @Success 200 {object} models.Response
// @Router /api/v1/health [get]
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("health check succeeded")
	utils.RespondWithJSON(w, http.StatusOK, models.SuccessResponse(map[string]string{
		"status": "healthy",
	}, "Service is healthy"))
}
