CREATE TABLE IF NOT EXISTS latest_processed (
  latest_processed_id INTEGER PRIMARY KEY,
  date_updated INTEGER
);

IF EXISTS (SELECT * FROM latest_processed WHERE latest_processed_id = 0)
UPDATE latest_processed SET date_updated = 0000000000 WHERE latest_processed_id = 0
ELSE
INSERT INTO latest_processed (latest_processed_id, date_updated) VALUES (0, 0000000000)
END IF
