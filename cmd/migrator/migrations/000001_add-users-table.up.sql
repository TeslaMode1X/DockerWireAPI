CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    role varchar NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);