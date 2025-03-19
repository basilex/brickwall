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

-- name: UserSelectByUsername :one
select * from users u where u.username = @username;

-- name: UserUpdateCredentialsByID :one
update users
   set username = @username, password = @password where id = @id returning *;

-- name: UserUpdateIsBlockedByID :one
update users
   set is_blocked = @is_blocked where id = @id returning *;

-- name: UserUpdateIsCheckedByID :one
update users
   set is_checked = @is_checked where id = @id returning *;

-- name: UserUpdateVisitedAtByID :one
update users
   set visited_at = @visited_at where id = @id returning *;

-- name: UserUpdateCheckedAtByID :one
update users
   set checked_at = @checked_at where id = @id returning *;

-- name: UserDeleteByID :one
delete from users where id = @id returning id;