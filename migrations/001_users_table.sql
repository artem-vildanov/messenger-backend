-- +migrate Up
create table users (
    id serial primary key,
    username varchar(50) not null,
    password_hash varchar(60) not null
);
-- +migrate Down
