-- name: CreateLog :one
INSERT INTO logs (level, message, fields, created_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-------------Region S&B Config start-------------
-- name: CreateConfig :one
insert into snb_config (endpoint, facility, device)
values ($1, $2, $3)
returning *;

-- name: GetConfig :one
select *
from snb_config
where facility = $1
  and device = $2;
-------------Region S&B Config end---------------

-------------Region Integrator Config start-------------
-- name: GetIntegratorConfig :one
select *
from integrator_config
where client_id = $1;
-------------Region Integrator Config end---------------

-------------Region Authentication start-------------
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
-------------Region Authentication end---------------