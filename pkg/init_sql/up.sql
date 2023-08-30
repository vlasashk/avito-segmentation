CREATE TABLE IF NOT EXISTS users
(
    "id" BIGSERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS segments
(
    "id"   BIGSERIAL PRIMARY KEY,
    "name" varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_segments
(
    user_id    BIGINT REFERENCES users (id) NOT NULL,
    segment_id BIGINT REFERENCES segments (id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    PRIMARY KEY (user_id, segment_id)
);
