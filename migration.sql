CREATE DATABASE mygram;

USE mygram;

DROP INDEX if exists idx_username;
DROP INDEX if exists idx_email;
DROP INDEX if exists idx_user_id;
DROP INDEX if exists idx_comment_user_id;
DROP INDEX if exists idx_comment_photo_id;

drop table if exists "socialmedia";
drop table if exists "comment";
drop table if exists "photo";
drop table if exists "user";


create table if not exists "user" (
  -- id INT PRIMARY KEY,
  id serial NOT NULL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  age INT NOT NULL,
  CHECK (age > 8),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

CREATE INDEX idx_username ON "user" (username);
CREATE INDEX idx_email ON "user" (email);

create table if not exists photo (
  -- id INT PRIMARY KEY,
  id serial NOT NULL PRIMARY KEY,
  user_id INT,
  title VARCHAR(255) NOT NULL,
  caption VARCHAR(255) NOT NULL,
  photo_url VARCHAR(255) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES "user"(id),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

CREATE INDEX idx_user_id ON photo (user_id);

create table if not exists comment (
  -- id INT PRIMARY KEY,
  id serial NOT NULL PRIMARY KEY,
  user_id INT,
  photo_id INT,
  message TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES "user"(id),
  FOREIGN KEY (photo_id) REFERENCES photo(id),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);

CREATE INDEX idx_comment_user_id ON comment (user_id);
CREATE INDEX idx_comment_photo_id ON comment (photo_id);

create table if not exists socialmedia (
  -- id INT PRIMARY KEY,
  id serial NOT NULL PRIMARY KEY,
  user_id INT,
  name VARCHAR(255) NOT NULL,
  social_media_url VARCHAR(255) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES "user"(id),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);


CREATE TYPE account_role AS ENUM ('admin', 'normal');

CREATE TYPE activity_type AS ENUM ('login', 'logout');
create table if not exists user_activities(
	id uuid primary key not null default uuid_generate_v4(),
	user_id INT not null,
	type activity_type not null,
	created_at timestamptz not null default now(),
	updated_at timestamptz not null default now(),
	deleted_at timestamptz
);
CREATE INDEX user_activity_deleted_at ON user_activities(deleted_at);