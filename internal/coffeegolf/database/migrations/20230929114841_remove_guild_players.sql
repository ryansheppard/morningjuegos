-- migrate:up
ALTER TABLE tournament DROP CONSTRAINT tournament_guild_id_fkey;
ALTER TABLE tournament_placement DROP CONSTRAINT tournament_placement_player_id_fkey;
ALTER TABLE round DROP CONSTRAINT round_player_id_fkey;
DROP TABLE players;
DROP TABLE guilds;

-- migrate:down
CREATE TABLE players (
    player_id BIGINT PRIMARY KEY,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE guilds (
    guild_id BIGINT PRIMARY KEY,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE tournament ADD CONSTRAINT tournament_guild_id_fkey FOREIGN KEY (guild_id) REFERENCES guilds (guild_id);
ALTER TABLE tournament_placement ADD CONSTRAINT tournament_placement_player_id_fkey FOREIGN KEY (player_id) REFERENCES players (player_id);
ALTER TABLE round ADD CONSTRAINT round_player_id_fkey FOREIGN KEY (player_id) REFERENCES players (player_id);