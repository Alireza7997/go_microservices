-- +migrate Up
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    first_name  VARCHAR(64),
    last_name   VARCHAR(64),
    username    VARCHAR(128) NOT NULL UNIQUE,
    email       VARCHAR(64) NOT NULL UNIQUE,
    password    VARCHAR(256) NOT NULL,
    joined_date BIGINT NOT NULL DEFAULT 0
);

-- +migrate Down
DROP TABLE users;
