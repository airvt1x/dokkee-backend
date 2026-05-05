CREATE TYPE audit_event_type AS ENUM (
    'USER_REGISTERED',
    'USER_LOGIN',
    'DOCUMENT_UPLOADED',
    'DOCUMENT_ANALYSIS_STARTED',
    'DOCUMENT_ANALYSIS_COMPLETED',
    'DOCUMENT_ANALYSIS_FAILED',
    'RESULT_ACCESSED',
    'PROFILE_ACCESSED',
    'PROFILE_UPDATED',
    'API_REQUEST'
);

CREATE TABLE audit_log (
    id           BIGSERIAL        PRIMARY KEY,
    event_type   audit_event_type NOT NULL,
    user_id_hash VARCHAR(64)      NOT NULL,
    doc_id_hash  VARCHAR(64),
    ip_hash      VARCHAR(64),
    success      BOOLEAN          NOT NULL DEFAULT TRUE,
    error_code   VARCHAR(50),
    created_at   TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_user ON audit_log(user_id_hash);
CREATE INDEX idx_audit_time ON audit_log(created_at DESC);
