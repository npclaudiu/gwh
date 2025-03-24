create table if not exists gwh_registry (
    id integer not null primary key,
    key text unique not null,
    value text not null
);
