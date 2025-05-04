package handler

// ----------------------------------------
// User
// ----------------------------------------

// swagger:parameters Register
type registerRequestWrapper struct { //nolint:unused
	// in: body
	Body registerRequest
}

// swagger:parameters UpdateUser
type updateUserRequestWrapper struct { //nolint:unused
	// User ID to update
	// in: query
	// required: true
	UserID string `json:"user_id"`

	// in: body
	Body updateRequest
}

// swagger:parameters SomeRequest
type someRequestRequestWrapper struct { //nolint:unused
	// User ID making the request
	// in: query
	// required: true
	UserID string `json:"user_id"`
}

// Return status
// swagger:response SuccessResponse
type successResponseWrapper struct { //nolint:unused
	// in: body
	Response struct { //nolint
		// example: success
		Status string `json:"status"`
	}
}

// Return error message
// swagger:response commonResponse
type commonResponseWrapper struct { //nolint:unused
	// in: body
	Response struct { //nolint
		// example: error
		Status string `json:"status"`
		// example: error message
		Message string `json:"message"`
	}
}
