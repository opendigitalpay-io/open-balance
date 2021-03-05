package errorz

type Response struct {
	ErrorFields []Field `json:"errors"`
}

type Field struct {
	Category Category `json:"category"`
	Code     Code     `json:"code"`
	DocURL   DocURL   `json:"doc_url"`
	Message  string   `json:"message"`
}

func NewNotFoundError(err error) Response {
	return Response{[]Field{{
		Category: InvalidRequestError,
		Code:     NotFound,
		DocURL:   APIDocURL,
		Message:  err.Error(),
	}}}
}

func NewInvalidValueError(err error) Response {
	return Response{[]Field{{
		Category: InvalidRequestError,
		Code:     InvalidValue,
		DocURL:   APIDocURL,
		Message:  err.Error(),
	}}}
}

func NewInvalidJSONError(err error) Response {
	return Response{[]Field{{
		Category: InvalidRequestError,
		Code:     InvalidJSONBody,
		DocURL:   APIDocURL,
		Message:  err.Error(),
	}}}
}

func NewIdemKeyError(err error) Response {
	return Response{[]Field{{
		Category: InvalidRequestError,
		Code:     IdempotencyKeyInUse,
		DocURL:   APIDocURL,
		Message:  err.Error(),
	}}}
}

func NewTransactionError(err error) Response {
	return Response{[]Field{{
		Category: InvalidRequestError,
		Code: InvalidTransaction,
		DocURL: APIDocURL,
		Message: err.Error(),
	}}}
}

func NewInternalError(err error) Response {
	return Response{[]Field{{
		Category: APIError,
		Code:     ServiceUnavailable,
		DocURL:   APIDocURL,
		Message:  "Something went wrong, please try again.",
	}}}
}
