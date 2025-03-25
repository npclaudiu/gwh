create table if not exists git_refs (
    id integer not null primary key,
    key text unique not null,
    value text not null
);
