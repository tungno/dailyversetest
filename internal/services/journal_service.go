/**
 *  JournalService provides business logic for managing journal entries.
 *  It handles input validation, date formatting, and delegates data persistence
 *  operations to the JournalRepository.
 *
 *  @interface JournalServiceInterface
 *  @struct   JournalService
 *
 *  @methods
 *  - CreateJournal(ctx, journal)                - Creates a new journal entry after validation and formatting.
 *  - GetJournal(ctx, userEmail, journalID)      - Retrieves a specific journal entry by user email and journal ID.
 *  - UpdateJournal(ctx, journal)                - Updates an existing journal entry.
 *  - DeleteJournal(ctx, userEmail, journalID)   - Deletes a journal entry by its ID.
 *  - GetAllJournals(ctx, userEmail)             - Fetches all journal entries associated with a specific user.
 *
 *  @dependencies
 *  - repositories.JournalRepository: Interface for data persistence operations.
 *  - models.Journal: Defines the structure of a journal entry.
 *  - time.Parse: Used for validating and formatting date strings.
 *
 *  @file      journal_service.go
 *  @project   DailyVerse
 *  @framework Go Business Logic Layer
 *  @purpose   Handles journal-related business logic and interacts with the repository layer.
 *
 *  @example
 *  ```
 *  journal := &models.Journal{
 *      Email: "user@example.com",
 *      Date:  "2024-12-05",
 *      Content: "This is a sample journal entry.",
 *  }
 *
 *  err := journalService.CreateJournal(context.Background(), journal)
 *  if err != nil {
 *      log.Fatalf("Failed to create journal: %v", err)
 *  }
 *  ```
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package services

import (
	"context"
	"fmt"
	"time"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
)

// JournalServiceInterface defines the contract for journal services.
type JournalServiceInterface interface {
	// CreateJournal creates a new journal entry.
	CreateJournal(ctx context.Context, journal *models.Journal) error

	// GetJournal retrieves a specific journal entry by user email and journal ID.
	GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error)

	// UpdateJournal updates an existing journal entry.
	UpdateJournal(ctx context.Context, journal *models.Journal) error

	// DeleteJournal deletes a journal entry by its ID and user email.
	DeleteJournal(ctx context.Context, userEmail, journalID string) error

	// GetAllJournals fetches all journal entries for a specific user.
	GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error)
}

// JournalService implements JournalServiceInterface.
type JournalService struct {
	JournalRepo repositories.JournalRepository // Repository for journal data persistence.
}

// NewJournalService initializes a new JournalService instance.
func NewJournalService(journalRepo repositories.JournalRepository) JournalServiceInterface {
	return &JournalService{JournalRepo: journalRepo}
}

// CreateJournal validates and creates a new journal entry.
// Validates the date format (YYYY-MM-DD) and stores the journal in the repository.
func (js *JournalService) CreateJournal(ctx context.Context, journal *models.Journal) error {
	// Validate and format the journal's date.
	journalDate, err := time.Parse("2006-01-02", journal.Date)
	if err != nil {
		return fmt.Errorf("Invalid date format. Please use YYYY-MM-DD.")
	}
	journal.Date = journalDate.Format("2006-01-02")

	// Delegate creation to the repository.
	return js.JournalRepo.CreateJournal(ctx, journal)
}

// GetJournal retrieves a specific journal entry by user email and journal ID.
func (js *JournalService) GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error) {
	return js.JournalRepo.GetJournal(ctx, userEmail, journalID)
}

// UpdateJournal updates an existing journal entry.
func (js *JournalService) UpdateJournal(ctx context.Context, journal *models.Journal) error {
	return js.JournalRepo.UpdateJournal(ctx, journal)
}

// DeleteJournal deletes a journal entry by its ID and associated user email.
func (js *JournalService) DeleteJournal(ctx context.Context, userEmail, journalID string) error {
	return js.JournalRepo.DeleteJournal(ctx, userEmail, journalID)
}

// GetAllJournals fetches all journal entries associated with a specific user.
func (js *JournalService) GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error) {
	return js.JournalRepo.GetAllJournals(ctx, userEmail)
}
