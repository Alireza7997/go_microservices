-- +migrate UP
CREATE TABLE users (
    Id         INTEGER NOT NULL AUTOINCREMENT PRIMARY KEY,
	first_name  VARCHAR(64),
	last_name   VARCHAR(64),  
	username   VARCHAR(128) NOT NULL UNIQUE,
	email      VARCHAR(64) NOT NULL UNIQUE,
	password   VARCHAR(256) NOT NULL,
	JoinDate   INTEGER NOT NULL
);

-- +migrate Down
DROP TABLE users