CREATE TABLE IF NOT EXISTS "user" (
  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE IF NOT EXISTS "book" (
  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  "title" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL,
  "author" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE IF NOT EXISTS "order" (
  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "date" TIMESTAMPTZ NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE IF NOT EXISTS "order_detail" (
  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  "order_id" UUID NOT NULL,
  "book_id" UUID NOT NULL,
  "quantity" INT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

ALTER TABLE "order" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE;

ALTER TABLE "order_detail" ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id") ON DELETE CASCADE;

ALTER TABLE "order_detail" ADD FOREIGN KEY ("book_id") REFERENCES "book" ("id") ON DELETE CASCADE;