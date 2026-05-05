CREATE TYPE document_status AS ENUM ('queued', 'processing', 'completed', 'failed');

CREATE TABLE documents (
    id            SERIAL          PRIMARY KEY,
    user_id       INTEGER         NOT NULL
                  REFERENCES auth_credentials(id) ON DELETE CASCADE,
    original_name VARCHAR(255)    NOT NULL,
    s3_key        VARCHAR(512)    NOT NULL UNIQUE,
    mime_type     VARCHAR(100)    NOT NULL,
    file_size     BIGINT          NOT NULL CHECK (file_size > 0),
    status        document_status NOT NULL DEFAULT 'queued',
    error_msg     TEXT,
    created_at    TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_status  ON documents(status);
