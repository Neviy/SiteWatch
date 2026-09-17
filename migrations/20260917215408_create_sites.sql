-- +goose Up

CREATE TABLE sites (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    status BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    last_check_at TIMESTAMP,

    CONSTRAINT fk_sites_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);

-- +goose Down

DROP TABLE sites;