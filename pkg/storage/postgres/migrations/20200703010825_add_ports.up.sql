CREATE TABLE IF NOT EXISTS "ports" (
  "id" varchar PRIMARY KEY,
  "name" varchar,
  "city" varchar,
  "country" varchar,
  "alias" varchar[],
  "regions" varchar[],
  "coordinates" jsonb,
  "province" varchar,
  "timezone" varchar,
  "unlocs" varchar[],
  "code" varchar
);

