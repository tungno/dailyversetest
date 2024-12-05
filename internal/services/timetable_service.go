// internal/services/timetable_service.go
package services

import (
	"context"
	"fmt"
	ics "github.com/arran4/golang-ical"
	"strings"
	"time"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
)

type TimetableServiceInterface interface {
	ImportTimetable(ctx context.Context, userEmail, icsContent string) error
}

type TimetableService struct {
	EventRepo repositories.EventRepository
}

func NewTimetableService(eventRepo repositories.EventRepository) TimetableServiceInterface {
	return &TimetableService{
		EventRepo: eventRepo,
	}
}

func (ts *TimetableService) ImportTimetable(ctx context.Context, userEmail, icsContent string) error {
	// Parse the ICS content
	cal, err := ics.ParseCalendar(strings.NewReader(icsContent))
	if err != nil {
		return fmt.Errorf("Failed to parse ICS content: %v", err)
	}

	for _, event := range cal.Events() {
		// Extract event details
		summary := event.GetProperty(ics.ComponentPropertySummary).Value
		description := event.GetProperty(ics.ComponentPropertyDescription).Value
		location := event.GetProperty(ics.ComponentPropertyLocation).Value

		dtStartProp := event.GetProperty(ics.ComponentPropertyDtStart)
		dtEndProp := event.GetProperty(ics.ComponentPropertyDtEnd)

		if dtStartProp == nil || dtEndProp == nil {
			continue
		}

		dtStart, err := time.Parse(time.RFC3339, dtStartProp.Value)
		if err != nil {
			continue
		}

		dtEnd, err := time.Parse(time.RFC3339, dtEndProp.Value)
		if err != nil {
			continue
		}

		// Create event model
		newEvent := models.Event{
			Email:         userEmail,
			Title:         summary,
			Description:   description,
			Date:          dtStart.Format("2006-01-02"),
			StartTime:     dtStart.Format("15:04"),
			EndTime:       dtEnd.Format("15:04"),
			EventTypeID:   "private",
			Status:        "confirmed",
			StreetAddress: location,
		}

		// Save event to repository
		if err := ts.EventRepo.CreateEvent(ctx, &newEvent); err != nil {
			return fmt.Errorf("Failed to save event: %v", err)
		}
	}

	return nil
}
