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
    balance numeric(16,5),
    user_id bigint);

-- Добавление таблицы счетов
create table if not exists cards (
    id bigserial primary key,
    number varchar(255),
    expiration_month int,
    expiration_year int,
    cvv varchar(255),
    user_id bigint,
    account_id bigint);

-- Добавление таблицы операций
create table if not exists operations (
    id bigserial primary key,
    sum_value numeric(16,5),
    operation_type varchar(50),
    user_id bigint,
    account_id bigint);

-- Добавление таблицы кредитов
create table if not exists credits (
    id bigserial primary key,
    amount numeric(16,5),
    rate numeric(16,5),
    month_count int,
    start_date timestamp(6),
    user_id bigint,
    account_id bigint);

-- Добавление таблицы графика платежей
create table if not exists payment_schedules (
    id bigserial primary key,
    expiration_time timestamp(6),
    amount numeric(16,5),
    payment_status smallint check (payment_status between 0 and 1),
    user_id bigint,
    credit_id bigint);

alter table if exists accounts drop constraint if exists accounts_user_id cascade;
alter table if exists accounts add constraint accounts_user_id foreign key (user_id) references users;

alter table if exists cards drop constraint if exists cards_user_id cascade;
alter table if exists cards add constraint cards_user_id foreign key (user_id) references users;

alter table if exists cards drop constraint if exists cards_account_id cascade;
alter table if exists cards add constraint cards_account_id foreign key (account_id) references accounts;

alter table if exists operations drop constraint if exists operations_user_id cascade;
alter table if exists operations add constraint operations_user_id foreign key (user_id) references users;

alter table if exists operations drop constraint if exists operations_account_id cascade;
alter table if exists operations add constraint operations_account_id foreign key (account_id) references accounts;

alter table if exists credits drop constraint if exists credits_user_id cascade;
alter table if exists credits add constraint credits_user_id foreign key (user_id) references users;

alter table if exists credits drop constraint if exists credits_account_id cascade;
alter table if exists credits add constraint credits_account_id foreign key (account_id) references accounts;

alter table if exists payment_schedules drop constraint if exists payment_schedules_user_id cascade;
alter table if exists payment_schedules add constraint payment_schedules_user_id foreign key (user_id) references users;

alter table if exists payment_schedules drop constraint if exists payment_schedules_credit_id cascade;
alter table if exists payment_schedules add constraint payment_schedules_credit_id foreign key (credit_id) references accounts;
