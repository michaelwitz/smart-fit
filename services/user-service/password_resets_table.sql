-- Create PASSWORD_RESETS table for storing password reset tokens
CREATE TABLE PASSWORD_RESETS (
    email VARCHAR(255) NOT NULL,
    token VARCHAR(32) PRIMARY KEY,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on email for faster lookups
CREATE INDEX idx_password_resets_email ON PASSWORD_RESETS(email);

-- Create an index on expires_at for cleanup queries
CREATE INDEX idx_password_resets_expires_at ON PASSWORD_RESETS(expires_at);
