--
-- Entity user
--
create table users (
    id              varchar(32)     not null default xid() primary key,
    username        varchar(64)     not null unique,
    password        varchar(255)    not null,

    is_blocked      bool            not null default false,
    is_checked      bool            not null default false,

    blocked_at      timestamp,
    checked_at      timestamp,
    visited_at      timestamp,
    created_at      timestamp       not null default timezone('utc', now()),
    updated_at      timestamp
);

create trigger users_updated_at
	before update on users for each row
	execute procedure trigger_updated_at();

insert into users(username, password, is_checked) values
    ('sys', crypt('passw!rd', gen_salt('bf', 12)), true),
    ('admin', crypt('passw0rd', gen_salt('bf', 12)), true),
    ('basilex', crypt('passw-rd', gen_salt('bf', 12)), true);
--
-- Entity role
--
create table role (
    id              varchar(32)     not null default xid() primary key,
    name            varchar(255)    not null unique,
    created_at      timestamp       not null default timezone('utc', now()),
    updated_at      timestamp
);

create trigger users_updated_at
	before update on role for each row
	execute procedure trigger_updated_at();

insert into role(name) values
    ('SYS'),
    ('ADMIN'),
    ('EDITOR'),
    ('AUDITOR'),
    ('MANAGER'),
    ('REPORTER'),
    ('FINANCIER');
--
-- Entity profile
--
create table profile (
    id              varchar(32)     not null default xid() primary key,
    user_id         varchar(32)     not null unique references users(id) on delete restrict,
    firstname       varchar(255)    not null,
    lastname        varchar(255)    not null,
    gender          varchar(32)     not null default 'unknown' check(gender in('unknown', 'male', 'female')),
    birthday        date            not null default '1000-01-01'::date,
    avatar_url      varchar(255)    not null default '/assets/images/person.jpg',    
    enable_2fa      bool            not null default false,
    secret_2fa      varchar(255),
    created_at      timestamp       not null default timezone('utc', now()),
    updated_at      timestamp
);

create trigger profile_updated_at
	before update on profile for each row
	execute procedure trigger_updated_at();

do $$
declare
    v_user_id profile.id%type;
begin
    select id into v_user_id from users where username = 'sys';
    insert into profile(user_id, firstname, lastname) values(v_user_id, 'Superuser', 'Administrator');

    select id into v_user_id from users where username = 'admin';
    insert into profile(user_id, firstname, lastname) values(v_user_id, 'System', 'Administrator');

    select id into v_user_id from users where username = 'basilex';
    insert into profile(user_id, firstname, lastname, gender, birthday) values(v_user_id, 'Alexander', 'Vasilenko', 'male', '1965-04-03'::date);
end $$;
--
-- Entity contact
--
create table contact (
    id              varchar(32)     not null default xid() primary key,
    user_id         varchar(32)     not null references users(id) on delete cascade,
    class           varchar(32)     not null check(class in ('email', 'phone', 'mobile', 'telegram', 'viber', 'signal')),
    content         varchar(255)    not null,
    created_at      timestamp       not null default timezone('utc', now()),
    updated_at      timestamp
);

create index contact_user_id on contact(user_id);

create unique index contact_class_content_unq on contact(user_id, class, content);
create unique index contact_class_email_unq on contact(content) where class = 'email';
create unique index contact_class_phone_unq on contact(content) where class = 'phone';
create unique index contact_class_mobile_unq on contact(content) where class = 'mobile';
create unique index contact_class_telegram_unq on contact(content) where class = 'telegram';
create unique index contact_class_viber_unq on contact(content) where class = 'viber';
create unique index contact_class_signal_unq on contact(content) where class = 'signal';

create trigger contact_updated_at
	before update on contact for each row
	execute procedure trigger_updated_at();

do $$
declare
    v_user_id users.id%type;
begin
    select id into v_user_id from users where username = 'basilex';

    insert into contact(user_id, class, content) values(v_user_id, 'email', 'alexander.vasilenko@gmail.com');
    insert into contact(user_id, class, content) values(v_user_id, 'email', 'alexander.vasilenko@icloud.com');
    insert into contact(user_id, class, content) values(v_user_id, 'mobile', '+380952066922');
    insert into contact(user_id, class, content) values(v_user_id, 'telegram', '@Basilex');
end $$;
--
-- *** VIEWS LAYER
--
-- View v_user_profile
--
create or replace view v_user_profile as
    select u.id         as user_id,
           p.id         as profile_id,
           u.username   as user_username,
           u.is_blocked as user_is_blocked,
           u.is_checked as user_is_checked,
           p.firstname  as profile_firstname,
           p.lastname   as profile_lastname,
           p.gender     as profile_gender,
           p.birthday   as profile_birthday,
           p.avatar_url as profile_avatar_url,
           p.enable_2fa as profile_enable_2fa,
           p.secret_2fa as profile_secret_2fa
      from users u
      left join profile p on p.user_id = u.id;

