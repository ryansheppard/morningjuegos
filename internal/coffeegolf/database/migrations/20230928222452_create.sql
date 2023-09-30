-- migrate:up
CREATE TABLE players (
    player_id BIGINT PRIMARY KEY,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE guilds (
    guild_id BIGINT PRIMARY KEY,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE tournament (
    id SERIAL PRIMARY KEY,
    guild_id BIGINT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (guild_id) REFERENCES guilds (guild_id)
);

CREATE TABLE tournament_placement (
    tournament_id INT NOT NULL,
    player_id BIGINT NOT NULL,
    tournament_placement INTEGER NOT NULL,
    strokes INTEGER NOT NULL,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tournament_id, player_id),
    FOREIGN KEY (tournament_id) REFERENCES tournament (id),
    FOREIGN KEY (player_id) REFERENCES players (player_id)
);


CREATE TABLE round (
    id SERIAL PRIMARY KEY,
    tournament_id INT NOT NULL,
    player_id BIGINT NOT NULL,
    total_strokes INTEGER NOT NULL,
    original_date VARCHAR(20) NOT NULL,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (tournament_id) REFERENCES tournament (id),
    FOREIGN KEY (player_id) REFERENCES players (player_id)
);

CREATE TABLE hole (
    id SERIAL PRIMARY KEY,
    round_id BIGINT NOT NULL,
    hole_number INTEGER NOT NULL,
    color VARCHAR(20) NOT NULL,
    strokes INTEGER NOT NULL,
    inserted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (round_id) REFERENCES round (id)
);

-- migrate:down
DROP TABLE players;
DROP TABLE guilds;
DROP TABLE tournament;
DROP TABLE tournament_placement;
DROP TABLE round;
DROP TABLE hole;