CREATE TABLE accounts (
  id        SERIAL primary key,
  token     VARCHAR,
  uid       VARCHAR,
  nickname  VARCHAR,
  name      VARCHAR,
  avatar_url VARCHAR,
  rating    jsonb DEFAULT '{}'
);

CREATE INDEX accounts_rating_idx ON accounts USING GIN (rating);

CREATE UNIQUE INDEX accounts_nickname ON accounts (nickname);
