CREATE EXTENSION citext;

CREATE TABLE IF NOT EXISTS movies (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), 
    title text NOT NULL,
    year integer NOT NULL,
    runtime integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);