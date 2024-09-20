package service

import (
	"recommendation/internal/model"
	"recommendation/internal/repository"
)

type UserInteractionService struct {
	userBehaviorRepository repository.UserBehaviorRepository
}

func NewUserInteractionService(u repository.UserBehaviorRepository) *UserInteractionService {
	return &UserInteractionService{
		userBehaviorRepository: u,
	}
}

func (s *UserInteractionService) Collect(userInteraction model.UserInteraction) error {
	return s.userBehaviorRepository.AddUserInteractions(userInteraction)
}
