--create types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'metrickind') THEN
        CREATE TYPE metrickind AS ENUM ('counter', 'gauge');
    END IF;
END$$;

--create tables
CREATE TABLE IF NOT EXISTS metric(
    id    varchar(255) primary key,
    name  varchar(255) not null,
    kind  metrickind not null,
    value double precision
);