-- +goose Up

-- Enable uuid extension
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- +goose StatementEnd

-- Create table
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    start_date DATE NOT NULL,
    stop_date DATE CHECK (stop_date > start_date),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT statement_timestamp(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT statement_timestamp()
    );
-- +goose StatementEnd

-- Indexes
-- +goose StatementBegin
CREATE INDEX idx_user_subscriptions_service_name ON user_subscriptions(service_name);
CREATE INDEX idx_user_subscriptions_start_date ON user_subscriptions(start_date);
CREATE UNIQUE INDEX idx_user_subscriptions_user_id_service_name ON user_subscriptions(user_id, service_name);
-- +goose StatementEnd

-- Triggers
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    new.updated_at = statement_timestamp();
    RETURN new;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER update_user_subscriptions_updated_at BEFORE UPDATE ON user_subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION months_between(start_date date, end_date date)
RETURNS integer AS $$
BEGIN
RETURN (EXTRACT(year FROM end_date) - EXTRACT(year FROM start_date)) * 12 +
       (EXTRACT(month FROM end_date) - EXTRACT(month FROM start_date));
END;
$$ LANGUAGE 'plpgsql' IMMUTABLE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS months_between();
DROP TRIGGER IF EXISTS update_user_subscriptions_updated_at on user_subscriptions;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_user_subscriptions_user_id_service_name;
DROP INDEX IF EXISTS idx_user_subscriptions_start_date;
DROP INDEX IF EXISTS idx_user_subscriptions_service_name;
DROP TABLE IF EXISTS user_subscriptions;
-- +goose StatementEnd