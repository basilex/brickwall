-- name: UserNew :one
insert into users(
    username, password
) values(
    @username, @password
) returning *;

-- name: UserSelect :many
select *
  from users u
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: UserSelectByID :one
select * from users u where u.id = @id;

-- name: UserDeleteByID :one
delete from users u where u.id = @id returning id;