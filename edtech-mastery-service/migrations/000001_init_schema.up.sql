-- events: append-only, unique (source, event_id)
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL,
    source VARCHAR(255) NOT NULL,
    type VARCHAR(64) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source, event_id)
);

CREATE INDEX idx_events_created_at ON events(created_at);
CREATE INDEX idx_events_type ON events(type);

-- outbox for worker consumption (postgres-backed queue)
CREATE TABLE IF NOT EXISTS event_outbox (
    id BIGSERIAL PRIMARY KEY,
    event_db_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_status ON event_outbox(status) WHERE status = 'pending';

-- student_mastery: materialized per standard
CREATE TABLE IF NOT EXISTS student_mastery (
    student_id VARCHAR(255) NOT NULL,
    standard_id VARCHAR(255) NOT NULL,
    mastery_score DECIMAL(5,4) NOT NULL CHECK (mastery_score >= 0 AND mastery_score <= 1),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (student_id, standard_id)
);

CREATE INDEX idx_student_mastery_student ON student_mastery(student_id);

-- class_rollups: materialized per class
CREATE TABLE IF NOT EXISTS class_rollups (
    class_id VARCHAR(255) PRIMARY KEY,
    completion_rate DECIMAL(5,4) NOT NULL CHECK (completion_rate >= 0 AND completion_rate <= 1),
    avg_score DECIMAL(5,2),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- risk_flags: at-risk students per class
CREATE TABLE IF NOT EXISTS risk_flags (
    id BIGSERIAL PRIMARY KEY,
    student_id VARCHAR(255) NOT NULL,
    class_id VARCHAR(255) NOT NULL,
    reason VARCHAR(255) NOT NULL,
    computed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(student_id, class_id, reason)
);

CREATE INDEX idx_risk_flags_class ON risk_flags(class_id);
