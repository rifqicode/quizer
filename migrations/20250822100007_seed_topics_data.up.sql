-- Insert default topics for testing
INSERT INTO topics (name, description, is_active) VALUES
('Database', 'Questions about database concepts, SQL, and database management systems', true),
('Programming', 'General programming concepts, algorithms, and best practices', true),
('Web Development', 'Frontend and backend web development technologies', true),
('Data Structures', 'Arrays, linked lists, trees, graphs, and other data structures', true),
('System Design', 'Scalability, architecture, and system design principles', true),
('Security', 'Cybersecurity, encryption, and secure coding practices', true),
('DevOps', 'CI/CD, containerization, cloud platforms, and infrastructure', true),
('Machine Learning', 'ML algorithms, data science, and artificial intelligence', true)
ON CONFLICT (name) DO NOTHING;
