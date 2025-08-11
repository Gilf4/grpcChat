-- +goose Up
create table users (
    id serial primary key ,
    email text not null unique,
    pass_hash bytea not null,
    name text not null,
    created_at timestamp not null default now()
);

-- +goose Down
drop table users;
