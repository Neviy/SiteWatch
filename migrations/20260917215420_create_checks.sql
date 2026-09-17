-- +goose Up

CREATE TABLE checks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    site_id BIGINT NOT NULL,
    status_code INT NOT NULL,
    response_time BIGINT NOT NULL,
    error TEXT NOT NULL DEFAULT '',
    checked_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_checks_site
        FOREIGN KEY (site_id)
        REFERENCES sites(id)
);

-- +goose Down

DROP TABLE checks;