package usecase

import (
	"context"

	"Cloud/pkg/logger"
	"Cloud/pkg/types"
)

type UseCase struct {
	repo   Repository
	logger logger.ILogger
}

type (
	Repository interface {
		FindUserByID(ctx context.Context, ID int) (*types.User, error)
		CreateUser(ctx context.Context, ID int) error
		UpdateUser(ctx context.Context, updateUser *types.User) error
	}
)

func New( //nolint:revive
	log logger.ILogger,
	repo Repository,
) *UseCase {
	return &UseCase{
		logger: log,
		repo:   repo,
	}
}
