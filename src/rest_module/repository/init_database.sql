-- Добавление таблицы пользователей
create table if not exists users (
    id bigserial primary key,
    password varchar(255),
    username varchar(50),
    email varchar(50));

-- Добавление таблицы счетов
create table if not exists accounts (
    id bigserial primary key,
    name varchar(255),
    bank varchar(255),
    user_id bigint);

alter table if exists accounts drop constraint if exists accounts_user_id cascade;
alter table if exists accounts add constraint accounts_user_id foreign key (user_id) references users;
