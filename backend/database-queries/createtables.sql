CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    designation VARCHAR(255),
    experience INT,
    skills TEXT,  -- Comma-separated skills
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL  -- Use hashed passwords
);


CREATE TABLE proposals (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    coordinator_id INT,  -- Foreign key to users table
    deadline DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coordinator_id) REFERENCES users(id)
);


CREATE TABLE subtasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    proposal_id INT,  -- Foreign key to proposals table
    topic_name VARCHAR(255) NOT NULL,
    description TEXT,
    skills_required TEXT,  -- Comma-separated skills
    submission_date DATE,
    status ENUM('pending', 'completed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (proposal_id) REFERENCES proposals(id)
);


CREATE TABLE submissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    subtask_id INT,  -- Foreign key to subtasks table
    user_id INT,  -- Foreign key to users table (who uploaded the task)
    file_path VARCHAR(255) NOT NULL,
    submission_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    review_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    reviewer_id INT,  -- Foreign key to users table
    review_date TIMESTAMP,
    FOREIGN KEY (subtask_id) REFERENCES subtasks(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (reviewer_id) REFERENCES users(id)
);


