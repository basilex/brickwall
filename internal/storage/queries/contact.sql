-- name: ContactNew :one
insert into Contact(
    user_id, class, content
) values(
    @user_id, @class, @content
) returning *;

-- name: ContactCount :one
select count(*) from contact;

-- name: ContactCountByUserID :one
select count(*) from contact where user_id = @user_id;

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

-- name: ContactUpdateByID :one
update contact
   set class = @class, content = @content
 where id = @id returning *;

-- name: ContactDeleteByID :one
delete from users u where u.id = @id returning id;