-- name: CreateUser :one
insert into users (name, username, password, permissions)
values ($1, $2, $3, $4)
returning *;

-- name: GetUser :one
select *
from users
where username = $1;

-- name: GetUsers :many
select *
from users;

-- name: DeleteUser :execresult
delete
from users
where id = $1;