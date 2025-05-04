package repository

import (
	"context"

	"Cloud/pkg/types"
)

func (r *Repository) FindUserByID(ctx context.Context, ID int) (*types.User, error) {
	sql := `SELECT id, tokens, tokens_cap, rate_per_minute, created_at, updated_at, last_addition_at
			FROM users
			WHERE id = $1`

	row := r.Conn().QueryRow(ctx, sql, ID)

	var u types.User
	if err := row.Scan(&u.ID, &u.Tokens, &u.TokensCap, &u.RatePerMinute, &u.CreatedAt, &u.UpdatedAt,
		&u.LastAdditionAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *Repository) CreateUser(ctx context.Context, ID int) error {
	sql := `INSERT INTO users(id) VALUES ($1)`

	_, err := r.Conn().Exec(ctx, sql, ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, updateUser *types.User) error {
	sql := `UPDATE users SET tokens = $1, tokens_cap = $2, rate_per_minute = $3, updated_at = $4, last_addition_at = $5
            WHERE id = $6`

	_, err := r.Conn().Exec(
		ctx,
		sql,
		updateUser.Tokens,
		updateUser.TokensCap,
		updateUser.RatePerMinute,
		updateUser.UpdatedAt,
		updateUser.LastAdditionAt,
		updateUser.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
