package services

import (
	"context"

	"booking-service/internal/model"
)

type (
	MentorService interface {
		GetMentorByID(ctx context.Context, id string) (model.Mentor, error)
		GetMentors(ctx context.Context) ([]model.Mentor, error)
		CreatedMentor(ctx context.Context, mentor model.Mentor) (model.Mentor, error)
		UpdateMentor(ctx context.Context, mentor model.Mentor) (model.Mentor, error)
		DeleteMentor(ctx context.Context, id string) error
	}
)

type (
	// ScheduleService -.
	ScheduleService interface {
		GetAvailableSchedules(ctx context.Context, from string, to string) ([]model.Schedule, error)
		GetAppointmentOfSchedule(ctx context.Context, scheduleID string) (model.Appointment, error)
		ScheduleAppointment(ctx context.Context, learnerID string, scheduleAt string, content string) (model.Appointment, error)
		UnscheduleAppointment(ctx context.Context, scheduleID string) error
		GetSchedulesByLearnerID(ctx context.Context, learnerID string) ([]model.Schedule, error)
		GetScheduleByID(ctx context.Context, id string) (model.Schedule, error)
		GetSchedulesByMentorID(ctx context.Context, mentorID string) ([]model.Schedule, error)
		GetSchedulesByTime(ctx context.Context, scheduleAt string) ([]model.Schedule, error)
		GetSchedules(ctx context.Context, from string, to string, limit int, offset int) ([]model.Schedule, int64, error)
		CreateSchedule(ctx context.Context, schedule model.Schedule) (model.Schedule, error)
		UpdateSchedule(ctx context.Context, schedule model.Schedule) (model.Schedule, error)
		DeleteSchedule(ctx context.Context, id string) error
		GetScheduledAppointments(ctx context.Context, from string, to string, limit int, offset int) ([]model.Appointment, int64, error)
	}
)

type (
	// AppointmentService -.
	AppointmentService interface {
		// GetByScheduleID -.
		// Erors:
		// 	- ErrNotFound
		// 	- ErrInternal
		GetLearnerAppointmentsInWeek(ctx context.Context, year int, week int, learnerID string) ([]model.Appointment, error)
		GetAppointmentByScheduleID(ctx context.Context, scheduleID string) (model.Appointment, error)
		GetAppointments(ctx context.Context, from string, to string, limit int, offset int) ([]model.Appointment, int64, error)
		CreateAppointment(ctx context.Context, appointment model.Appointment) (model.Appointment, error)
		UpdateAppointment(ctx context.Context, appointment model.Appointment) (model.Appointment, error)
		DeleteAppointmentByScheduleID(ctx context.Context, scheduleID string) error
	}
)
