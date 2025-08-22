CREATE TABLE IF NOT EXISTS topics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_topics_name ON topics(name);
CREATE INDEX IF NOT EXISTS idx_topics_is_active ON topics(is_active);
