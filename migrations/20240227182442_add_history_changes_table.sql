-- +goose Up
create table history_changes
(
    id         serial primary key,
    entity     text      not null,
    entity_id  serial    not null,
    value      jsonb,
    created_at timestamp not null default now()

);

-- +goose Down
drop table history_changes;

