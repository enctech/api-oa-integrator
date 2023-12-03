alter table users
    drop column name;
alter table users
    rename column permissions to permission;
alter table users
    alter permission type varchar,
    alter permission set default '';