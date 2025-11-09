-- +goose Up
-- +goose StatementBegin

CREATE TABLE balances (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    balance INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

ALTER TABLE balances
    ADD CONSTRAINT balances_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(max_id);

CREATE TABLE balance_operations (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    balance_id VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    type VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

ALTER TABLE balance_operations
    ADD CONSTRAINT balance_operations_balance_id_fkey
    FOREIGN KEY (balance_id) REFERENCES balances(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE balance_operations DROP CONSTRAINT IF EXISTS balance_operations_balance_id_fkey;
DROP TABLE balance_operations;
ALTER TABLE balances DROP CONSTRAINT IF EXISTS balances_user_id_fkey;
DROP TABLE balances;
-- +goose StatementEnd
