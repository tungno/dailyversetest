/**
 *  TimetableService provides functionalities for managing and importing timetables.
 *  It parses ICS (iCalendar) content to extract events and saves them into the system using
 *  the EventRepository.
 *
 *  @file       timetable_service.go
 *  @package    services
 *
 *  @interfaces
 *  - TimetableServiceInterface - Defines the contract for timetable-related operations.
 *
 *  @methods
 *  - NewTimetableService(eventRepo)                   - Creates a new instance of TimetableService.
 *  - ImportTimetable(ctx, userEmail, icsContent)      - Parses and imports events from ICS content.
 *
 *  @dependencies
 *  - EventRepository: Handles CRUD operations for events.
 *  - "github.com/arran4/golang-ical": Provides ICS parsing capabilities.
 *  - models.Event: Represents the data structure for an event.
 *
 *  @behaviors
 *  - Parses ICS (iCalendar) content to extract event details such as title, description, location, and timing.
 *  - Saves each extracted event into the database using the EventRepository.
 *  - Ignores events with missing or invalid start and end times.
 *
 *  @example
 *  Import Timetable:
 *  ```
 *  timetableService := NewTimetableService(eventRepo)
 *  err := timetableService.ImportTimetable(ctx, "user@example.com", icsContent)
 *  if err != nil {
 *      log.Fatal("Failed to import timetable:", err)
 *  }
 *  ```
 *
 *  @errors
 *  - Returns an error if the ICS content cannot be parsed.
 *  - Returns an error if saving an event to the repository fails.
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
	ics "github.com/arran4/golang-ical"
	"strings"
	"time"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
)

// TimetableServiceInterface defines the operations for managing timetables.
type TimetableServiceInterface interface {
	// ImportTimetable parses ICS content and imports events for a specific user.
	ImportTimetable(ctx context.Context, userEmail, icsContent string) error
}

// TimetableService provides implementation of TimetableServiceInterface.
type TimetableService struct {
	EventRepo repositories.EventRepository // Repository for event data operations.
}

// NewTimetableService initializes a new instance of TimetableService.
func NewTimetableService(eventRepo repositories.EventRepository) TimetableServiceInterface {
	return &TimetableService{
		EventRepo: eventRepo,
	}
}

// ImportTimetable parses ICS content and saves the extracted events to the database.
// Parameters:
//   - ctx: The context for handling deadlines and cancellations.
//   - userEmail: The email of the user for whom the timetable is being imported.
//   - icsContent: The raw ICS content to be parsed.
//
// Returns:
//   - error: Returns an error if parsing or saving fails.
func (ts *TimetableService) ImportTimetable(ctx context.Context, userEmail, icsContent string) error {
	// Parse the ICS content.
	cal, err := ics.ParseCalendar(strings.NewReader(icsContent))
	if err != nil {
		return fmt.Errorf("Failed to parse ICS content: %v", err)
	}

	// Iterate over the events in the calendar.
	for _, event := range cal.Events() {
		// Extract event details.
		summary := event.GetProperty(ics.ComponentPropertySummary).Value
		description := event.GetProperty(ics.ComponentPropertyDescription).Value
		location := event.GetProperty(ics.ComponentPropertyLocation).Value

		dtStartProp := event.GetProperty(ics.ComponentPropertyDtStart)
		dtEndProp := event.GetProperty(ics.ComponentPropertyDtEnd)

		// Skip events with missing start or end time.
		if dtStartProp == nil || dtEndProp == nil {
			continue
		}

		// Parse start and end times.
		dtStart, err := time.Parse(time.RFC3339, dtStartProp.Value)
		if err != nil {
			continue
		}

		dtEnd, err := time.Parse(time.RFC3339, dtEndProp.Value)
		if err != nil {
			continue
		}

		// Create an event model.
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

		// Save the event to the repository.
		if err := ts.EventRepo.CreateEvent(ctx, &newEvent); err != nil {
			return fmt.Errorf("Failed to save event: %v", err)
		}
	}

	return nil
}
