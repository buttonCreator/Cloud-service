package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Cloud/pkg/types"
	"Cloud/pkg/validator"
	"Cloud/usecase"
)

type registerRequest struct {
	ID int `json:"id" validate:"required,min=1"`
}

type updateRequest struct {
	// required: true
	Tokens int `json:"tokens"          validate:"required,min=1"`
	// required: true
	TokensCap int `json:"tokens_cap"      validate:"required,min=1"`
	// required: true
	RatePerMinute int `json:"rate_per_minute" validate:"required,min=1"`
}

// swagger:route POST /v1/user user Register
// Register user
//
// Responses:
//
//	200: SuccessResponse
//	400: commonResponse
//	500: commonResponse
func (a *api) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	val := validator.New[registerRequest]()
	req, err := val.ValidateRequest(ctx, r)
	if err != nil {
		a.handleUseCaseError(w, ctx, usecase.ErrValidationFailed)
		return
	}

	if err = a.app.Register(ctx, req.ID); err != nil {
		a.handleUseCaseError(w, ctx, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Status: successStatus,
	})
}

// swagger:route PUT /v1/user user UpdateUser
// Update data of user
//
// Responses:
//
//	200: SuccessResponse
//	400: commonResponse
//	404: commonResponse
//	500: commonResponse
func (a *api) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	strUserID := r.URL.Query().Get("user_id")
	if strUserID == "" {
		a.handleUseCaseError(w, ctx, usecase.ErrValidationFailed)
		return
	}

	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		a.handleUseCaseError(w, ctx, usecase.ErrValidationFailed)
		return
	}

	val := validator.New[updateRequest]()
	req, err := val.ValidateRequest(ctx, r)
	if err != nil {
		a.handleUseCaseError(w, ctx, usecase.ErrValidationFailed)
		return
	}

	if err = a.app.UpdateUser(ctx, transformUpdateUserRequestToType(req, userID)); err != nil {
		a.handleUseCaseError(w, ctx, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Status: successStatus,
	})
}

// swagger:route Get /v1/user/request user SomeRequest
// Simulate some request from user
//
// Responses:
//
//	200: SuccessResponse
//	400: commonResponse
//	404: commonResponse
//	500: commonResponse
func (a *api) someRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	strUserID := r.URL.Query().Get("user_id")
	if strUserID == "" {
		a.handleUseCaseError(w, ctx, usecase.ErrValidationFailed)
		return
	}

	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		a.handleUseCaseError(w, ctx, usecase.ErrValidationFailed)
		return
	}

	if err = a.app.SomeRequest(ctx, userID); err != nil {
		a.handleUseCaseError(w, ctx, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Status: successStatus,
	})
}

func transformUpdateUserRequestToType(req *updateRequest, userID int) *types.User {
	return &types.User{
		ID:            userID,
		Tokens:        req.Tokens,
		TokensCap:     req.TokensCap,
		RatePerMinute: req.RatePerMinute,
	}
}
