--create types
CREATE TYPE metric_kind AS ENUM ('counter', 'gauge');

--create tables
CREATE TABLE metric(
    id    varchar(255) primary key,
    name  varchar(255) not null,
    kind  metric_kind not null,
    value double precision not null
);