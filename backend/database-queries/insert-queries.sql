INSERT INTO users (name, designation, experience, skills, email, password) 
VALUES ('John Doe', 'Developer', 5, 'Go,Python', 'john@example.com', 'hashed_password');


INSERT INTO proposals (title, description, coordinator_id, deadline) 
VALUES ('Proposal Title', 'Proposal Description', 1, '2024-10-01');


INSERT INTO subtasks (proposal_id, topic_name, description, skills_required, submission_date) 
VALUES (1, 'Subtask Topic', 'Subtask Description', 'Go,HTML', '2024-09-30');


UPDATE submissions 
SET review_status = 'approved', reviewer_id = 1, review_date = NOW() 
WHERE id = 1;


