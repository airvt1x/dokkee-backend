package dokkee

type User struct {
	Id         int     `json:"id" db:"id"`
	Username   string  `json:"username" binding:"required"`
	Password   string  `json:"password" binding:"required"`
	FirstName  string  `json:"first_name" binding:"required"`
	LastName   string  `json:"last_name" binding:"required"`
	MiddleName string  `json:"middle_name"`
	Email      string  `json:"email" binding:"required"`
	Phone      string  `json:"phone" binding:"required"`
	Balance    float64 `json:"balance"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

/*
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
*/
