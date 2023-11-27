-- name: CreateUser :one
insert into users (username, password, permission)
values ($1, $2, $3)
returning *;

-- name: GetUser :one
select *
from users
where username = $1;

-- name: DeleteUser :execresult
delete
from users
where id = $1;