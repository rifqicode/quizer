CREATE TABLE IF NOT EXISTS user_question_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(user_id, question_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_question_history_user_id ON user_question_history(user_id);
CREATE INDEX IF NOT EXISTS idx_user_question_history_question_id ON user_question_history(question_id);
CREATE INDEX IF NOT EXISTS idx_user_question_history_created_at ON user_question_history(created_at);
