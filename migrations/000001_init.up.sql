
CREATE TABLE users (
    id serial PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL CHECK (char_length(first_name) > 0),
    last_name VARCHAR(100) NOT NULL CHECK (char_length(last_name) > 0),
    middle_name VARCHAR(100),
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(15) CHECK (phone ~ '^\+?[0-9]+$' AND char_length(phone) BETWEEN 10 AND 15) UNIQUE,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
