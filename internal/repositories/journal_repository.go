package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

type JournalRepository interface {
	CreateJournal(ctx context.Context, journal *models.Journal) error
	GetJournal(ctx context.Context, userEmail, journalID string) (*models.Journal, error)
	UpdateJournal(ctx context.Context, journal *models.Journal) error
	DeleteJournal(ctx context.Context, userEmail, journalID string) error
	GetAllJournals(ctx context.Context, userEmail string) ([]models.Journal, error)
}
