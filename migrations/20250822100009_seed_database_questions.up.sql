-- Insert sample database questions and answers
-- First, get the topic ID for Database topic and insert questions

INSERT INTO questions (topic_id, question_text, difficulty, explanation, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Database' LIMIT 1), 'What does SQL stand for?', 'easy', 'SQL stands for Structured Query Language, which is a standard language for managing and manipulating relational databases.', true),
((SELECT id FROM topics WHERE name = 'Database' LIMIT 1), 'Which of the following is NOT a type of database constraint?', 'intermediate', 'LOOP is not a database constraint. Common constraints include PRIMARY KEY (ensures uniqueness), FOREIGN KEY (maintains referential integrity), and CHECK (validates data against conditions).', true),
((SELECT id FROM topics WHERE name = 'Database' LIMIT 1), 'What is the purpose of database normalization?', 'intermediate', 'Database normalization is the process of organizing data to reduce redundancy and improve data integrity. It involves dividing large tables into smaller ones and defining relationships between them.', true),
((SELECT id FROM topics WHERE name = 'Database' LIMIT 1), 'Which SQL command is used to retrieve data from a database?', 'easy', 'SELECT is the SQL command used to retrieve data from one or more tables in a database. INSERT adds data, UPDATE modifies data, and DELETE removes data.', true),
((SELECT id FROM topics WHERE name = 'Database' LIMIT 1), 'What is ACID in the context of database transactions?', 'hard', 'ACID stands for Atomicity, Consistency, Isolation, and Durability. These are four key properties that guarantee reliable processing of database transactions.', true);

-- Insert question options for question 1: "What does SQL stand for?"
INSERT INTO question_options (question_id, option_text, is_correct) VALUES
((SELECT id FROM questions WHERE question_text = 'What does SQL stand for?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'Structured Query Language', true),
((SELECT id FROM questions WHERE question_text = 'What does SQL stand for?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'Simple Query Language', false),
((SELECT id FROM questions WHERE question_text = 'What does SQL stand for?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'Standard Query Language', false),
((SELECT id FROM questions WHERE question_text = 'What does SQL stand for?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'System Query Language', false);

-- Insert question options for question 2: "Which of the following is NOT a type of database constraint?"
INSERT INTO question_options (question_id, option_text, is_correct) VALUES
((SELECT id FROM questions WHERE question_text = 'Which of the following is NOT a type of database constraint?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'PRIMARY KEY', false),
((SELECT id FROM questions WHERE question_text = 'Which of the following is NOT a type of database constraint?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'FOREIGN KEY', false),
((SELECT id FROM questions WHERE question_text = 'Which of the following is NOT a type of database constraint?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'CHECK', false),
((SELECT id FROM questions WHERE question_text = 'Which of the following is NOT a type of database constraint?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'LOOP', true);

-- Insert question options for question 3: "What is the purpose of database normalization?"
INSERT INTO question_options (question_id, option_text, is_correct) VALUES
((SELECT id FROM questions WHERE question_text = 'What is the purpose of database normalization?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'To increase data redundancy', false),
((SELECT id FROM questions WHERE question_text = 'What is the purpose of database normalization?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'To reduce data redundancy and improve data integrity', true),
((SELECT id FROM questions WHERE question_text = 'What is the purpose of database normalization?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'To make queries slower', false),
((SELECT id FROM questions WHERE question_text = 'What is the purpose of database normalization?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'To increase storage space', false);

-- Insert question options for question 4: "Which SQL command is used to retrieve data from a database?"
INSERT INTO question_options (question_id, option_text, is_correct) VALUES
((SELECT id FROM questions WHERE question_text = 'Which SQL command is used to retrieve data from a database?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'INSERT', false),
((SELECT id FROM questions WHERE question_text = 'Which SQL command is used to retrieve data from a database?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'UPDATE', false),
((SELECT id FROM questions WHERE question_text = 'Which SQL command is used to retrieve data from a database?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'SELECT', true),
((SELECT id FROM questions WHERE question_text = 'Which SQL command is used to retrieve data from a database?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'DELETE', false);

-- Insert question options for question 5: "What is ACID in the context of database transactions?"
INSERT INTO question_options (question_id, option_text, is_correct) VALUES
((SELECT id FROM questions WHERE question_text = 'What is ACID in the context of database transactions?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'A set of database optimization techniques', false),
((SELECT id FROM questions WHERE question_text = 'What is ACID in the context of database transactions?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'Atomicity, Consistency, Isolation, Durability', true),
((SELECT id FROM questions WHERE question_text = 'What is ACID in the context of database transactions?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'A type of database index', false),
((SELECT id FROM questions WHERE question_text = 'What is ACID in the context of database transactions?' AND topic_id = (SELECT id FROM topics WHERE name = 'Database' LIMIT 1) LIMIT 1), 'Advanced Cognitive Information Database', false);
