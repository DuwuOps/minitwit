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

create table if not exists message (
  message_id serial primary key,
  author_id integer not null,
  text text not null,
  pub_date integer,
  flagged integer
);
