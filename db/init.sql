DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id serial primary key,
  email text unique,
  password text,
  salt text
);