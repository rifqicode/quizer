-- Remove the sample database questions and their options
-- First remove question options, then questions

DELETE FROM question_options WHERE question_id IN (
    SELECT id FROM questions WHERE topic_id = (
        SELECT id FROM topics WHERE name = 'Database' LIMIT 1
    ) AND question_text IN (
        'What does SQL stand for?',
        'Which of the following is NOT a type of database constraint?',
        'What is the purpose of database normalization?',
        'Which SQL command is used to retrieve data from a database?',
        'What is ACID in the context of database transactions?'
    )
);

DELETE FROM questions WHERE topic_id = (
    SELECT id FROM topics WHERE name = 'Database' LIMIT 1
) AND question_text IN (
    'What does SQL stand for?',
    'Which of the following is NOT a type of database constraint?',
    'What is the purpose of database normalization?',
    'Which SQL command is used to retrieve data from a database?',
    'What is ACID in the context of database transactions?'
);
