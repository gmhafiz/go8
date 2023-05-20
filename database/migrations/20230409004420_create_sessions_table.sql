-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions
(
    token   TEXT PRIMARY KEY,
    user_id BIGINT      CONSTRAINT session_user_fk REFERENCES users ON DELETE CASCADE ,
    data    BYTEA       NOT NULL,
    expiry  TIMESTAMPTZ NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
