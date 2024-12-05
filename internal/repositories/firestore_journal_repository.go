/**
 *  FirestoreJournalRepository implements the JournalRepository interface, providing
 *  methods for managing journal data in a Firestore database. It allows users to perform
 *  CRUD operations (Create, Read, Update, Delete) on journals stored under user-specific
 *  collections.
 *
 *  @struct   FirestoreJournalRepository
 *  @inherits None
 *
 *  @methods
 *  - NewFirestoreJournalRepository(client)          - Creates a new FirestoreJournalRepository instance.
 *  - CreateJournal(ctx, journal)                   - Adds a new journal to the user's collection.
 *  - GetJournal(ctx, userEmail, journalID)         - Retrieves a specific journal by its ID.
 *  - UpdateJournal(ctx, journal)                   - Updates an existing journal in Firestore.
 *  - DeleteJournal(ctx, userEmail, journalID)      - Deletes a journal by its ID.
 *  - GetAllJournals(ctx, userEmail)                - Retrieves all journals for a specific user.
 *
 *  @dependencies
 *  - cloud.google.com/go/firestore: Provides the Firestore client for database operations.
 *  - google.golang.org/api/iterator: Handles Firestore document iteration.
 *  - models.Journal: Defines the structure of a journal object.
 *
 *  @file      firestore_journal_repository.go
 *  @project   DailyVerse
 *  @framework Go with Firestore integration
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package repositories

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreJournalRepository provides Firestore-based implementation of JournalRepository.
type FirestoreJournalRepository struct {
	Client *firestore.Client // Firestore client for database operations.
}

// NewFirestoreJournalRepository initializes a new FirestoreJournalRepository instance.
func NewFirestoreJournalRepository(client *firestore.Client) JournalRepository {
	return &FirestoreJournalRepository{Client: client}
}

// CreateJournal adds a new journal to the user's Firestore collection.
func (jr *FirestoreJournalRepository) CreateJournal(ctx context.Context, journal *models.Journal) error {
	userDocRef := jr.Client.Collection("users").Doc(journal.Email).Collection("journals")

	// Add journal data to Firestore.
	docRef, _, err := userDocRef.Add(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to create journal: %v", err)
	}

	// Update the journal with its generated ID.
	journal.JournalID = docRef.ID
	_, err = docRef.Set(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to update journal with JournalID: %v", err)
	}

	return nil
}

// GetJournal retrieves a specific journal by its ID from Firestore.
func (jr *FirestoreJournalRepository) GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error) {
	docRef := jr.Client.Collection("users").Doc(userEmail).Collection("journals").Doc(journalID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("Journal not found: %v", err)
	}

	// Map Firestore data to a Journal model.
	var journal models.Journal
	err = doc.DataTo(&journal)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse journal data: %v", err)
	}

	return &journal, nil
}

// UpdateJournal updates an existing journal in the Firestore collection.
func (jr *FirestoreJournalRepository) UpdateJournal(ctx context.Context, journal *models.Journal) error {
	docRef := jr.Client.Collection("users").Doc(journal.Email).Collection("journals").Doc(journal.JournalID)
	_, err := docRef.Set(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to update journal: %v", err)
	}
	return nil
}

// DeleteJournal removes a journal from Firestore by its ID.
func (jr *FirestoreJournalRepository) DeleteJournal(ctx context.Context, userEmail, journalID string) error {
	docRef := jr.Client.Collection("users").Doc(userEmail).Collection("journals").Doc(journalID)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete journal: %v", err)
	}
	return nil
}

// GetAllJournals retrieves all journals for a specific user from Firestore.
func (jr *FirestoreJournalRepository) GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error) {
	userDocRef := jr.Client.Collection("users").Doc(userEmail).Collection("journals")
	iter := userDocRef.Documents(ctx)

	var journals []models.Journal

	// Iterate through documents and map them to Journal models.
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve journals: %v", err)
		}

		var journal models.Journal
		err = doc.DataTo(&journal)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse journal data: %v", err)
		}

		// Include the document ID in the journal.
		journal.JournalID = doc.Ref.ID
		journals = append(journals, journal)
	}

	return journals, nil
}
