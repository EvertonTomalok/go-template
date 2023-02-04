CREATE TABLE IF NOT EXISTS people
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cpf INT NOT NULL
);


CREATE UNIQUE INDEX if not exists idx_cpf ON people(cpf);