package repositories

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreJournalRepository struct {
	Client *firestore.Client
}

func NewFirestoreJournalRepository(client *firestore.Client) JournalRepository {
	return &FirestoreJournalRepository{Client: client}
}

func (jr *FirestoreJournalRepository) CreateJournal(ctx context.Context, journal *models.Journal) error {
	userDocRef := jr.Client.Collection("users").Doc(journal.Email).Collection("journals")
	docRef, _, err := userDocRef.Add(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to create journal: %v", err)
	}

	journal.JournalID = docRef.ID
	_, err = docRef.Set(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to update journal with JournalID: %v", err)
	}

	return nil
}

func (jr *FirestoreJournalRepository) GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error) {
	docRef := jr.Client.Collection("users").Doc(userEmail).Collection("journals").Doc(journalID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("Journal not found: %v", err)
	}

	var journal models.Journal
	err = doc.DataTo(&journal)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse journal data: %v", err)
	}

	return &journal, nil
}

func (jr *FirestoreJournalRepository) UpdateJournal(ctx context.Context, journal *models.Journal) error {
	docRef := jr.Client.Collection("users").Doc(journal.Email).Collection("journals").Doc(journal.JournalID)
	_, err := docRef.Set(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to update journal: %v", err)
	}
	return nil
}

func (jr *FirestoreJournalRepository) DeleteJournal(ctx context.Context, userEmail, journalID string) error {
	docRef := jr.Client.Collection("users").Doc(userEmail).Collection("journals").Doc(journalID)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete journal: %v", err)
	}
	return nil
}

func (jr *FirestoreJournalRepository) GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error) {
	userDocRef := jr.Client.Collection("users").Doc(userEmail).Collection("journals")
	iter := userDocRef.Documents(ctx)

	var journals []models.Journal

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

		journal.JournalID = doc.Ref.ID
		journals = append(journals, journal)
	}

	return journals, nil
}
