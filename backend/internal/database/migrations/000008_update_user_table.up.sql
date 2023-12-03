alter table users
    add column if not exists name varchar;
alter table users
    rename column permission to permissions;
alter table users
    alter permissions type text[] using array [permissions],
    alter permissions set default '{}';