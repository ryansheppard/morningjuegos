-- migrate:up
ALTER TABLE tournament ALTER COLUMN inserted_at TYPE timestamp with time zone USING inserted_at::timestamp with time zone;
ALTER TABLE tournament_placement ALTER COLUMN inserted_at TYPE timestamp with time zone USING inserted_at::timestamp with time zone;
ALTER TABLE round ALTER COLUMN inserted_at TYPE timestamp with time zone USING inserted_at::timestamp with time zone;
ALTER TABLE hole ALTER COLUMN inserted_at TYPE timestamp with time zone USING inserted_at::timestamp with time zone;

-- migrate:down
ALTER TABLE tournament ALTER COLUMN inserted_at TYPE timestamp without time zone USING inserted_at::timestamp without time zone;
ALTER TABLE tournament_placement ALTER COLUMN inserted_at TYPE timestamp without time zone USING inserted_at::timestamp without time zone;
ALTER TABLE round ALTER COLUMN inserted_at TYPE timestamp without time zone USING inserted_at::timestamp without time zone;
ALTER TABLE hole ALTER COLUMN inserted_at TYPE timestamp without time zone USING inserted_at::timestamp without time zone;