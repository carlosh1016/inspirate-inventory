// Package errors defines the domain-level error type used across usecases.
// Handlers map DomainError to HTTP responses; nothing in domain or usecase
// imports net/http.
package errors

// Code identifies the category of a DomainError, used by the HTTP layer to
// pick a status code without inspecting error text.
type Code string

const (
	CodeValidation   Code = "validation_error"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeBusinessRule Code = "business_rule_violation"
	CodeInternal     Code = "internal_error"
	CodeRateLimit    Code = "rate_limited"
)

// DomainError is the only error type usecases should return for expected,
// user-facing failures. Title and Detail are Spanish, meant to be shown to
// the end user as-is.
type DomainError struct {
	Code   Code
	Title  string
	Detail string
	Fields map[string][]string
	// Extra carries a structured, error-specific payload the client can act
	// on programmatically (e.g. per-item stock shortfall detail), surfaced
	// as an RFC 7807 extension member. Most errors leave this nil.
	Extra   any
	Wrapped error
}

func (e *DomainError) Error() string {
	if e.Detail != "" {
		return e.Detail
	}
	return e.Title
}

func (e *DomainError) Unwrap() error {
	return e.Wrapped
}

// NewValidation reports malformed input, optionally with per-field messages.
func NewValidation(title, detail string, fields map[string][]string) *DomainError {
	return &DomainError{Code: CodeValidation, Title: title, Detail: detail, Fields: fields}
}

// NewNotFound reports that the requested resource does not exist.
func NewNotFound(title, detail string) *DomainError {
	return &DomainError{Code: CodeNotFound, Title: title, Detail: detail}
}

// NewConflict reports a uniqueness or state conflict (e.g. duplicate correo).
func NewConflict(title, detail string) *DomainError {
	return &DomainError{Code: CodeConflict, Title: title, Detail: detail}
}

// NewBusinessRule reports a violation of a domain invariant that isn't a
// plain field-validation error (e.g. "no puedes desactivar al último admin").
func NewBusinessRule(title, detail string) *DomainError {
	return &DomainError{Code: CodeBusinessRule, Title: title, Detail: detail}
}

// NewUnauthorized reports missing or invalid credentials.
func NewUnauthorized(title, detail string) *DomainError {
	return &DomainError{Code: CodeUnauthorized, Title: title, Detail: detail}
}

// NewForbidden reports that the caller is authenticated but not allowed.
func NewForbidden(title, detail string) *DomainError {
	return &DomainError{Code: CodeForbidden, Title: title, Detail: detail}
}

// NewRateLimit reports that the caller exceeded a rate limit.
func NewRateLimit(title, detail string) *DomainError {
	return &DomainError{Code: CodeRateLimit, Title: title, Detail: detail}
}

// NewInternal wraps an unexpected error. Detail should be a generic Spanish
// message safe to show publicly; the real error goes in Wrapped for logging.
func NewInternal(title, detail string, err error) *DomainError {
	return &DomainError{Code: CodeInternal, Title: title, Detail: detail, Wrapped: err}
}
