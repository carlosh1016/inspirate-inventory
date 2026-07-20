package ventas

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	idempotencykeysrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/idempotency_keys"
)

const (
	createVentaEndpoint = "POST /ventas"
	idempotencyKeyTTL   = 24 * time.Hour
)

// cachedResponse is what a successful idempotency lookup returns: the exact
// bytes+status originally sent to the client, replayed verbatim.
type cachedResponse struct {
	Status int32
	Body   []byte
}

// checkIdempotency looks up key (already known non-empty — the handler
// treats "" the same as no header). If a row exists and hasn't expired:
//   - usuario_id or endpoint differ from the current request -> Conflict
//     ("Idempotency-Key ya usado con otro request").
//   - request_hash differs -> Conflict ("Idempotency-Key ya usado con un
//     cuerpo distinto").
//   - otherwise -> the cached response, to be replayed as-is.
//
// Returns (nil, nil) when no live row exists for key (fresh request).
func (s *Service) checkIdempotency(ctx context.Context, key string, requesterID int64, requestHash string) (*cachedResponse, error) {
	row, err := s.IdempotencyKeys.GetByKey(ctx, key)
	if err != nil {
		if errors.Is(err, idempotencykeysrepo.ErrNotFound) {
			return nil, nil
		}
		return nil, internalErr(err)
	}

	if row.UsuarioID != requesterID || row.Endpoint != createVentaEndpoint {
		return nil, domainerrors.NewConflict(
			"Idempotency-Key en uso",
			"Esta llave de idempotencia ya se usó con otra solicitud.",
		)
	}
	if row.RequestHash != requestHash {
		return nil, domainerrors.NewConflict(
			"Idempotency-Key en uso",
			"Esta llave de idempotencia ya se usó con un cuerpo de solicitud distinto.",
		)
	}

	return &cachedResponse{Status: row.ResponseStatus, Body: row.ResponseBody}, nil
}

// StoreIdempotencyResponse persists the status+body sent to the client for
// key, so a future replay can return the same response (response_body is
// JSONB, so Postgres re-serializes whitespace on storage — the JSON content
// is unaffected). Called by the HTTP handler *after* it builds the real
// response DTO — deliberately
// outside CreateVenta's own transaction: only the handler can produce the
// byte-identical response a replay must match (the DTO shape lives in
// internal/http/handlers/ventas, which this usecase package must not
// import), so the store can't happen pre-commit inside createVentaTx
// without duplicating that shaping logic here and risking the two drifting
// apart. The trade-off: if this insert fails after a successful venta
// commit, a client retry with the same key creates a second (harmless,
// correctly consolidated) venta rather than replaying — acceptable, since
// the alternative bug (replaying a response that doesn't match what the
// client actually received) is strictly worse.
func (s *Service) StoreIdempotencyResponse(ctx context.Context, key string, requesterID int64, requestHash string, status int, body []byte) error {
	return s.IdempotencyKeys.Insert(ctx, idempotencykeysrepo.InsertParams{
		Key:            key,
		UsuarioID:      requesterID,
		Endpoint:       createVentaEndpoint,
		ResponseStatus: int32(status),
		ResponseBody:   body,
		RequestHash:    requestHash,
		ExpiresAt:      time.Now().Add(idempotencyKeyTTL),
	})
}
