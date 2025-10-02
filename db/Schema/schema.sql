create table Usuario (
    id_usuario SERIAL PRIMARY KEY,
    nombre_usuario VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    contrasena VARCHAR(255) NOT NULL
);

CREATE TABLE Tarjeta (
    id_tarjeta SERIAL PRIMARY KEY,
    pregunta VARCHAR(255) NOT NULL,
    respuesta VARCHAR(255) NOT NULL,
    opcion_a VARCHAR(255) NOT NULL,
    opcion_b VARCHAR(255) NOT NULL,
    opcion_c VARCHAR(255) NOT NULL,
    id_tema INTEGER NOT NULL
);

CREATE TABLE Tema (
    id_tema SERIAL PRIMARY KEY,
    nombre_tema VARCHAR(100) NOT NULL UNIQUE
);