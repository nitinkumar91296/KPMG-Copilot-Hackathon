package models

import "time"

// User struct represents a user in the system
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Experience  int    `json:"experience"`
	Skills      string `json:"skills"` // Comma-separated list
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// LoginCredentials struct for user login
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Proposal struct represents a proposal
type Proposal struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CoordinatorID int       `json:"coordinator_id"`
	Deadline      string    `json:"deadline"` // YYYY-MM-DD format
	CreatedAt     time.Time `json:"created_at"`
}

// Subtask struct represents a subtask
type Subtask struct {
	ID             int    `json:"id"`
	ProposalID     int    `json:"proposal_id"`
	TopicName      string `json:"topic_name"`
	Description    string `json:"description"`
	SkillsRequired string `json:"skills_required"` // Comma-separated list
	SubmissionDate string `json:"submission_date"` // YYYY-MM-DD format
	Status         string `json:"status"`
}

// Submission struct represents a subtask submission
type Submission struct {
	ID             int       `json:"id"`
	SubtaskID      int       `json:"subtask_id"`
	UserID         int       `json:"user_id"`
	FilePath       string    `json:"file_path"`
	SubmissionDate time.Time `json:"submission_date"`
}

// SubmissionReview struct for reviewing a submission
type SubmissionReview struct {
	SubmissionID int       `json:"submission_id"`
	ReviewerID   int       `json:"reviewer_id"`
	ReviewStatus string    `json:"review_status"`
	ReviewDate   time.Time `json:"review_date"`
}

type SubtaskReport struct {
	SubtaskID      int    `json:"subtask_id"`
	TopicName      string `json:"topic_name"`
	Description    string `json:"description"`
	SkillsRequired string `json:"skills_required"`
	SubmissionDate string `json:"submission_date"`
	Status         string `json:"status"`
	ReviewerID     int    `json:"reviewer_id"`
	ReviewStatus   string `json:"review_status"`
}

type ProposalReport struct {
	ProposalID       int             `json:"proposal_id"`
	Title            string          `json:"title"`
	Description      string          `json:"description"`
	Deadline         string          `json:"deadline"`
	ApprovedSubtasks []SubtaskReport `json:"approved_subtasks"`
}

type Notification struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	ProposalID uint      `json:"proposal_id,omitempty"`
	SubtaskID  uint      `json:"subtask_id,omitempty"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
 
type SubmissionDetails struct {
    SubmissionID   int    `json:"submission_id"`
    UserID         int    `json:"user_id"`
    UserName       string `json:"user_name"`
    SubtaskID      int    `json:"subtask_id"`
    SubtaskTopic   string `json:"subtask_topic"`
    ProposalID     int    `json:"proposal_id"`
    ProposalTitle  string `json:"proposal_title"`
    CoordinatorID  int    `json:"coordinator_id"`
    CoordinatorName string `json:"coordinator_name"`
    FilePath       string `json:"file_path"`
    SubmissionDate string `json:"submission_date"`
    ReviewStatus   string `json:"review_status"`
}
