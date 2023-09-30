-- migrate:up
ALTER TABLE round ADD COLUMN first_round BOOLEAN NOT NULL DEFAULT FALSE;

-- migrate:down
ALTER TABLE round DROP COLUMN first_round;