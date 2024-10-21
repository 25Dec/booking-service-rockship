package model

type Event struct {
	ID        string      `json:"id"` //PK
	Payload   interface{} `json:"payload"`
	EventType string      `json:"event_type"`
	Status    string      `json:"status"` //success or failded
	CreatedAt int64       `json:"created_at"`
	CreatedBy string      `json:"created_by"`
}

// event type enum
const (
	CreateAppointment = "create_appointment"
	UpdateAppointment = "update_appointment"
	DeleteAppointment = "delete_appointment"
	CreateSchedule    = "create_schedule"
	UpdateSchedule    = "update_schedule"
	DeleteSchedule    = "delete_schedule"
)
