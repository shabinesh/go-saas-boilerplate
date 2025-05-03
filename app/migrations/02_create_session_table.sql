CREATE TABLE sessions
(
    session_id       VARCHAR(128) PRIMARY KEY,
    user_id          UUID                     NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_accessed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    ip_address       INET,
    user_agent       TEXT,
    is_active        BOOLEAN                  DEFAULT TRUE,
    data             JSONB,
    device_id        VARCHAR(64),
    location         VARCHAR(255),
    refresh_token    VARCHAR(255)
);

CREATE INDEX idx_user_sessions ON sessions (user_id, is_active);
CREATE INDEX idx_expires ON sessions (expires_at);
CREATE INDEX idx_last_accessed ON sessions (last_accessed_at);

-- Trigger to automatically update last_accessed_at
CREATE OR REPLACE FUNCTION update_session_last_accessed()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.last_accessed_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER session_last_accessed
    BEFORE UPDATE
    ON sessions
    FOR EACH ROW
EXECUTE FUNCTION update_session_last_accessed();

-- Recommended maintenance procedure
COMMENT ON TABLE sessions IS 'Stores web application session data with automatic cleanup';

