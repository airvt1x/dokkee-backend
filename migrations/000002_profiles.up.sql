CREATE TABLE user_profiles (
    user_id     INTEGER      PRIMARY KEY
                REFERENCES auth_credentials(id) ON DELETE CASCADE,
    first_name  VARCHAR(100) NOT NULL CHECK (char_length(first_name) > 0),
    last_name   VARCHAR(100) NOT NULL CHECK (char_length(last_name) > 0),
    middle_name VARCHAR(100),
    email       VARCHAR(100) NOT NULL UNIQUE,
    phone       VARCHAR(15)  UNIQUE
                CHECK (phone ~ '^\+?[0-9]+$'
                AND char_length(phone) BETWEEN 10 AND 15),
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_balances (
    user_id    INTEGER        PRIMARY KEY
               REFERENCES auth_credentials(id) ON DELETE CASCADE,
    balance    DECIMAL(10, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0),
    updated_at TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
);
