package migrations

import (
	"database/sql"
	"fmt"
)

// RunMigrations will execute all database migrations
func RunMigrations(db *sql.DB) error {
	if err := createUsersTable(db); err != nil {
		return err
	}
	if err := createProposalsTable(db); err != nil {
		return err
	}
	if err := createSubtasksTable(db); err != nil {
		return err
	}
	if err := createSubmissionsTable(db); err != nil {
		return err
	}
	if err := createNotificationsTable(db); err != nil {
		return err
	}

	fmt.Println("Migrations completed successfully!")
	return nil
}

func createUsersTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        designation VARCHAR(255),
        experience INT,
        skills TEXT,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}
	fmt.Println("Users table created.")
	return nil
}

func createProposalsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS proposals (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        coordinator_id INT,
        deadline DATE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (coordinator_id) REFERENCES users(id)
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating proposals table: %v", err)
	}
	fmt.Println("Proposals table created.")
	return nil
}

func createSubtasksTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS subtasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        proposal_id INT,
        topic_name VARCHAR(255) NOT NULL,
        description TEXT,
        skills_required TEXT,
        submission_date DATE,
        status ENUM('pending', 'completed') DEFAULT 'pending',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (proposal_id) REFERENCES proposals(id)
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating subtasks table: %v", err)
	}
	fmt.Println("Subtasks table created.")
	return nil
}

func createSubmissionsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS submissions (
        id INT AUTO_INCREMENT PRIMARY KEY,
        subtask_id INT,
        user_id INT,
        file_path VARCHAR(255) NOT NULL,
        submission_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        review_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
        reviewer_id INT,
        review_date TIMESTAMP,
        FOREIGN KEY (subtask_id) REFERENCES subtasks(id),
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (reviewer_id) REFERENCES users(id)
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating submissions table: %v", err)
	}
	fmt.Println("Submissions table created.")
	return nil
}

func createNotificationsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS notifications (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT,
        proposal_id INT,
        subtask_id INT,
        message VARCHAR(255),
        status ENUM('unread', 'read') DEFAULT 'unread',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (proposal_id) REFERENCES proposals(id),
        FOREIGN KEY (subtask_id) REFERENCES subtasks(id)
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating notifications table: %v", err)
	}
	fmt.Println("notifications table created.")
	return nil
}
