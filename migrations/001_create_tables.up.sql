CREATE TABLE "users" (
  "id" varchar(36) PRIMARY KEY,
  "username" varchar(30) NOT NULL UNIQUE,
  "password" varchar(60) NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT NOW(),
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "post" (
  "id" varchar(36) PRIMARY KEY,
  "description" varchar(300),
  "photos" text[],
  "created_at" timestamp NOT NULL DEFAULT NOW(),
  "created_by" varchar(36) NOT NULL REFERENCES "users" ("id"),
  "updated_at" timestamp,
  "updated_by" varchar(36) REFERENCES "users" ("id"),
  "deleted_at" timestamp,
  "deleted_by" varchar(36) REFERENCES "users" ("id")
);