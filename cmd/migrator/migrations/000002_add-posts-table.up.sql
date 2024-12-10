CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title varchar(255) NOT NULL,
    body TEXT NOT NULL,
    status varchar,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);