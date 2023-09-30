-- migrate:up
ALTER TABLE hole ALTER COLUMN round_id TYPE integer USING round_id::integer;

-- migrate:down

ALTER TABLE hole ALTER COLUMN round_id TYPE BIGINT USING round_id::BIGINT;