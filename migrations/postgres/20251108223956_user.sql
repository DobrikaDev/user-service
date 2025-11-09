-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    max_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    geolocation VARCHAR(255) NOT NULL,
    age INT NOT NULL,
    sex VARCHAR(255) NOT NULL,
    about TEXT NOT NULL,
    role VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    reputation_group_id INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
