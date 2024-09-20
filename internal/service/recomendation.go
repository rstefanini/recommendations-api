package service

import (
	"time"

	"recommendation/internal/model"
	"recommendation/internal/repository"
)

const (
	week      = 7 * 24 * time.Hour
	ratingTop = 3
)

type ProductRecommendationService struct {
	userBehaviorRepository repository.UserBehaviorRepository
}

func NewProductRecommendationService(u repository.UserBehaviorRepository) *ProductRecommendationService {
	return &ProductRecommendationService{
		userBehaviorRepository: u,
	}
}

func (s *ProductRecommendationService) GetProductRecommendation(u model.UserID) (*model.ProductsRecommendation, error) {
	ui, err := s.userBehaviorRepository.GetUserInteractionsSince(u, week)
	if err != nil {
		return nil, err
	}

	ph := buildProductHits(ui)

	return &model.ProductsRecommendation{
		Products: ph.GetTop(ratingTop),
	}, nil
}

func buildProductHits(ui *model.UserInteraction) *model.ProductHits {
	ph := make(model.ProductHits)

	for _, v := range ui.Interactions {
		ph[v.Product] = model.Hits{
			Product: v.Product,
			Hits:    ph[v.Product].Hits + 1,
		}
	}

	return &ph
}
