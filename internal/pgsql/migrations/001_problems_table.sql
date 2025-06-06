CREATE TABLE problems (
    uuid string primary key,
    name string,
    kind string,
    status string,
    startTime timestamp,
    solvedTime timestamp,
);

---- create above / drop below ----

drop table problems;
