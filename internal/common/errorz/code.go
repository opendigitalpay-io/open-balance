package errorz

type Code string

const (
	InvalidValue             Code = "INVALID_VALUE"
	InvalidJSONBody          Code = "INVALID_JSON_BODY"
	MissingRequiredParameter Code = "MISSING_REQUIRED_PARAMETER"
	NotFound                 Code = "NOT_FOUND"
	ServiceUnavailable       Code = "SERVICE_UNAVAILABLE"
	IdempotencyKeyInUse      Code = "IDEMPOTENCY_KEY_IN_USE"
	InvalidTransaction       Code = "INVALID_TRANSACTION"
)
