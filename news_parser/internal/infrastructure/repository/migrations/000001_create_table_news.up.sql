create table "news"
(
    id serial primary key,
    source varchar(100) not null,
    category varchar(100) not null,
    article_id int not null,
    title varchar(500) not null,
    url varchar(100) not null,
    data_json text,
    published_at timestamp with time zone not null,
    created_dt timestamp default now() not null,
    is_send bool default false
);

CREATE INDEX idx_news_published_at
    ON news (published_at);