-- 创建用户
CREATE USER xhz WITH PASSWORD '***';
-- 创建 db
create database xiaohuazhu with owner postgres;
-- 给权限
GRANT ALL PRIVILEGES ON DATABASE xiaohuazhu TO xhz;

-- 这个在自己的 db 中执行
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO xhz;
-- 授权序列给对应用户即可
GRANT USAGE,SELECT,UPDATE ON ALL SEQUENCES IN SCHEMA public TO xhz;

create table common_time
(
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp
);

create table account
(
    id       serial      not null primary key,
    username varchar(20) not null unique,
    password varchar(80) not null,
    profile_picture varchar(60)
)
    inherits (common_time);

create table record_money
(
    id       serial      not null primary key,
    -- 发布人
    account_id int not null ,
    share bool not null default true,
    money   numeric(10,2) not null default 0,
    describe   varchar(200) not null,
    image  varchar(60) not null
)
    inherits (common_time);


create table account_friend
(
    id       serial      not null primary key,
    -- 当前人
    account_id int not null,
    friend_ids int[]
)
    inherits (common_time);


create table account_friend_apply
(
    id         serial not null primary key,
    -- 被申请人
    account_id int    not null,
    friend_id  int,
    status     int
)
    inherits (common_time);


create table goal
(
    id       serial      not null primary key,
    -- 小组成员
    account_ids int[],
    -- 管理者，才能进行操作
    leader int not null ,
    -- 当前目标
    money   numeric(8,2) not null default 1000,
    -- 类型 1：个人，2小组
    type int not null default 1,
    name varchar(8) not null default '目标'
)
    inherits (common_time);