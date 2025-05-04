package validator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Validator[T any] struct{}

func New[T any]() Validator[T] {
	return Validator[T]{}
}

func (*Validator[T]) ValidateRequest(ctx context.Context, r *http.Request) (*T, error) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	validate := validator.New()
	if err := validate.StructCtx(ctx, req); err != nil {
		return nil, err
	}

	return &req, nil
}
