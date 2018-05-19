package contracts

import "time"

type EventCreatedEvent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	LocationID string    `json:"location_id"`
	Start      time.Time `json:"start_date"`
	End        time.Time `json:"end_date"`
}

func (e *EventCreatedEvent) EventName() string {
	return "eventCreated"
}
