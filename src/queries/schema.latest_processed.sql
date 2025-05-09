CREATE TABLE IF NOT EXISTS latest_processed (
  latest_processed_id INTEGER PRIMARY KEY,
  date_updated INTEGER
);

INSERT INTO latest_processed (latest_processed_id, date_updated) VALUES (0, 0000000000) ON CONFLICT DO NOTHING;