CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE role_types AS ENUM('user', 'admin');

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    role role_types NOT NULL, 
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);