package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strconv"
	"time"

	"github.com/major75/online-subscriptions/internal/models"
	"github.com/major75/online-subscriptions/internal/repository/subscriptions"
	"github.com/major75/online-subscriptions/internal/utils"
	"github.com/major75/online-subscriptions/pkg/logger"
	"github.com/major75/online-subscriptions/pkg/types"
	dtvalidator "github.com/major75/online-subscriptions/pkg/validator"
)

type UserHandler struct {
	logger           logger.Logger
	subscriptionRepo subscriptions.UserSubscription
}

func NewUserHandler(l logger.Logger, r subscriptions.UserSubscription) *UserHandler {
	return &UserHandler{
		logger:           l,
		subscriptionRepo: r,
	}
}

// CreateUserSubscription godoc
// @Summary Create new subscription
// @Description Create new user's subscription
// @Tags subscriptions
// @Produce json
// @Param request body models.CreateUserSubscriptionRequest true "Create user subscription request"
// @Success 201 {object} models.Response{data=models.CreateUserSubscriptionResponse}
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/subscriptions [post]
func (u *UserHandler) CreateUserSubscription(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("Calling CreateUserSubscription")

	var req models.CreateUserSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("Incorrect create subscription parameters: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	u.logger.Debug("CreateUserSubscriptionRequest", "data", req)

	err := dtvalidator.Validate(req)
	if err != nil {
		msg := fmt.Sprintf("Validation error: %v", err)

		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			msg = fmt.Sprintf("Validation error: %v", utils.FormatValidationError(errs))
		}

		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	if req.StopDate != nil && (!req.StartDate.Time.Before(req.StopDate.Time)) {
		utils.RespondWithErrorStatus(w, u.logger, "Start_date should be before stop_date", http.StatusBadRequest)
		return
	}

	response, err := u.subscriptionRepo.Create(r.Context(), req)
	if err != nil {
		msg := fmt.Sprintf("Create subscription query error: %v", err)
		status := http.StatusInternalServerError

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				status = http.StatusConflict
			}
		}

		utils.RespondWithErrorStatus(w, u.logger, msg, status)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.SuccessResponse(response, "Subscription created successfully"))
	u.logger.Info("Called CreateUserSubscription")
}

// GetSubscription godoc
// @Summary Get subscription details
// @Description Get user's online subscription details
// @Tags subscriptions
// @Produce json
// @Success 200 {object} models.Response{data=models.UserSubscription}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/subscriptions/{id} [get]
func (u *UserHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("Calling GetSubscription")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		msg := fmt.Sprintf("Incorrect ID: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	ret, err := u.subscriptionRepo.GetSubscription(r.Context(), uint32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.RespondWithErrorStatus(w, u.logger, "No subscription fetched", http.StatusNotFound)
			return
		}

		msg := fmt.Sprintf("Get subscription query error: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.SuccessResponse(ret, "Subscription fetched successfully"))
	u.logger.Info("Called GetSubscription")
}

// GetUserSubscriptions godoc
// @Summary List user subscriptions
// @Description List user's online subscriptions
// @Tags users
// @Produce json
// @Success 200 {object} models.Response{data=models.UserSubscription}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/{userID}/subscriptions [get]
func (u *UserHandler) GetUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("Calling GetUserSubscriptions")

	idStr := chi.URLParam(r, "userID")
	parsedUUID, err := uuid.Parse(idStr)
	if err != nil {
		msg := fmt.Sprintf("Incorrect UUID: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	ret, err := u.subscriptionRepo.GetSubscriptions(r.Context(), parsedUUID.String())
	if err != nil {
		msg := fmt.Sprintf("Get user subscriptions query error: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusInternalServerError)
		return
	}

	if len(ret) == 0 {
		utils.RespondWithErrorStatus(w, u.logger, "No user subscriptions fetched", http.StatusNotFound)
	} else {
		utils.RespondWithJSON(w, http.StatusOK, models.SuccessResponse(ret, "User subscriptions fetched successfully"))
	}

	u.logger.Info("Called GetUserSubscriptions")
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update user's online subscription
// @Tags subscriptions
// @Produce json
// @Param request body models.CreateUserSubscriptionRequest true "Update user subscription request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/subscriptions/{id} [put]
func (u *UserHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("Calling UpdateSubscription")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		msg := fmt.Sprintf("Incorrect ID: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	var req models.CreateUserSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("Incorrect update parameters: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	u.logger.Debug("CreateUserSubscriptionRequest", "data", req)

	err = dtvalidator.Validate(req)
	if err != nil {
		msg := fmt.Sprintf("Validation error: %v", err)

		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			msg = fmt.Sprintf("Validation error: %v", utils.FormatValidationError(errs))
		}

		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	response, err := u.subscriptionRepo.Update(r.Context(), uint32(id), req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.RespondWithErrorStatus(w, u.logger, "No subscription found", http.StatusNotFound)
			return
		}

		msg := fmt.Sprintf("Update subscription query error: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.SuccessResponse(response, "Subscription updated successfully"))

	u.logger.Info("Called UpdateSubscription")
}

// PatchSubscription godoc
// @Summary Patch subscription
// @Description Patch user's online subscription
// @Tags subscriptions
// @Produce json
// @Param request body models.PatchUserSubscriptionRequest true "Patch user subscription request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/subscriptions/{id} [patch]
func (u *UserHandler) PatchSubscription(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("Calling PatchSubscription")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		msg := fmt.Sprintf("Incorrect ID: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	var req models.PatchUserSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("Incorrect patch parameters: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	u.logger.Debug("PatchUserSubscriptionRequest", "data", req)

	err = dtvalidator.Validate(req)
	if err != nil {
		msg := fmt.Sprintf("Validation error: %v", err)

		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			msg = fmt.Sprintf("Validation error: %v", utils.FormatValidationError(errs))
		}

		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	response, err := u.subscriptionRepo.Patch(r.Context(), uint32(id), req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			utils.RespondWithErrorStatus(w, u.logger, "No subscription patched", http.StatusNotFound)
			return
		}

		msg := fmt.Sprintf("Patch subscription query error: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.SuccessResponse(response, "Subscription patched successfully"))
	u.logger.Info("Called PatchSubscription")
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete user's online subscription
// @Tags subscriptions
// @Produce json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/subscriptions/{id} [delete]
func (u *UserHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("Calling DeleteSubscription")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		msg := fmt.Sprintf("Incorrect ID: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
		return
	}

	rowsCount, err := u.subscriptionRepo.DeleteSubscription(r.Context(), uint32(id))
	if err != nil {
		msg := fmt.Sprintf("Delete subscription query error: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusInternalServerError)
		return
	}

	if rowsCount == 0 {
		utils.RespondWithErrorStatus(w, u.logger, "No subscriptions deleted", http.StatusNotFound)
	} else {
		utils.RespondWithJSON(w, http.StatusOK, models.SuccessResponse(nil, "User subscription successfully deleted"))
	}

	u.logger.Info("Called DeleteSubscription")
}

// GetSubscriptionsTotal godoc
// @Summary Get subscriptions total for selected time interval.
// @Description Get subscriptions total for selected time interval and optional: user_id and service_name. If time range ends in the future then the total will be calculated till now date month
// @Tags subscriptions
// @Produce json
// @Param date_from query string true "Filter by date from (ISO 8601 format, example(07-2026))"
// @Param date_to query string true "Filter by date to (ISO 8601 format, example(07-2026))"
// @Param user_id query string false "Filter by user ID (UUID, example(60601fee-2bf1-4721-ae6f-7636e79a0cba))"
// @Param service_name query string false "Filter by service name"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/subscriptions/total [get]
func (u *UserHandler) GetSubscriptionsTotal(w http.ResponseWriter, r *http.Request) {
	u.logger.Debug("Calling GetSubscriptionsTotal")
	var filter models.SubscriptionsTotalFilter
	if v := r.URL.Query().Get("date_from"); v != "" {
		if t, err := time.Parse(types.DATE_FORMAT, v); err == nil {
			filter.DateFrom = t
		} else {
			utils.RespondWithErrorStatus(w, u.logger, "Invalid date_from parameter value", http.StatusBadRequest)
			return
		}
	} else {
		utils.RespondWithErrorStatus(w, u.logger, "Invalid date_from parameter", http.StatusBadRequest)
		return
	}

	if v := r.URL.Query().Get("date_to"); v != "" {
		if t, err := time.Parse(types.DATE_FORMAT, v); err == nil {
			filter.DateTo = t
		} else {
			utils.RespondWithErrorStatus(w, u.logger, "Invalid date_to parameter value", http.StatusBadRequest)
			return
		}
	} else {
		utils.RespondWithErrorStatus(w, u.logger, "Invalid date_to parameter", http.StatusBadRequest)
		return
	}

	if !filter.DateFrom.Before(filter.DateTo) {
		utils.RespondWithErrorStatus(w, u.logger, "Date_from should be before date_to", http.StatusBadRequest)
		return
	}

	if v := r.URL.Query().Get("user_id"); v != "" {
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			msg := fmt.Sprintf("Incorrect UUID: %v", err)
			utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusBadRequest)
			return
		}
		filter.UserID = &parsedUUID
	}

	if v := r.URL.Query().Get("service_name"); v != "" {
		filter.ServiceName = &v
	}

	total, err := u.subscriptionRepo.GetSubscriptionsTotal(r.Context(), filter)
	if err != nil {
		msg := fmt.Sprintf("Query error: %v", err)
		utils.RespondWithErrorStatus(w, u.logger, msg, http.StatusInternalServerError)
		return
	}

	u.logger.Info("Called GetSubscriptionsTotal")
	utils.RespondWithJSON(w, http.StatusOK, models.SuccessResponse(total, "Totals fetched successfully"))
}
