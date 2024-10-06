-- +goose Up
-- +goose StatementBegin

create extension if not exists pgcrypto;

create table if not exists players
(
    id bigserial primary key,
    email varchar(254) unique not null,
    username varchar(64) unique not null,
    hashed_password char(60) not null,
    registration_date timestamp default current_timestamp not null
);

comment on column players.hashed_password
    is 'Use bcrypt to hash and store with the salt (60 char limit is implementation-dependent)';

create table if not exists game_sessions
(
    id bigserial primary key,
    player_id bigint not null,
    start_time timestamp not null,
    end_time timestamp not null,
    points smallint not null
);

alter table game_sessions drop constraint if exists fk_player_of_session;
alter table game_sessions add constraint fk_player_of_session foreign key (id) references players;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists game_sessions;

drop table if exists players;

drop extension if exists pgcrypto;

-- +goose StatementEnd
