CREATE TABLE problems (
    uuid varchar(255) primary key,
    name varchar(255),
    kind varchar(255),
    status varchar(255),
    startTime timestamptz,
    solvedTime timestamptz
);

---- create above / drop below ----

DROP TABLE problems;
