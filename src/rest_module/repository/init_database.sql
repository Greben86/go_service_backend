-- name: create-users-table
create table if not exists users (
    id bigserial primary key,
    password varchar(255),
    username varchar(50),
    email varchar(50));

-- name: create-accounts-table
create table if not exists accounts (
    id bigserial primary key,
    name varchar(255),
    bank varchar(255),
    user_id bigint);

-- name: create-accounts-constraint
alter table if exists accounts add constraint accounts_user_id foreign key (user_id) references users;
