package services

import (
	"booking-service/internal/model"
	"booking-service/internal/repositories"
	"context"
	"fmt"
)

// Appointment -.
type appointmentServiceImpl struct {
	AppointmentRepo repositories.AppointmentRepository
}

// NewAppointmentService -.
func NewAppointmentService(appointmentRepo repositories.AppointmentRepository) AppointmentService {
	return &appointmentServiceImpl{
		AppointmentRepo: appointmentRepo,
	}
}

func (s *appointmentServiceImpl) GetAppointmentByScheduleID(ctx context.Context, scheduleID string) (model.Appointment, error) {
	appointment, err := s.AppointmentRepo.GetByScheduleID(ctx, scheduleID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("AppointmentServiceImpl - GetAppointmentByScheduleID - s.AppointmentRepo.GetByScheduleID: %w", err)
	}
	return appointment, nil
}

func (s *appointmentServiceImpl) GetLearnerAppointmentsInWeek(ctx context.Context, year int, week int, learnerID string) ([]model.Appointment, error) {
	appointments, err := s.AppointmentRepo.GetManyLearnAppointmentsInWeek(ctx, year, week, learnerID)
	if err != nil {
		return nil, fmt.Errorf("AppointmentServiceImpl - GetLearnerAppointmentsInWeek - s.AppointmentRepo.GetManyByLearnerID: %w", err)
	}
	return appointments, nil
}

// GetAppointments -.
func (s *appointmentServiceImpl) GetAppointments(ctx context.Context, from string, to string, limit int, offset int) ([]model.Appointment, int64, error) {
	appointments, count, err := s.AppointmentRepo.GetMany(ctx, from, to, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("AppointmentServiceImpl - GetAppointments - s.AppointmentRepo.GetAll: %w", err)
	}
	return appointments, count, nil
}

// CreateAppointment -.
func (s *appointmentServiceImpl) CreateAppointment(ctx context.Context, appointment model.Appointment) (model.Appointment, error) {
	appointment, err := s.AppointmentRepo.Create(ctx, appointment)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("AppointmentServiceImpl - CreateAppointment - s.AppointmentRepo.Create: %w", err)
	}
	return appointment, nil
}

// UpdateAppointment -.
func (s *appointmentServiceImpl) UpdateAppointment(ctx context.Context, appointment model.Appointment) (model.Appointment, error) {
	appointment, err := s.AppointmentRepo.Update(ctx, appointment)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("AppointmentServiceImpl - UpdateAppointment - s.AppointmentRepo.Update: %w", err)
	}
	return appointment, nil
}

// DeleteAppointment -.
func (s *appointmentServiceImpl) DeleteAppointmentByScheduleID(ctx context.Context, id string) error {
	err := s.AppointmentRepo.DeleteByScheduleID(ctx, id)
	if err != nil {
		return fmt.Errorf("AppointmentServiceImpl - DeleteAppointment - s.AppointmentRepo.Delete: %w", err)
	}
	return nil
}
