CREATE TABLE IF NOT EXISTS integration_detail
(
    id                 uuid PRIMARY KEY,
    exam_id            uuid        NOT NULL,
    request_type       varchar(50) NOT NULL,
    request_payload    jsonb       NOT NULL,
    response_payload   jsonb,
    request_status     varchar(20) DEFAULT 'pending',
    response_time      timestamp,
    error_message      text,
    integration_date   timestamp   DEFAULT now(),
    external_reference varchar(255),
    processing_time_ms int,
    created_at timestamp DEFAULT now()
);
