

-- name: GetUsuarioById :one
SELECT id_usuario, nombre_usuario, email, contrasena
FROM Usuario
WHERE id_usuario = $1;

-- name: ListUsuarios :many
SELECT id_usuario, nombre_usuario, email, contrasena
FROM Usuario
ORDER BY nombre_usuario;

-- name: CreateUsuario :one
INSERT INTO Usuario (nombre_usuario, email, contrasena)
VALUES ($1, $2, $3)
RETURNING id_usuario, nombre_usuario, email, contrasena;

-- name: UpdateUsuario :exec
UPDATE Usuario
SET nombre_usuario = $2, email = $3, contrasena = $4
WHERE id_usuario = $1;

-- name: DeleteUsuario :exec
DELETE FROM Usuario
WHERE id_usuario = $1;

-- name: GetTarjetaById :one
SELECT id_tarjeta, pregunta, respuesta, opcion_a, opcion_b, opcion_c, id_tema 
FROM Tarjeta
WHERE id_tarjeta = $1;

-- name: ListTarjetasByTema :many
SELECT id_tarjeta, pregunta, respuesta, opcion_a, opcion_b, opcion_c, id_tema 
FROM Tarjeta
WHERE id_tema = $1
ORDER BY random();

-- name: ListTarjetas :many
SELECT id_tarjeta, pregunta, respuesta, opcion_a, opcion_b, opcion_c, id_tema 
FROM Tarjeta
ORDER BY random();

-- name: CreateTarjeta :one
INSERT INTO Tarjeta (pregunta, respuesta, opcion_a, opcion_b, opcion_c, id_tema)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id_tarjeta, pregunta, respuesta, opcion_a, opcion_b, opcion_c, id_tema;   

-- name: UpdateTarjeta :exec
UPDATE Tarjeta
SET pregunta = $2, respuesta = $3, opcion_a = $4, opcion_b = $5, opcion_c = $6, id_tema = $7
WHERE id_tarjeta = $1;

-- name: DeleteTarjeta :exec
DELETE FROM Tarjeta
WHERE id_tarjeta = $1;

-- name: GetTemaById :one
SELECT id_tema, nombre_tema
FROM Tema
WHERE id_tema = $1;

-- name: ListTemas :many
SELECT id_tema, nombre_tema
FROM Tema
ORDER BY nombre_tema;

-- name: CreateTema :one
INSERT INTO Tema (nombre_tema)
VALUES ($1)
RETURNING id_tema, nombre_tema;

-- name: UpdateTema :exec
UPDATE Tema
SET nombre_tema = $2
WHERE id_tema = $1;

-- name: DeleteTema :exec
DELETE FROM Tema
WHERE id_tema = $1;