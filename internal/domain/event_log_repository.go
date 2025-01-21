package domain

type EventLogRepository interface {
	Create(eventLog *EventLog) error
	GetByID(id string) (*EventLog, error)
	GetByUserID(userID string) ([]EventLog, error)
	Update(eventLog *EventLog) error
	Delete(id string) error
	GetByReferenceID(referenceID string) (*EventLog, error)
}
