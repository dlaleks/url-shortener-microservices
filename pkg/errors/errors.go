package errors

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error codes
const (
	// Common errors
	CodeInternal         = "INTERNAL_ERROR"
	CodeValidation       = "VALIDATION_ERROR"
	CodeNotFound         = "NOT_FOUND"
	CodeAlreadyExists    = "ALREADY_EXISTS"
	CodeUnauthorized     = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeRateLimit        = "RATE_LIMITED"
	CodeTimeout          = "TIMEOUT"
	
	// URL service specific
	CodeInvalidURL       = "INVALID_URL"
	CodeURLNotAccessible = "URL_NOT_ACCESSIBLE"
	CodeCustomCodeTaken  = "CUSTOM_CODE_TAKEN"
	CodeURLExpired       = "URL_EXPIRED"
	CodePasswordRequired = "PASSWORD_REQUIRED"
	CodeInvalidPassword  = "INVALID_PASSWORD"
	
	// User service specific
	CodeEmailTaken       = "EMAIL_TAKEN"
	CodeUsernameTaken    = "USERNAME_TAKEN"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeEmailNotVerified = "EMAIL_NOT_VERIFIED"
	CodeInvalidToken     = "INVALID_TOKEN"
	CodeExpiredToken     = "EXPIRED_TOKEN"
	CodeInvalidAPIKey    = "INVALID_API_KEY"
	CodeAPIKeyExpired    = "API_KEY_EXPIRED"
	
	// Analytics service specific
	CodeInvalidDateRange = "INVALID_DATE_RANGE"
	CodeNoAnalyticsData  = "NO_ANALYTICS_DATA"
)

// AppError represents an application error with code, message, and optional fields
type AppError struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Field     string                 `json:"field,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Cause     error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithField adds a field to the error
func (e *AppError) WithField(field string) *AppError {
	e.Field = field
	return e
}

// WithDetail adds a detail to the error
func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithCause adds the underlying cause
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// New creates a new AppError
func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Newf creates a new AppError with formatted message
func Newf(code, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// Wrapf wraps an existing error with formatted message
func Wrapf(err error, code, format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// Common error constructors
func Internal(message string) *AppError {
	return New(CodeInternal, message)
}

func Internalf(format string, args ...interface{}) *AppError {
	return Newf(CodeInternal, format, args...)
}

func Validation(message string) *AppError {
	return New(CodeValidation, message)
}

func Validationf(format string, args ...interface{}) *AppError {
	return Newf(CodeValidation, format, args...)
}

func NotFound(message string) *AppError {
	return New(CodeNotFound, message)
}

func NotFoundf(format string, args ...interface{}) *AppError {
	return Newf(CodeNotFound, format, args...)
}

func AlreadyExists(message string) *AppError {
	return New(CodeAlreadyExists, message)
}

func AlreadyExistsf(format string, args ...interface{}) *AppError {
	return Newf(CodeAlreadyExists, format, args...)
}

func Unauthorized(message string) *AppError {
	return New(CodeUnauthorized, message)
}

func Unauthorizedf(format string, args ...interface{}) *AppError {
	return Newf(CodeUnauthorized, format, args...)
}

func Forbidden(message string) *AppError {
	return New(CodeForbidden, message)
}

func Forbiddenf(format string, args ...interface{}) *AppError {
	return Newf(CodeForbidden, format, args...)
}

func RateLimit(message string) *AppError {
	return New(CodeRateLimit, message)
}

func RateLimitf(format string, args ...interface{}) *AppError {
	return Newf(CodeRateLimit, format, args...)
}

// HTTP status code mapping
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case CodeValidation:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound
	case CodeAlreadyExists:
		return http.StatusConflict
	case CodeUnauthorized, CodeInvalidCredentials, CodeInvalidToken, CodeExpiredToken:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeRateLimit:
		return http.StatusTooManyRequests
	case CodeTimeout:
		return http.StatusRequestTimeout
	case CodeInvalidURL, CodeCustomCodeTaken, CodePasswordRequired, CodeInvalidPassword:
		return http.StatusBadRequest
	case CodeURLNotAccessible, CodeURLExpired:
		return http.StatusUnprocessableEntity
	case CodeEmailTaken, CodeUsernameTaken:
		return http.StatusConflict
	case CodeEmailNotVerified:
		return http.StatusForbidden
	case CodeInvalidAPIKey, CodeAPIKeyExpired:
		return http.StatusUnauthorized
	case CodeInvalidDateRange:
		return http.StatusBadRequest
	case CodeNoAnalyticsData:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// gRPC status code mapping
func (e *AppError) GRPCStatus() *status.Status {
	var code codes.Code
	
	switch e.Code {
	case CodeValidation, CodeInvalidURL, CodeCustomCodeTaken, CodePasswordRequired, CodeInvalidPassword, CodeInvalidDateRange:
		code = codes.InvalidArgument
	case CodeNotFound, CodeNoAnalyticsData:
		code = codes.NotFound
	case CodeAlreadyExists, CodeEmailTaken, CodeUsernameTaken:
		code = codes.AlreadyExists
	case CodeUnauthorized, CodeInvalidCredentials, CodeInvalidToken, CodeExpiredToken, CodeInvalidAPIKey, CodeAPIKeyExpired:
		code = codes.Unauthenticated
	case CodeForbidden, CodeEmailNotVerified:
		code = codes.PermissionDenied
	case CodeRateLimit:
		code = codes.ResourceExhausted
	case CodeTimeout:
		code = codes.DeadlineExceeded
	case CodeURLNotAccessible:
		code = codes.FailedPrecondition
	case CodeURLExpired:
		code = codes.NotFound
	default:
		code = codes.Internal
	}
	
	return status.New(code, e.Message)
}

// IsAppError checks if error is AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError converts error to AppError if possible
func AsAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// FromGRPCError converts gRPC error to AppError
func FromGRPCError(err error) *AppError {
	if err == nil {
		return nil
	}
	
	st, ok := status.FromError(err)
	if !ok {
		return Internal("unknown gRPC error")
	}
	
	var code string
	switch st.Code() {
	case codes.InvalidArgument:
		code = CodeValidation
	case codes.NotFound:
		code = CodeNotFound
	case codes.AlreadyExists:
		code = CodeAlreadyExists
	case codes.Unauthenticated:
		code = CodeUnauthorized
	case codes.PermissionDenied:
		code = CodeForbidden
	case codes.ResourceExhausted:
		code = CodeRateLimit
	case codes.DeadlineExceeded:
		code = CodeTimeout
	case codes.FailedPrecondition:
		code = CodeValidation
	default:
		code = CodeInternal
	}
	
	return &AppError{
		Code:    code,
		Message: st.Message(),
		Cause:   err,
	}
}

// ValidationErrors holds multiple validation errors
type ValidationErrors struct {
	Errors []*AppError `json:"errors"`
}

// Error implements the error interface
func (v *ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", v.Errors[0].Message)
}

// Add adds a validation error
func (v *ValidationErrors) Add(field, code, message string) {
	v.Errors = append(v.Errors, &AppError{
		Code:    code,
		Message: message,
		Field:   field,
	})
}

// Addf adds a validation error with formatted message
func (v *ValidationErrors) Addf(field, code, format string, args ...interface{}) {
	v.Errors = append(v.Errors, &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Field:   field,
	})
}

// HasErrors returns true if there are validation errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// HTTPStatus returns appropriate HTTP status code
func (v *ValidationErrors) HTTPStatus() int {
	return http.StatusBadRequest
}

// NewValidationErrors creates a new ValidationErrors
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]*AppError, 0),
	}
}