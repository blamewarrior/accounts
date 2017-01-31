CREATE TABLE users (
  id        SERIAL primary key,
  token     VARCHAR NOT NULL,
  uid       VARCHAR NOT NULL,
  nickname  VARCHAR NOT NULL,
  avatar_url VARCHAR NOT NULL,

  name      VARCHAR
);

CREATE UNIQUE INDEX accounts_nickname ON users (nickname);
