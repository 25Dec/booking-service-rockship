package model

// Appointment -.
type Appointment struct {
	LearnerID  string      `json:"learner_id" gorm:"primaryKey;column:learner_id"`
	ScheduleID string      `json:"schedule_id" gorm:"primaryKey;column:schedule_id"`
	Schedule   Schedule    `json:"schedule" gorm:"-"` // Omit from GORM processing
	Detail     interface{} `json:"detail" gorm:"column:detail;type:jsonb"`
	WeekNumber int         `json:"week_number" gorm:"column:week_number"`
	Year       int         `json:"year" gorm:"column:year"`
	CreatedAt  int64       `json:"-" gorm:"column:created_at"`
	Content    string      `json:"content" gorm:"column:content"`
}
