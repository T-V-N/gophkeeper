BEGIN;

CREATE TABLE IF NOT EXISTS 
USERS
(
    uid serial primary key,
    email varchar unique,
    password_hash varchar, 
    created_at timestamp default current_timestamp,
    confirmed_at timestamp default null,
    verification_code varchar,
    totp_secret varchar,
    totp_enabled boolean default false
);

CREATE TABLE IF NOT EXISTS 
LOG_PASSWORDS
(
    id serial primary key,
    uid integer references users(uid),
    password_hash varchar,
    resource_name varchar,
    login_hash varchar,
    entry_hash varchar not null,
    is_deleted boolean default false
);

CREATE TABLE IF NOT EXISTS 
CARDS
(
    id serial primary key,
    uid integer references users(uid),
    card_number_hash varchar not null,
    valid_until_hash varchar not null,
    CVV_hash varchar not null,
    last_four_digits varchar not null,
    entry_hash varchar not null,
    is_deleted boolean default false
);

CREATE TABLE IF NOT EXISTS 
FILES
(
    id serial primary key,
    uid integer references users(uid),
    file_name varchar,
    s3_link varchar,
    committed_at timestamp default null,
    is_deleted boolean default false
);

CREATE TABLE IF NOT EXISTS 
TEXT_NOTES
(
    id serial primary key,
    uid integer references users(uid),
    note_name varchar,
    note_text_hash varchar,
    entry_hash varchar not null,
    is_deleted boolean default false
);

COMMIT;
