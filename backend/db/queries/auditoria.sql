-- name: InsertAuditoria :exec
INSERT INTO auditoria (usuario_id, accion, tabla_afectada, registro_id, datos_antes, datos_despues, ip, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
