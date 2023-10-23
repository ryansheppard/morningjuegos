-- migrate:up
ALTER TABLE tournament ALTER COLUMN start_time TYPE timestamp with time zone USING start_time::timestamp with time zone;
ALTER TABLE tournament ALTER COLUMN end_time TYPE timestamp with time zone USING end_time::timestamp with time zone;

-- migrate:down
ALTER TABLE tournament ALTER COLUMN start_time TYPE timestamp without time zone USING start_time::timestamp without time zone;
ALTER TABLE tournament ALTER COLUMN end_time TYPE timestamp without time zone USING end_time::timestamp without time zone;