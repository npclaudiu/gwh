create table if not exists gwh_registry (
    id integer not null primary key,
    key text unique not null,
    value text not null
);

create table if not exists gwh_git_repositories (
    id integer not null primary key,
    name text unique not null,
    path text unique not null
);
