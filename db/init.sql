DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id serial primary key,
  email text unique not null,
  password text not null,
  salt text not null
);