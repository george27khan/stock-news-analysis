create table "news"
(
    id serial primary key,
    source varchar(100) not null,
    category varchar(100) not null,
    url varchar(100) not null,
    data_json text,
    created_dt timestamp default now() not null,
    is_send bool default false
);