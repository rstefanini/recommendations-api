package handler

import (
	"encoding/json"
	"net/http"

	"recommendation/internal/model"
)

type ProductRecommendationService interface {
	GetProductRecommendation(u model.UserID) (*model.ProductsRecommendation, error)
}

type ProductRecommendationHandler struct {
	recommendationService ProductRecommendationService
}

func NewProductRecommendationHandler(r ProductRecommendationService) *ProductRecommendationHandler {
	return &ProductRecommendationHandler{
		recommendationService: r,
	}
}

func (h *ProductRecommendationHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := model.UserID(r.PathValue("user_id"))
	if err := userID.Validate(); err != nil {
		err.WriteJSONError(w)
		return
	}

	prodRec, err := h.recommendationService.GetProductRecommendation(userID)
	if err != nil {
		model.NewInternalServerApiError(err.Error()).WriteJSONError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prodRec)
}
