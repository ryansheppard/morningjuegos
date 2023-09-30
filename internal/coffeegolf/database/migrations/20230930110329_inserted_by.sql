-- migrate:up
ALTER TABLE round ADD inserted_by varchar(50) NOT NULL DEFAULT '';
ALTER TABLE hole ADD inserted_by varchar(50) NOT NULL DEFAULT '';
ALTER TABLE tournament ADD inserted_by varchar(50) NOT NULL DEFAULT '';
ALTER TABLE tournament_placement ADD inserted_by varchar(50) NOT NULL DEFAULT '';


-- migrate:down
ALTER TABLE round DROP COLUMN inserted_by;
ALTER TABLE hole DROP COLUMN inserted_by;
ALTER TABLE tournament DROP COLUMN inserted_by;
ALTER TABLE tournament_placement DROP COLUMN inserted_by;