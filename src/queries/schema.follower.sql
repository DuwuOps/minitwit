CREATE TABLE IF NOT EXISTS follower (
  follower_id INTEGER,
  following_id INTEGER
);


ALTER TABLE follower DROP CONSTRAINT IF EXISTS follower_pkey;

ALTER TABLE follower ADD PRIMARY KEY (follower_id, following_id);


ALTER TABLE follower DROP CONSTRAINT IF EXISTS fk_follower_id;

ALTER TABLE follower
ADD CONSTRAINT fk_follower_id
FOREIGN KEY (follower_id)
REFERENCES users(user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE follower DROP CONSTRAINT IF EXISTS fk_following_id;

ALTER TABLE follower
ADD CONSTRAINT fk_following_id
FOREIGN KEY (following_id)
REFERENCES users(user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;