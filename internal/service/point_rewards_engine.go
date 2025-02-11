package service

import (
	"context"

	"github.com/google/uuid"
)

func (s *PointsService) CalculateRewards(ctx context.Context, customerID, programID uuid.UUID, points int) (int, error) {
	return 0, nil
}
