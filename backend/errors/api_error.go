package errors

type CodeError string

const (
	BadRequest          CodeError = "BAD_REQUEST"
	Unauthorized        CodeError = "UNAUTHORIZED"
	Forbidden           CodeError = "FORBIDDEN"
	NotFound            CodeError = "NOT_FOUND"
	MethodNotAllowed    CodeError = "METHOD_NOT_ALLOWED"
	Conflict            CodeError = "CONFLICT"
	UnprocessableEntity CodeError = "UNPROCESSABLE_ENTITY"
	TooManyRequests     CodeError = "TOO_MANY_REQUESTS"
	InternalServerError CodeError = "INTERNAL_SERVER_ERROR"
	BadGateway          CodeError = "BAD_GATEWAY"
	ServiceUnavailable  CodeError = "SERVICE_UNAVAILABLE"
	GatewayTimeout      CodeError = "GATEWAY_TIMEOUT"
	RefreshExpired      CodeError = "REFRESH_EXPIRED"
	TokenInvalid        CodeError = "TOKEN_INVALID"
	ValidationError     CodeError = "VALIDATION_ERROR"
	UnknownError        CodeError = "UNKNOWN_ERROR"
)

type ApiError struct {
	Status  int
	Code    CodeError
	Message string
	Details string
}

func (e *ApiError) Error() string {
	return e.Message
}

func StatusFromCode(code CodeError) int {
	switch code {
	case BadRequest:
		return 400
	case Unauthorized, RefreshExpired, TokenInvalid:
		return 401
	case Forbidden:
		return 403
	case NotFound:
		return 404
	case MethodNotAllowed:
		return 405
	case Conflict:
		return 409
	case UnprocessableEntity, ValidationError:
		return 422
	case TooManyRequests:
		return 429
	case BadGateway:
		return 502
	case ServiceUnavailable:
		return 503
	case GatewayTimeout:
		return 504
	case InternalServerError, UnknownError:
		return 500
	default:
		return 500
	}
}

func NewApiError(code CodeError, message string, details string) *ApiError {
	return &ApiError{
		Status:  StatusFromCode(code),
		Code:    code,
		Message: message,
		Details: details,
	}
}
