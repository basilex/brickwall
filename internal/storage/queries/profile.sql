-- name: ProfileNew :one
insert into profile(
    user_id, firstname, lastname, gender, birthday, avatar_url, enable_2fa, secret_2fa
) values(
    @user_id, @firstname, @lastname, @gender, @birthday, @avatar_url, @enable_2fa, @secret_2fa
) returning *;

-- name: ProfileSelect :many
select *
  from profile p
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: IsUserProfileExists :one
select case when exists (
    select * from profile p where p.user_id = @user_id
) then cast(1 as bit) else cast(0 as bit) end;

-- name: ProfileSelectByID :one
select * from profile p where p.id = @id;

-- name: ProfileSelectByUserID :one
select * from profile p where p.user_id = @id;

-- name: ProfileDeleteByID :one
delete from profile p where p.id = @id returning id;