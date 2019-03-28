drop table if exists forum_db.post;
drop table if exists forum_db.vote;
drop type if exists forum_db.vote_type;
drop table if exists forum_db.thread;
drop table if exists forum_db.forum;
drop table if exists forum_db.user;
drop extension if exists citext;

create extension citext;

create table forum_db.user (
  nickname citext primary key,
  fullname varchar(500) not null,
  email citext unique not null,
  about text
);

create table forum_db.forum (
  slug citext primary key,
  title varchar(256) not null,
  author_id citext references forum_db.user(nickname) not null
);

create table forum_db.thread (
  id serial primary key,
  author_id citext references forum_db.user(nickname) not null,
  created timestamp with time zone default now(),
  forum_id citext references forum_db.forum(slug) not null,
  message text not null,
  slug citext unique,
  title varchar(256) not null
);

create table forum_db.post (
  id bigserial primary key,
  author_id citext references forum_db.user(nickname) not null,
  created timestamp with time zone default now(),
  forum_id citext references forum_db.forum(slug),
  isEdited boolean default false not null,
  message text not null,
  parent_id bigint references forum_db.post(id),
  thread_id integer references forum_db.thread(id) not null,

  path integer[] not null
);

create table forum_db.vote (
  user_id citext references forum_db.user(nickname) not null,
  thread_id integer references forum_db.thread(id) not null,
  value integer default 0,
  primary key (user_id, thread_id)
);