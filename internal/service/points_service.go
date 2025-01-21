package service

import (
	"errors"
	"go-playground/internal/domain"
)

type PointsService struct {
	pointsRepo domain.PointsRepository
	eventRepo  domain.EventLogRepository
}

func NewPointsService(pointsRepo domain.PointsRepository, eventRepo domain.EventLogRepository) *PointsService {
	return &PointsService{
		pointsRepo: pointsRepo,
		eventRepo:  eventRepo,
	}
}

func (s *PointsService) GetBalance(userID string) (*domain.PointsBalance, error) {
	return s.pointsRepo.GetByUserID(userID)
}

func (s *PointsService) UpdateBalance(userID string, points int) error {
	balance, err := s.pointsRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	if balance == nil {
		// Create new balance if it doesn't exist
		balance = &domain.PointsBalance{
			UserID:      userID,
			TotalPoints: points,
		}
		return s.pointsRepo.Create(balance)
	}

	// Update existing balance
	balance.TotalPoints += points
	if balance.TotalPoints < 0 {
		return errors.New("insufficient points balance")
	}

	if err := s.pointsRepo.Update(balance); err != nil {
		return err
	}

	// Log the balance update event
	event := &domain.EventLog{
		EventType: "balance_update",
		UserID:    userID,
		Details: map[string]interface{}{
			"points_changed": points,
			"new_balance":    balance.TotalPoints,
		},
	}
	return s.eventRepo.Create(event)
}
