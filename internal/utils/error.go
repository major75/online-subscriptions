package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/major75/online-subscriptions/internal/models"
	"github.com/major75/online-subscriptions/pkg/logger"
)

var (
	ErrDecodingBody = errors.New("body decoding error")
)

func RespondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithErrorStatus(w http.ResponseWriter, logger logger.Logger, message string, status int) {
	logger.Error(message)
	payload := models.ErrorResponse(nil, message)
	RespondWithJSON(w, status, payload)
}

func FormatValidationError(errs validator.ValidationErrors) string {
	fnSep := func(l int) string {
		if l > 0 {
			return ", "
		}
		return ""
	}

	var sOut string
	for _, f := range errs {
		var t string
		switch f.Tag() {
		case "required":
			t = fmt.Sprintf("%s=Required field", f.Field())
		case "gte":
			t = fmt.Sprintf("%s=Value should be greater or equal to: %s", f.Field(), f.Param())
		default:
			t = fmt.Sprintf("%s=Validation error in tag: '%s'", f.Field(), f.Tag())
		}

		sOut = strings.Join([]string{sOut, t}, fnSep(len(sOut)))
	}

	return sOut
}
