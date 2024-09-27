package main

import (
	"database/sql"
	"kpmg-hackathon/handlers"
	"kpmg-hackathon/migrations"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/kpmg")
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Run migrations
	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	defer db.Close()

	r := gin.Default()

	// CORS middleware to avoid CORS errors
	r.Use(cors.Default())

	// API routes
	// Auth Routes
	r.POST("/signup", handlers.SignupHandler)
	r.POST("/login", handlers.LoginHandler)

	// Proposal Routes
	r.POST("/proposal", handlers.CreateProposalHandler)
	r.GET("/proposals", handlers.GetProposalsHandler)

	// Subtask Routes
	r.POST("/subtask", handlers.CreateSubtaskHandler)
	r.GET("/proposals/:proposal_id/subtasks", handlers.GetSubtasksHandler)
	r.POST("/subtask/upload", handlers.UploadSubtaskHandler)
	r.GET("/submission/download/:submission_id", handlers.DownloadSubmissionFile)
	r.POST("/subtask/review", handlers.ReviewSubtaskHandler)
	r.GET("/coordinator/proposals/:userId", handlers.GetProposalsWithApprovedSubtasks)

	r.Run(":8080") // Start the server on port 8080
}
