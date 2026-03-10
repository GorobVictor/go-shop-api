CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE role AS ENUM(
    'member',
    'admin'
);