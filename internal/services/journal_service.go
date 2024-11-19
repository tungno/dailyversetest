// internal/services/journal_service.go
package services

import (
	"context"
	"fmt"
	"time"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
)

type JournalServiceInterface interface {
	CreateJournal(ctx context.Context, journal *models.Journal) error
	GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error)
	UpdateJournal(ctx context.Context, journal *models.Journal) error
	DeleteJournal(ctx context.Context, userEmail, journalID string) error
	GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error)
}

type JournalService struct {
	JournalRepo repositories.JournalRepository
}

func NewJournalService(journalRepo repositories.JournalRepository) JournalServiceInterface {
	return &JournalService{JournalRepo: journalRepo}
}

func (js *JournalService) CreateJournal(ctx context.Context, journal *models.Journal) error {
	// Parse and format the date
	journalDate, err := time.Parse("2006-01-02", journal.Date)
	if err != nil {
		return fmt.Errorf("Invalid date format. Please use YYYY-MM-DD.")
	}
	journal.Date = journalDate.Format("2006-01-02")

	return js.JournalRepo.CreateJournal(ctx, journal)
}

func (js *JournalService) GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error) {
	return js.JournalRepo.GetJournal(ctx, userEmail, journalID)
}

func (js *JournalService) UpdateJournal(ctx context.Context, journal *models.Journal) error {
	return js.JournalRepo.UpdateJournal(ctx, journal)
}

func (js *JournalService) DeleteJournal(ctx context.Context, userEmail, journalID string) error {
	return js.JournalRepo.DeleteJournal(ctx, userEmail, journalID)
}

func (js *JournalService) GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error) {
	return js.JournalRepo.GetAllJournals(ctx, userEmail)
}
