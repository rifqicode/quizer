-- Remove seeded topics
DELETE FROM topics WHERE name IN (
    'Database',
    'Programming', 
    'Web Development',
    'Data Structures',
    'System Design',
    'Security',
    'DevOps',
    'Machine Learning'
);
