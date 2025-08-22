CREATE TABLE IF NOT EXISTS quiz_answers (
    id SERIAL PRIMARY KEY,
    quiz_session_id INTEGER NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    selected_option_id INTEGER NOT NULL REFERENCES question_options(id) ON DELETE CASCADE,
    is_correct BOOLEAN NOT NULL,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    time_spent_seconds INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(quiz_session_id, question_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_quiz_answers_session_id ON quiz_answers(quiz_session_id);
CREATE INDEX IF NOT EXISTS idx_quiz_answers_question_id ON quiz_answers(question_id);
CREATE INDEX IF NOT EXISTS idx_quiz_answers_is_correct ON quiz_answers(is_correct);
CREATE INDEX IF NOT EXISTS idx_quiz_answers_answered_at ON quiz_answers(answered_at);
