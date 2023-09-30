-- migrate:up
ALTER DATABASE coffeegolf SET timezone TO 'America/New_York';

-- migrate:down
ALTER DATABASE coffeegolf SET timezone TO 'UTC';