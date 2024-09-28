package handlers

import (
	"fmt"
	"kpmg-hackathon/db"
	"kpmg-hackathon/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SignupHandler handles user registration
func SignupHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create user: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	var creds models.LoginCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := db.AuthenticateUser(creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": userID})
}

// CreateProposalHandler handles the creation of proposals
func CreateProposalHandler(c *gin.Context) {
	var proposal models.Proposal
	if err := c.ShouldBindJSON(&proposal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.CreateProposal(proposal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create proposal: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proposal created successfully"})
}

func GetProposalsHandler(c *gin.Context) {
	proposals, err := db.GetAllProposals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve proposals: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, proposals)
}

// CreateSubtaskHandler handles subtask creation
func CreateSubtaskHandler(c *gin.Context) {
	var subtask models.Subtask
	if err := c.ShouldBindJSON(&subtask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.CreateSubtask(subtask); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create subtask: %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtask created successfully"})
}

func GetSubtasksHandler(c *gin.Context) {
	proposalIDParam := c.Param("proposal_id") // Get the proposal ID from the URL
	proposalID, err := strconv.Atoi(proposalIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}

	subtasks, err := db.GetSubtasksByProposalID(proposalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subtasks: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, subtasks)
}

// UploadSubtaskHandler handles the uploading of subtasks
func UploadSubtaskHandler(c *gin.Context) {
	// Parse form data
	subtaskIDStr := c.PostForm("subtask_id")
	userIDStr := c.PostForm("user_id")

	// Convert subtask_id and user_id from string to int
	subtaskID, err := strconv.Atoi(subtaskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subtask ID"})
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// File handling
	file, err := c.FormFile("file_path")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File upload failed: %v", err.Error())})
		return
	}

	// Save the file to a directory on the server
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	// Create submission record
	submission := models.Submission{
		SubtaskID:      subtaskID,
		UserID:         userID,
		FilePath:       filePath,
		SubmissionDate: time.Now(),
	}

	// Save submission to database
	if err := db.UploadSubtask(submission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload subtask"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtask uploaded successfully"})
}

func DownloadSubmissionFile(c *gin.Context) {
	// Get the submission ID from the request URL
	submissionID := c.Param("submission_id")

	// Retrieve the submission from the database using the submissionID
	submission, err := db.GetSubmissionByID(submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submission: " + err.Error()})
		return
	}

	// Check if the file exists before sending it
	if _, err := os.Stat(submission.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file to the client
	c.File(submission.FilePath)
}

// ReviewSubtaskHandler handles the review of subtasks
func ReviewSubtaskHandler(c *gin.Context) {
	var review models.SubmissionReview
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.ReviewDate = time.Now()

	if err := db.ReviewSubtask(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review subtask"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtask reviewed successfully"})
}

// GetProposalsWithApprovedSubtasks handles the tracking and reporting request
func GetProposalsWithApprovedSubtasks(c *gin.Context) {
	userID := c.Param("userId")

	proposals, err := db.ProposalsWithApprovedSubtasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to review subtask"})
		return
	}

	c.JSON(http.StatusOK, proposals)
}

func GetNotifications(c *gin.Context) {
	userID := c.Param("user_id") // Assume user_id is passed as a URL parameter

	// Call the function to get notifications
	notifications, err := db.GetNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func GetReviewSubmissionsList(c *gin.Context) {
	userID := c.Param("user_id")

	notifications, err := db.GetReviewSubmissionsList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
