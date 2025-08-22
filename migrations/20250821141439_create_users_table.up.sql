-- Create users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

-- Create index for soft deletes
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Create index for email lookups
CREATE INDEX idx_users_email ON users(email);

-- Create index for active users
CREATE INDEX idx_users_is_active ON users(is_active);
