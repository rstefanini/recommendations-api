package handler

import (
	"encoding/json"
	"net/http"

	"recommendation/internal/model"
)

type UserInteractionsHandler struct {
	userInteractionsService UserInteractionService
}

type UserInteractionService interface {
	Collect(userInteraction model.UserInteraction) error
}

func NewUserInteractionHandler(u UserInteractionService) *UserInteractionsHandler {
	return &UserInteractionsHandler{
		userInteractionsService: u,
	}
}

// Collector handles the users products SKU interactions
func (h *UserInteractionsHandler) Collector(w http.ResponseWriter, r *http.Request) {
	var request model.UserInteraction

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		apiErr := model.NewBadRequestApiError("Malformed Body")
		apiErr.WriteJSONError(w)
		return
	}

	if err := request.Validate(); err != nil {
		err.WriteJSONError(w)
		return
	}

	if err := h.userInteractionsService.Collect(request); err != nil {
		model.NewInternalServerApiError(err.Error()).WriteJSONError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
