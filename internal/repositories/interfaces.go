package repositories

import (
	"booking-service/internal/model"
	"context"
)

type (

	// ScheduleRepository -.
	ScheduleRepository interface {
		GetByID(ctx context.Context, id string) (model.Schedule, error)
		GetAvailableSchedulesAtTime(ctx context.Context, scheduleAt string) ([]model.Schedule, error)
		GetAvailableSchedules(ctx context.Context, from string, to string) ([]model.Schedule, error)
		GetAllByTime(ctx context.Context, scheduleAt string) ([]model.Schedule, error)
		GetManyByLearnerID(ctx context.Context, learnerID string) ([]model.Schedule, error)
		GetManyByMentorID(ctx context.Context, mentorID string) ([]model.Schedule, error)
		GetMany(ctx context.Context, from string, to string, limit int, offset int) ([]model.Schedule, int64, error)
		Create(ctx context.Context, schedule model.Schedule) (model.Schedule, error)
		Update(ctx context.Context, schedule model.Schedule) (model.Schedule, error)
		Delete(ctx context.Context, id string) error
	}

	// MentorRepository -.
	MentorRepository interface {
		GetByID(ctx context.Context, id string) (model.Mentor, error)
		GetMany(ctx context.Context) ([]model.Mentor, error)
		Create(ctx context.Context, mentor model.Mentor) (model.Mentor, error)
		Update(ctx context.Context, mentor model.Mentor) (model.Mentor, error)
		Delete(ctx context.Context, id string) error
	}

	// AppointmentRepository -.
	AppointmentRepository interface {
		GetByScheduleID(ctx context.Context, scheduleID string) (model.Appointment, error)
		GetMany(ctx context.Context, from string, to string, limit int, offset int) ([]model.Appointment, int64, error)
		GetManyLearnAppointmentsInWeek(ctx context.Context, year int, week int, learnerID string) ([]model.Appointment, error)
		Create(ctx context.Context, appointment model.Appointment) (model.Appointment, error)
		Update(ctx context.Context, appointment model.Appointment) (model.Appointment, error)
		DeleteByScheduleID(ctx context.Context, scheduleID string) error
	}
)
