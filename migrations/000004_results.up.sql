CREATE TABLE analysis_results (
    id          SERIAL    PRIMARY KEY,
    document_id INTEGER   NOT NULL UNIQUE
                REFERENCES documents(id) ON DELETE CASCADE,
    result_json JSONB     NOT NULL,
    model_used  VARCHAR(100),
    tokens_used INTEGER   CHECK (tokens_used >= 0),
    analyzed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_results_json ON analysis_results USING GIN(result_json);
