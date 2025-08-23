CREATE TYPE "perms" AS ENUM (
  'viewer',
  'editor',
  'owner'
);

CREATE TYPE "roles" AS ENUM (
  'user',
  'admin'
);

CREATE TABLE "images" (
  "id" int PRIMARY KEY NOT NULL,
  "original_name" varchar NOT NULL,
  "name" varchar(40),
  "description" text,
  "path" varchar NOT NULL,
  "size" int,
  "publisher_id" int,
  "taken_at" timestamp,
  "uploaded_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp
);

CREATE TABLE "previews" (
  "id" int PRIMARY KEY NOT NULL,
  "path" varchar NOT NULL,
  "image_id" int NOT NULL
);

CREATE TABLE "albums" (
  "id" int PRIMARY KEY NOT NULL,
  "name" varchar(40) NOT NULL,
  "description" text,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp
);

CREATE TABLE "albums_images" (
  "album_id" int NOT NULL,
  "image_id" int NOT NULL
);

CREATE TABLE "album_actions" (
  "album_id" int NOT NULL,
  "user_id" int NOT NULL,
  "action" varchar NOT NULL,
  "timestamp" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "users" (
  "id" int PRIMARY KEY NOT NULL,
  "username" varchar(20) UNIQUE NOT NULL,
  "role" roles
);

CREATE TABLE "user_perms" (
  "user_id" int NOT NULL,
  "album_id" int NOT NULL,
  "permissions" perms NOT NULL
);

ALTER TABLE "images" ADD FOREIGN KEY ("publisher_id") REFERENCES "users" ("id");

ALTER TABLE "album_actions" ADD FOREIGN KEY ("album_id") REFERENCES "albums" ("id");

ALTER TABLE "album_actions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_perms" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_perms" ADD FOREIGN KEY ("album_id") REFERENCES "albums" ("id");

ALTER TABLE "albums_images" ADD FOREIGN KEY ("image_id") REFERENCES "images" ("id");

ALTER TABLE "albums_images" ADD FOREIGN KEY ("album_id") REFERENCES "albums" ("id");

ALTER TABLE "previews" ADD FOREIGN KEY ("image_id") REFERENCES "images" ("id");
