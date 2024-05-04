-- +goose Up
create table accesses
(
    id      serial primary key,
    name    text not null,
    role_id int  not null
);

-- +goose Down
drop table accesses;
