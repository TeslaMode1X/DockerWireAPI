CREATE TABLE IF NOT EXISTS follows (
    following_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    followed_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (following_user_id, followed_user_id)
);