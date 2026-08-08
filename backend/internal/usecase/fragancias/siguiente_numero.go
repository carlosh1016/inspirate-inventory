package fragancias

import (
	"context"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

// SiguienteNumero suggests the next available numero_genero for sedeID and
// genero (masculina/femenina numbered independently). Purely a UI default —
// the caller can still submit any other number on create/update.
func (s *Service) SiguienteNumero(ctx context.Context, sedeID int64, genero string) (int32, error) {
	if genero != "masculina" && genero != "femenina" {
		return 0, domainerrors.NewValidation(
			"Género inválido",
			"El género debe ser masculina o femenina.",
			map[string][]string{"genero": {"Valor no reconocido."}},
		)
	}

	next, err := s.Fragancias.NextNumeroGenero(ctx, sedeID, genero)
	if err != nil {
		return 0, internalErr(err)
	}
	return next, nil
}
