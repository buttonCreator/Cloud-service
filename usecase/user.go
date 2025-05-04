package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"Cloud/pkg/types"
)

func (u *UseCase) Register(ctx context.Context, ID int) error {
	user, err := u.repo.FindUserByID(ctx, ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return wrapError(err)
	}

	if user != nil {
		return ErrDuplicate
	}

	return wrapError(u.repo.CreateUser(ctx, ID))
}

func (u *UseCase) UpdateUser(ctx context.Context, updateUser *types.User) error {
	user, err := u.repo.FindUserByID(ctx, updateUser.ID)
	if err != nil {
		return wrapError(err)
	}

	timeNow := time.Now()
	user.Tokens = updateUser.Tokens
	user.TokensCap = updateUser.TokensCap
	user.RatePerMinute = updateUser.RatePerMinute
	user.UpdatedAt = &timeNow

	return wrapError(u.repo.UpdateUser(ctx, user))
}

func (u *UseCase) SomeRequest(ctx context.Context, ID int) error {
	user, err := u.repo.FindUserByID(ctx, ID)
	if err != nil {
		return wrapError(err)
	}

	timeNow := time.Now()
	timeDiff := timeNow.Minute() - user.LastAdditionAt.Minute()
	addTokens := timeDiff * user.RatePerMinute

	if timeDiff != 0 {
		if user.Tokens+addTokens > user.TokensCap {
			user.Tokens = user.TokensCap
		} else {
			user.Tokens += addTokens
		}

		user.LastAdditionAt = &timeNow
	}

	if user.Tokens == 0 {
		return wrapError(ErrLimitExceeded)
	}

	user.ID = ID
	user.Tokens--
	user.UpdatedAt = &timeNow

	return wrapError(u.repo.UpdateUser(ctx, user))
}
