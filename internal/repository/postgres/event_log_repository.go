package postgres

import (
	"database/sql"
	"go-playground/internal/domain"
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
func (r *EventLogRepository) Create(eventLog *domain.EventLog) error {
	query := `
		INSERT INTO event_log (event_type, user_id, details, reference_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, event_timestamp, created_at
	`
	return r.db.QueryRow(
		query,
		eventLog.EventType,
		eventLog.UserID,
		eventLog.Details,
		eventLog.ReferenceID,
	).Scan(&eventLog.ID, &eventLog.EventTimestamp, &eventLog.CreatedAt)
}

// GetByID retrieves an event log entry by its ID
func (r *EventLogRepository) GetByID(id string) (*domain.EventLog, error) {
	eventLog := &domain.EventLog{}
	query := `
		SELECT id, event_type, user_id, details, event_timestamp, created_at
		FROM event_log
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&eventLog.ID,
		&eventLog.EventType,
		&eventLog.UserID,
		&eventLog.Details,
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
		SELECT id, event_type, user_id, details, event_timestamp, created_at
		FROM event_log
		WHERE user_id = $1
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
			&eventLog.UserID,
			&eventLog.Details,
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

// Update modifies an existing event log entry
func (r *EventLogRepository) Update(eventLog *domain.EventLog) error {
	query := `
		UPDATE event_log
		SET event_type = $1, details = $2, reference_id = $3
		WHERE id = $4
		RETURNING updated_at
	`
	return r.db.QueryRow(query, eventLog.EventType, eventLog.Details, eventLog.ReferenceID, eventLog.ID).Scan(&eventLog.UpdatedAt)
}

// Delete removes an event log entry by its ID
func (r *EventLogRepository) Delete(id string) error {
	query := `DELETE FROM event_log WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetByReferenceID retrieves an event log entry by its reference ID
func (r *EventLogRepository) GetByReferenceID(referenceID string) (*domain.EventLog, error) {
	eventLog := &domain.EventLog{}
	query := `
		SELECT id, event_type, user_id, details, event_timestamp, created_at
		FROM event_log
		WHERE reference_id = $1
	`
	err := r.db.QueryRow(query, referenceID).Scan(
		&eventLog.ID,
		&eventLog.EventType,
		&eventLog.UserID,
		&eventLog.Details,
		&eventLog.EventTimestamp,
		&eventLog.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return eventLog, err
}

// Add other methods for EventLogRepository as needed
