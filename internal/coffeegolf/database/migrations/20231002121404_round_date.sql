-- migrate:up
ALTER TABLE round ADD round_date DATE;
UPDATE round SET round_date = to_date(substring(original_date, 1,3) || ' ' || right(original_date, 2), 'Mon DD') + interval '2023 year' WHERE original_date != '';

CREATE FUNCTION update_round_date() returns trigger as $$
BEGIN  
    IF new.round_date is NULL THEN
        new.round_date := to_date(substring(new.original_date, 1,3) || ' ' || right(new.original_date, 2), 'Mon DD') + interval '2023 year';
    END IF;
    RETURN new;
END
$$ language plpgsql;

CREATE TRIGGER update_round_date_trigger BEFORE INSERT OR UPDATE ON round FOR EACH ROW EXECUTE PROCEDURE update_round_date();

-- migrate:down
DROP TRIGGER update_round_date_trigger ON round;
DROP FUNCTION update_round_date();
ALTER TABLE round DROP COLUMN round_date;