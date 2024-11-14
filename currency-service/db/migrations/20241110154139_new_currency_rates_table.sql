-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currency_rates (
    day date NOT NULL,
    currency VARCHAR(255) NOT NULL,
    rate double NOT NULL,
    PRIMARY KEY (day, currency)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currency_rates;
-- +goose StatementEnd
