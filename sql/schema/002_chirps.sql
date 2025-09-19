-- +goose up

CREATE TABLE chirps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    body VARCHAR(256) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose down

DROP TABLE chirps;