package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"go-playground/server/domain"
	"log"
)

// EventLogRepository struct
type EventLogRepository struct {
	db *sql.DB
}

// NewEventLogRepository creates a new EventLogRepository
func NewEventLogRepository(db *sql.DB) *EventLogRepository {
	return &EventLogRepository{db: db}
}

// Create inserts a new event log entry
func (r *EventLogRepository) Create(ctx context.Context, eventLog *domain.EventLog) error {

	jsonDetails, err := json.Marshal(eventLog.Details)
	if err != nil {
		return err
	}

	log.Println("eventLog: ", eventLog)

	query := `
		INSERT INTO event_log (event_type, actor_id, actor_type, details, reference_id, event_timestamp)
		VALUES ($1::event_type, $2, $3, $4::jsonb, $5, CURRENT_TIMESTAMP)
		RETURNING id, event_timestamp, created_at
	`

	err = r.db.QueryRow(
		query,
		eventLog.EventType,
		eventLog.ActorID,
		eventLog.ActorType,
		jsonDetails, // Pass marshaled JSON instead of map
		eventLog.ReferenceID,
	).Scan(&eventLog.ID, &eventLog.EventTimestamp, &eventLog.CreatedAt)

	if err != nil {
		log.Printf("Error creating event log: %v", err)
		return err
	}

	return nil
}

// GetByID retrieves an event log entry by its ID
func (r *EventLogRepository) GetByID(id string) (*domain.EventLog, error) {
	eventLog := &domain.EventLog{}
	query := `
		SELECT id, event_type, actor_id, actor_type, details, reference_id, event_timestamp, created_at
		FROM event_log
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&eventLog.ID,
		&eventLog.EventType,
		&eventLog.ActorID,
		&eventLog.ActorType,
		&eventLog.Details,
		&eventLog.ReferenceID,
		&eventLog.EventTimestamp,
		&eventLog.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return eventLog, err
}

// GetByUserID retrieves all event logs for a specific user
func (r *EventLogRepository) GetByUserID(userID string) ([]domain.EventLog, error) {
	query := `
		SELECT id, event_type, actor_id, actor_type, details, reference_id, event_timestamp, created_at
		FROM event_log
		WHERE actor_id = $1
		ORDER BY event_timestamp DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventLogs []domain.EventLog
	for rows.Next() {
		var eventLog domain.EventLog
		err := rows.Scan(
			&eventLog.ID,
			&eventLog.EventType,
			&eventLog.ActorID,
			&eventLog.ActorType,
			&eventLog.Details,
			&eventLog.ReferenceID,
			&eventLog.EventTimestamp,
			&eventLog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		eventLogs = append(eventLogs, eventLog)
	}
	return eventLogs, nil
}

// GetByReferenceID retrieves an event log entry by its reference ID
func (r *EventLogRepository) GetByReferenceID(referenceID string) (*domain.EventLog, error) {
	eventLog := &domain.EventLog{}
	query := `
		SELECT id, event_type, actor_id, actor_type, details, reference_id, event_timestamp, created_at
		FROM event_log
		WHERE reference_id = $1
	`
	err := r.db.QueryRow(query, referenceID).Scan(
		&eventLog.ID,
		&eventLog.EventType,
		&eventLog.ActorID,
		&eventLog.ActorType,
		&eventLog.Details,
		&eventLog.ReferenceID,
		&eventLog.EventTimestamp,
		&eventLog.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return eventLog, err
}

// Add other methods for EventLogRepository as needed
