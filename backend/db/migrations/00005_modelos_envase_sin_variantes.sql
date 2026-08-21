-- +goose Up
ALTER TABLE modelos_envase ADD COLUMN tiene_variantes BOOLEAN NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE modelos_envase DROP COLUMN tiene_variantes;
