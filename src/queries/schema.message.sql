CREATE TABLE IF NOT EXISTS message (
  message_id SERIAL PRIMARY KEY,
  author_id INTEGER NOT NULL,
  text TEXT NOT NULL,
  pub_date INTEGER,
  flagged INTEGER
);


ALTER TABLE message DROP CONSTRAINT IF EXISTS fk_author_id;

ALTER TABLE message
ADD CONSTRAINT fk_author_id
FOREIGN KEY (author_id)
REFERENCES users(user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;