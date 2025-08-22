CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    topic_id INTEGER NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('easy', 'intermediate', 'hard')),
    explanation TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_questions_topic_id ON questions(topic_id);
CREATE INDEX IF NOT EXISTS idx_questions_difficulty ON questions(difficulty);
CREATE INDEX IF NOT EXISTS idx_questions_is_active ON questions(is_active);
CREATE INDEX IF NOT EXISTS idx_questions_topic_difficulty_active ON questions(topic_id, difficulty, is_active);
