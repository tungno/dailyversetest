// internal/services/journal_service.go
package services

import (
	"context"
	"fmt"
	"time"

	"proh2052-group6/pkg/models"

	"google.golang.org/api/iterator"
)

type JournalServiceInterface interface {
	CreateJournal(ctx context.Context, journal *models.Journal) error
	GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error)
	UpdateJournal(ctx context.Context, journal *models.Journal) error
	DeleteJournal(ctx context.Context, userEmail, journalID string) error
	GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error)
}

type JournalService struct {
	DB DatabaseInterface
}

func NewJournalService(db DatabaseInterface) JournalServiceInterface {
	return &JournalService{DB: db}
}

func (js *JournalService) CreateJournal(ctx context.Context, journal *models.Journal) error {
	// Parse and format the date
	journalDate, err := time.Parse("2006-01-02", journal.Date)
	if err != nil {
		return fmt.Errorf("Invalid date format. Please use YYYY-MM-DD.")
	}
	journal.Date = journalDate.Format("2006-01-02")

	userDocRef := js.DB.Collection("users").Doc(journal.Email).Collection("journals")
	docRef, _, err := userDocRef.Add(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to create journal")
	}

	journal.JournalID = docRef.ID
	_, err = docRef.Set(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to update journal with JournalID")
	}

	return nil
}

func (js *JournalService) GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error) {
	doc, err := js.DB.Collection("users").Doc(userEmail).Collection("journals").Doc(journalID).Get(ctx)
	if err != nil || !doc.Exists() {
		return nil, fmt.Errorf("Journal not found")
	}

	var journal models.Journal
	err = doc.DataTo(&journal)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse journal data")
	}

	return &journal, nil
}

func (js *JournalService) UpdateJournal(ctx context.Context, journal *models.Journal) error {
	docRef := js.DB.Collection("users").Doc(journal.Email).Collection("journals").Doc(journal.JournalID)
	_, err := docRef.Set(ctx, journal)
	if err != nil {
		return fmt.Errorf("Failed to update journal")
	}
	return nil
}

func (js *JournalService) DeleteJournal(ctx context.Context, userEmail, journalID string) error {
	docRef := js.DB.Collection("users").Doc(userEmail).Collection("journals").Doc(journalID)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete journal")
	}
	return nil
}

func (js *JournalService) GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error) {
	userDocRef := js.DB.Collection("users").Doc(userEmail).Collection("journals")
	iter := userDocRef.Documents(ctx)

	var journals []models.Journal

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve journals")
		}

		var journal models.Journal
		err = doc.DataTo(&journal)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse journal data")
		}

		journal.JournalID = doc.Ref.ID
		journals = append(journals, journal)
	}

	return journals, nil
}
