package model

type Schedule struct {
	ID       string `json:"id"`
	MentorID string `json:"mentor_id"`
	Interval struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"interval" gorm:"-"` //this field does not store in db
	Wday      string `json:"wday"  gorm:"-"` //this field does not store in db
	StartAt   string `json:"start_at"`       //this field does store in db can be used to generate intervals and wday
	CreatedAt int64  `json:"-" gorm:"created_at"`
}

type CreateScheduleRequest struct {
	MentorID string `json:"mentor_id" binding:"required"`
	StartAt  string `json:"start_at" binding:"required" time_format:"2006-01-02 15:04:05+07"`
}

type ScheduleAppointmentRequest struct {
	ScheduleAt string `json:"schedule_at" form:"start_at" time_format:"2006-01-02 15:04:05+07"`
	LearnerID  string `json:"learner_id" binding:"required"`
	Content    string `json:"content"`
}
