CREATE TABLE IF NOT EXISTS campaigns
(
    id serial primary key,
    name text not null
);

INSERT INTO campaigns (name) VALUES ('Первая запись');

CREATE TABLE IF NOT EXISTS items
(
    id serial primary key,
    campaign_id integer references campaigns (id),
    name text not null,
    description text,
    priority serial,
    removed boolean,
    created_at timestamp
    );

CREATE INDEX IF NOT EXISTS item_name_idx ON items (name);
