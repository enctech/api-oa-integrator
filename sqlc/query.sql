-- name: CreateLog :one
INSERT INTO logs (module, info, extra)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetConfig :one
select *
from snb_config
where facility = $1
  and device = $2;

-- name: CreateUser :one
insert into users (username, password, permission)
values ($1, $2, $3)
returning *;

-- name: GetUser :one
select * from users where username = $1;

-- name: DeleteUser :execresult
delete from users where id = $1;