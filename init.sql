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
    board varchar   not null,
    board_id uuid not null,
    status         varchar  not null,
    description varchar,
    assignee varchar,
    estimation varchar,
    updated_at timestamptz
);

create table if not exists main.reports
(
    id           uuid primary key     default gen_random_uuid(),
    board       varchar      not null,
    status      varchar[]    not null,
    assignee    varchar      not null,
    count       int          not null,
    estimation  varchar      not null,
    cards       varchar[]
);

