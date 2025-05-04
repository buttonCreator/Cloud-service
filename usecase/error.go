package usecase

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type appError uint8

const (
	ErrDuplicate appError = iota + 1
	ErrNoContent
	ErrValidationFailed
	ErrForbidden
	ErrLimitExceeded
)

func (err appError) Error() string {
	switch {
	case errors.Is(err, ErrDuplicate):
		return "duplicate"
	case errors.Is(err, ErrNoContent):
		return "not found"
	case errors.Is(err, ErrValidationFailed):
		return "validation failed"
	case errors.Is(err, ErrForbidden):
		return "forbidden"
	case errors.Is(err, ErrLimitExceeded):
		return "limit exceeded"
	}

	return ""
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNoContent
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23503" {
			return ErrValidationFailed
		}
	}

	return err
}
