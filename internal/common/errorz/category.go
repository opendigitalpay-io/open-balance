package errorz

type Category string

const (
	APIError            Category = "API_ERROR"
	InvalidRequestError Category = "INVALID_REQUEST_ERROR"
)
