CREATE TABLE IF NOT EXISTS quiz_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('easy', 'intermediate', 'hard')),
    total_questions INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'abandoned')),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    final_score INTEGER NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user_id ON quiz_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_topic_id ON quiz_sessions(topic_id);
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_status ON quiz_sessions(status);
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user_status ON quiz_sessions(user_id, status);
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_started_at ON quiz_sessions(started_at);
