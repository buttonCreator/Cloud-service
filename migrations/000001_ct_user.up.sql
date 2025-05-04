CREATE TABLE IF NOT EXISTS users
(
    id                INT PRIMARY KEY,
    tokens            INT DEFAULT 100,
    tokens_cap        INT DEFAULT 100,
    rate_per_minute   INT DEFAULT 1,
    created_at        TIMESTAMP NOT NULL DEFAULT now(),
    updated_at        TIMESTAMP NOT NULL DEFAULT now(),
    last_addition_at  TIMESTAMP NOT NULL DEFAULT now()
);
