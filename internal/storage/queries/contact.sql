-- name: ContactNew :one
insert into Contact(
    user_id, class, content
) values(
    @user_id, @class, @content
) returning *;

-- name: ContactSelect :many
select *
  from contact p
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: ContactSelectByID :one
select * from contact p where p.id = @id;

-- name: ContactSelectByUserID :many
select * from contact p where p.user_id = @user_id order by p.class;

-- name: ContactSelectByUserIDClass :many
select * from contact p
 where p.user_id = @user_id
   and p.class = @class
   and p.content = @content;

-- name: ContactDeleteByID :one
delete from users u where u.id = @id returning id;