-- +goose Up
create table users (
                      id serial primary key,
                      name text not null,
                      password text not null,
                      email text not null,
                      role int not null,
                      created_at timestamp not null default now(),
                      updated_at timestamp not null default now()
);

-- +goose Down
drop table users;

