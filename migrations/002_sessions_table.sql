-- +migrate Up
create table sessions (
    id varchar(60) primary key,
    expires_at timestamp not null,
    user_id integer not null,
    foreign key (user_id) references users(id)
);
-- +migrate Down