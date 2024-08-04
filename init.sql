CREATE EXTENSION pgcrypto;

create schema main authorization postgres;
grant all on schema main to postgres;

create table if not exists main.users
(
    id           uuid primary key     default gen_random_uuid(),
    display_name  varchar     not null,
    login         varchar,
    password      varchar,
    created_at   timestamptz not null default now()
);

create table if not exists main.boards
(
    id           uuid primary key     default gen_random_uuid(),
   -- creator_id uuid     references   users (id),
    created_at   timestamptz not null default now(),
    title        varchar     not null,
    columns      varchar[]
);

create table if not exists main.cards
(
    id           uuid primary key     default gen_random_uuid(),
    created_at   timestamptz not null default now(),
    title varchar   not null,
    description varchar,
    board_id       uuid         references   main.boards (id),
    status         varchar  not null                                -- нужно будет записавать статус колонки: "todo", "in_progress", "done"
    -- TODO: Добавить все поля из ТЗ
);

