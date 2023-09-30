-- migrate:up
ALTER TABLE round ADD COLUMN percentage varchar(20) NOT NULL DEFAULT FALSE;

-- migrate:down
ALTER TABLE round DROP COLUMN percentage;