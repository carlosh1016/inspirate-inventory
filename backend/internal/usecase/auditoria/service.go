// Package auditoria implements the read-only use cases over the audit log:
// listing with filters/pagination and fetching a single evento. The log itself
// is written by many other usecases (see repository/auditoria.Insert).
package auditoria

import (
	auditoriarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
)

const (
	defaultPageSize = 50
	maxPageSize     = 200
)

// Service provides read access to the audit log.
type Service struct {
	repo auditoriarepo.Repository
}

// NewService wires the auditoria repository.
func NewService(repo auditoriarepo.Repository) *Service {
	return &Service{repo: repo}
}

// normalizePaging clamps page/pageSize to sane bounds and returns limit/offset.
func normalizePaging(page, pageSize int) (limit, offset int32, normPage, normSize int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return int32(pageSize), int32((page - 1) * pageSize), page, pageSize
}
