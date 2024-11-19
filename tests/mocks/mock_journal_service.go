// tests/mocks/mock_journal_service.go
package mocks

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
)

type MockJournalService struct {
	Journals map[string]*models.Journal
}

func NewMockJournalService() *MockJournalService {
	return &MockJournalService{
		Journals: make(map[string]*models.Journal),
	}
}

func (mjs *MockJournalService) CreateJournal(ctx context.Context, journal *models.Journal) error {
	if _, exists := mjs.Journals[journal.JournalID]; exists {
		return fmt.Errorf("journal already exists")
	}
	mjs.Journals[journal.JournalID] = journal
	return nil
}

func (mjs *MockJournalService) GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error) {
	journal, exists := mjs.Journals[journalID]
	if !exists || journal.Email != userEmail {
		return nil, fmt.Errorf("journal not found")
	}
	return journal, nil
}

func (mjs *MockJournalService) UpdateJournal(ctx context.Context, journal *models.Journal) error {
	existingJournal, exists := mjs.Journals[journal.JournalID]
	if !exists || existingJournal.Email != journal.Email {
		return fmt.Errorf("journal not found")
	}
	mjs.Journals[journal.JournalID] = journal
	return nil
}

func (mjs *MockJournalService) DeleteJournal(ctx context.Context, userEmail, journalID string) error {
	journal, exists := mjs.Journals[journalID]
	if !exists || journal.Email != userEmail {
		return fmt.Errorf("journal not found")
	}
	delete(mjs.Journals, journalID)
	return nil
}

func (mjs *MockJournalService) GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error) {
	var journals []models.Journal
	for _, journal := range mjs.Journals {
		if journal.Email == userEmail {
			journals = append(journals, *journal)
		}
	}
	return journals, nil
}
