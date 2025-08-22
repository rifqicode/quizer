CREATE TABLE IF NOT EXISTS question_options (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_question_options_question_id ON question_options(question_id);
CREATE INDEX IF NOT EXISTS idx_question_options_is_correct ON question_options(is_correct);

-- Ensure only one correct answer per question
CREATE UNIQUE INDEX IF NOT EXISTS idx_question_options_unique_correct 
ON question_options(question_id) 
WHERE is_correct = true;
