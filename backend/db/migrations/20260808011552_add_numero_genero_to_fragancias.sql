-- +goose Up
ALTER TABLE fragancias ADD COLUMN numero_genero INTEGER;

-- +goose Down
ALTER TABLE fragancias DROP COLUMN numero_genero;
