create table if not exists users (
  user_id serial primary key,
  username text not null,
  email text not null,
  pw_hash text not null
);

create table if not exists follower (
  who_id integer,
  whom_id integer
);

ALTER TABLE follower ADD PRIMARY KEY (who_id, whom_id);

ALTER TABLE follower
ADD CONSTRAINT fk_who_id
FOREIGN KEY (who_id)
REFERENCES users(user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE follower
ADD CONSTRAINT fk_whom_id
FOREIGN KEY (whom_id)
REFERENCES users(user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;


create table if not exists message (
  message_id serial primary key,
  author_id integer not null,
  text text not null,
  pub_date integer,
  flagged integer
);

ALTER TABLE message
ADD CONSTRAINT fk_author_id
FOREIGN KEY (author_id)
REFERENCES users(user_id)
ON DELETE CASCADE
ON UPDATE CASCADE;