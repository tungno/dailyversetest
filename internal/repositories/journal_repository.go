/**
 *  JournalRepository defines the interface for data access operations related to journals.
 *  It abstracts the database layer, allowing CRUD operations (Create, Read, Update, Delete)
 *  on journal entries associated with specific users.
 *
 *  @interface JournalRepository
 *  @inherits None
 *
 *  @methods
 *  - CreateJournal(ctx, journal)                - Adds a new journal entry to the database.
 *  - GetJournal(ctx, userEmail, journalID)      - Retrieves a specific journal entry by its ID and user email.
 *  - UpdateJournal(ctx, journal)                - Updates an existing journal entry in the database.
 *  - DeleteJournal(ctx, userEmail, journalID)   - Deletes a journal entry by its ID and user email.
 *  - GetAllJournals(ctx, userEmail)             - Retrieves all journal entries associated with a specific user.
 *
 *  @dependencies
 *  - models.Journal: Defines the structure of a journal object.
 *  - context.Context: Manages request-scoped values, deadlines, and cancellations.
 *
 *  @file      journal_repository.go
 *  @project   DailyVerse
 *  @framework Go Interface for Repository Pattern
 *  @purpose   Database operations abstraction for journals.
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

// JournalRepository defines the interface for journal-related data operations.
type JournalRepository interface {
	// CreateJournal inserts a new journal entry into the database.
	CreateJournal(ctx context.Context, journal *models.Journal) error

	// GetJournal retrieves a specific journal entry by its ID and associated user email.
	GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error)

	// UpdateJournal modifies an existing journal entry in the database.
	UpdateJournal(ctx context.Context, journal *models.Journal) error

	// DeleteJournal removes a journal entry from the database by its ID and associated user email.
	DeleteJournal(ctx context.Context, userEmail, journalID string) error

	// GetAllJournals fetches all journal entries linked to a specific user's email.
	GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error)
}
