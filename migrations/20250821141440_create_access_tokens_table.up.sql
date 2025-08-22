-- Create access_tokens table
CREATE TABLE access_tokens (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    token VARCHAR(255) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    
    CONSTRAINT fk_access_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_access_tokens_deleted_at ON access_tokens(deleted_at);
CREATE INDEX idx_access_tokens_token ON access_tokens(token);
CREATE INDEX idx_access_tokens_user_id ON access_tokens(user_id);
CREATE INDEX idx_access_tokens_expires_at ON access_tokens(expires_at);
CREATE INDEX idx_access_tokens_is_revoked ON access_tokens(is_revoked);

-- Create composite index for token validation queries
CREATE INDEX idx_access_tokens_validation ON access_tokens(token, is_revoked, expires_at);

-- Create index for cleanup queries
CREATE INDEX idx_access_tokens_cleanup ON access_tokens(expires_at, is_revoked);
