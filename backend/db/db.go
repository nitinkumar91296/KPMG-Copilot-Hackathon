package db

import (
	"database/sql"
	"errors"
	"fmt"
	"kpmg-hackathon/models"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Initialize database connection
func init() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/kpmg")
	if err != nil {
		panic(err)
	}
}

// CreateUser inserts a new user into the database
func CreateUser(user models.User) error {
	query := "INSERT INTO users (name, designation, experience, skills, email, password) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, user.Name, user.Designation, user.Experience, user.Skills, user.Email, user.Password)
	return err
}

// AuthenticateUser checks if the user's credentials are correct
func AuthenticateUser(creds models.LoginCredentials) (int, error) {
	query := "SELECT id FROM users WHERE email = ? AND password = ?"
	var userID int
	err := db.QueryRow(query, creds.Email, creds.Password).Scan(&userID)
	if err != nil {
		return 0, errors.New("invalid credentials")
	}
	return userID, nil
}

// CreateProposal inserts a new proposal into the database
func CreateProposal(proposal models.Proposal) error {
	query := "INSERT INTO proposals (title, description, coordinator_id, deadline) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, proposal.Title, proposal.Description, proposal.CoordinatorID, proposal.Deadline)
	return err
}

func GetAllProposals() ([]models.Proposal, error) {
	var proposals []models.Proposal

	query := "SELECT id, title, description, coordinator_id, deadline, created_at FROM proposals"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var proposal models.Proposal
		var createdAtStr string

		if err := rows.Scan(&proposal.ID, &proposal.Title, &proposal.Description, &proposal.CoordinatorID, &proposal.Deadline, &createdAtStr); err != nil {
			return nil, err
		}

		// Convert the string to time.Time
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}
		proposal.CreatedAt = createdAt // Assign the parsed time to the proposal

		proposals = append(proposals, proposal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return proposals, nil
}

// CreateSubtask inserts a new subtask into the database
func CreateSubtask(subtask models.Subtask) error {
	query := "INSERT INTO subtasks (proposal_id, topic_name, description, skills_required, submission_date) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, subtask.ProposalID, subtask.TopicName, subtask.Description, subtask.SkillsRequired, subtask.SubmissionDate)
	return err
}

func GetSubtasksByProposalID(proposalID int) ([]models.Subtask, error) {
	var subtasks []models.Subtask

	query := "SELECT id, proposal_id, topic_name, description, skills_required, submission_date, status FROM subtasks WHERE proposal_id = ?"
	rows, err := db.Query(query, proposalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subtask models.Subtask
		if err := rows.Scan(&subtask.ID, &subtask.ProposalID, &subtask.TopicName, &subtask.Description, &subtask.SkillsRequired, &subtask.SubmissionDate, &subtask.Status); err != nil {
			return nil, err
		}
		subtasks = append(subtasks, subtask)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subtasks, nil
}

// UploadSubtask inserts a subtask submission into the database
func UploadSubtask(submission models.Submission) error {
	query := "INSERT INTO submissions (subtask_id, user_id, file_path, submission_date) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, submission.SubtaskID, submission.UserID, submission.FilePath, submission.SubmissionDate)
	return err
}

func GetSubmissionByID(submissionID string) (models.Submission, error) {
	var submission models.Submission
	var submissiondate string
	query := "SELECT id, subtask_id, user_id, file_path, submission_date FROM submissions WHERE id = ?"
	err := db.QueryRow(query, submissionID).Scan(&submission.ID, &submission.SubtaskID, &submission.UserID, &submission.FilePath, &submissiondate)

	// Convert the string to time.Time
	submissiondt, err := time.Parse("2006-01-02 15:04:05", submissiondate)
	if err != nil {
		return submission, fmt.Errorf("error parsing created_at: %v", err)
	}
	submission.SubmissionDate = submissiondt // Assign the parsed time to the proposal

	return submission, err
}

// ReviewSubtask updates the status of a subtask submission
func ReviewSubtask(review models.SubmissionReview) error {
	query := "UPDATE submissions SET review_status = ?, reviewer_id = ?, review_date = ? WHERE id = ?"
	_, err := db.Exec(query, review.ReviewStatus, review.ReviewerID, review.ReviewDate, review.SubmissionID)
	return err
}

func ProposalsWithApprovedSubtasks(userID string) ([]models.ProposalReport, error) {
	// SQL query to fetch all proposals where the user is the coordinator and fetch all approved subtasks
	query := `
    SELECT 
        p.id as proposal_id, p.title, p.description, p.deadline,
        s.id as subtask_id, s.topic_name, s.description, s.skills_required, s.submission_date, s.status,
        sub.reviewer_id, sub.review_status
    FROM proposals p
    LEFT JOIN subtasks s ON p.id = s.proposal_id
    LEFT JOIN submissions sub ON s.id = sub.subtask_id
    WHERE p.coordinator_id = ? AND sub.review_status = 'approved';
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error executing query: ", err)
		return nil, err
	}
	defer rows.Close()

	proposalsMap := make(map[int]*models.ProposalReport)

	for rows.Next() {
		var proposalID, subtaskID, reviewerID int
		var title, proposalDescription, deadline, topicName, subtaskDescription, skillsRequired, submissionDate, status, reviewStatus string

		if err := rows.Scan(&proposalID, &title, &proposalDescription, &deadline, &subtaskID, &topicName, &subtaskDescription, &skillsRequired, &submissionDate, &status, &reviewerID, &reviewStatus); err != nil {
			log.Println("Error scanning rows: ", err)
			return nil, err
		}

		// Check if the proposal is already added in the map
		if _, exists := proposalsMap[proposalID]; !exists {
			proposalsMap[proposalID] = &models.ProposalReport{
				ProposalID:       proposalID,
				Title:            title,
				Description:      proposalDescription,
				Deadline:         deadline,
				ApprovedSubtasks: []models.SubtaskReport{},
			}
		}

		// Add the approved subtask to the current proposal
		subtask := models.SubtaskReport{
			SubtaskID:      subtaskID,
			TopicName:      topicName,
			Description:    subtaskDescription,
			SkillsRequired: skillsRequired,
			SubmissionDate: submissionDate,
			Status:         status,
			ReviewerID:     reviewerID,
			ReviewStatus:   reviewStatus,
		}

		proposalsMap[proposalID].ApprovedSubtasks = append(proposalsMap[proposalID].ApprovedSubtasks, subtask)
	}

	// Convert the proposals map to a slice
	proposals := make([]models.ProposalReport, 0, len(proposalsMap))
	for _, proposal := range proposalsMap {
		proposals = append(proposals, *proposal)
	}

	return proposals, nil
}

func GetNotifications(userID string) ([]string, error) {
	notifications := []string{}

	// Query for new proposals
	queryNewProposals := `SELECT title FROM proposals WHERE created_at > NOW() - INTERVAL 1 DAY`
	rows, err := db.Query(queryNewProposals)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		notifications = append(notifications, fmt.Sprintf("New Proposal Submitted: %s", title))
	}

	// Query for new subtasks
	queryNewSubtasks := `SELECT topic_name FROM subtasks WHERE created_at > NOW() - INTERVAL 1 DAY`
	rows, err = db.Query(queryNewSubtasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var topicName string
		if err := rows.Scan(&topicName); err != nil {
			return nil, err
		}
		notifications = append(notifications, fmt.Sprintf("New Subtask Created: %s", topicName))
	}

	// Query for reviewed subtasks
	queryReviewedSubtasks := `SELECT subtask_id FROM submissions WHERE user_id = ? AND review_status IN ('approved', 'rejected')`
	rows, err = db.Query(queryReviewedSubtasks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subtaskId string
		if err := rows.Scan(&subtaskId); err != nil {
			return nil, err
		}

		// Query for reviewed subtasks
		queryGetSubtaskName := `SELECT topic_name FROM subtasks WHERE id = ?`
		rows, err = db.Query(queryGetSubtaskName, subtaskId)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var topicName string
			if err := rows.Scan(&topicName); err != nil {
				return nil, err
			}
			notifications = append(notifications, fmt.Sprintf("Subtask Reviewed: %s", topicName))
		}
	}

	return notifications, nil
}

func GetReviewSubmissionsList(userId, subtaskId string) ([]models.SubmissionDetails, error) {
	log.Print(userId)
	query := `
        SELECT 
			s.id AS submission_id,
			s.user_id,
			u.name AS user_name,
			s.subtask_id,
			st.topic_name AS subtask_topic,
			p.id AS proposal_id,
			p.title AS proposal_title,
			p.coordinator_id,
			cu.name AS coordinator_name,
			s.file_path,
			s.submission_date,
			s.review_status
		FROM 
			submissions s
		JOIN 
			subtasks st ON s.subtask_id = st.id
		JOIN 
			proposals p ON st.proposal_id = p.id
		JOIN 
			users u ON s.user_id = u.id
		JOIN 
			users cu ON p.coordinator_id = cu.id
		WHERE 
			p.coordinator_id = ? AND s.subtask_id = ?
		ORDER BY 
			s.submission_date DESC;
    `

	rows, err := db.Query(query, userId, subtaskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []models.SubmissionDetails
	for rows.Next() {
		var submission models.SubmissionDetails
		if err := rows.Scan(
			&submission.SubmissionID,
			&submission.UserID,
			&submission.UserName,
			&submission.SubtaskID,
			&submission.SubtaskTopic,
			&submission.ProposalID,
			&submission.ProposalTitle,
			&submission.CoordinatorID,
			&submission.CoordinatorName,
			&submission.FilePath,
			&submission.SubmissionDate,
			&submission.ReviewStatus,
		); err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func UpdateReviewStatus(submissionId int, reviewStatus string) error {
	// Update the review status for the specified submission
	updateSubmissionQuery := `UPDATE submissions SET review_status = ? WHERE id = ?`
	_, err := db.Exec(updateSubmissionQuery, reviewStatus, submissionId)
	if err != nil {
		return fmt.Errorf("failed to update submission review status: %v", err)
	}

	// If the review status is approved, update the subtask status to 'completed'
	if reviewStatus == "approved" {
		var subtaskId int
		selectSubtaskQuery := `SELECT subtask_id FROM submissions WHERE id = ?`
		err = db.QueryRow(selectSubtaskQuery, submissionId).Scan(&subtaskId)
		if err != nil {
			return fmt.Errorf("failed to get subtask id: %v", err)
		}

		updateSubtaskQuery := `UPDATE subtasks SET status = 'completed' WHERE id = ?`
		_, err = db.Exec(updateSubtaskQuery, subtaskId)
		if err != nil {
			return fmt.Errorf("failed to update subtask status: %v", err)
		}
	}

	return nil
}
