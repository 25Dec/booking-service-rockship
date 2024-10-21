package repositories

import (
	"booking-service/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// AppointmentRepo -.
type appointmentRepositoryImpl struct {
	db *gorm.DB
}

// New -.
func NewAppointmentRepositoryImpl(db *gorm.DB) AppointmentRepository {
	return &appointmentRepositoryImpl{
		db: db,
	}
}

// GetByID -.
func (r *appointmentRepositoryImpl) GetByScheduleID(ctx context.Context, scheduleID string) (model.Appointment, error) {
	var appointment model.Appointment
	err := r.db.First(&appointment, "schedule_id = ?", scheduleID).Error
	if err != nil {
		return model.Appointment{}, fmt.Errorf("AppointmentRepo - GetByScheduleID - r.db.First: %w", err)
	}
	return appointment, nil
}

// GetManyLearnAppointmentsInWeek -.
func (r *appointmentRepositoryImpl) GetManyLearnAppointmentsInWeek(ctx context.Context, year int, week int, learnerID string) ([]model.Appointment, error) {
	var appointments []model.Appointment
	err := r.db.Find(&appointments, "year = ? AND week_number = ? AND learner_id = ?", year, week, learnerID).Error
	if err != nil {
		return nil, fmt.Errorf("AppointmentRepo - GetManyLearnAppointmentsInWeek - r.db.Find: %w", err)
	}

	return appointments, nil
}

// GetAll -.
func (r *appointmentRepositoryImpl) GetMany(ctx context.Context, from string, to string, limit int, offset int) ([]model.Appointment, int64, error) {
	var appointments []model.Appointment
	var count int64
	err := r.db.Table("appointments").
		Select("*").
		Joins("JOIN schedules ON schedules.id = appointments.schedule_id").
		Where("schedules.start_at >= ? AND schedules.start_at <= ?", from, to).Limit(limit).Offset(offset).Find(&appointments).Error
	if err != nil {
		return nil, 0, fmt.Errorf("AppointmentRepo - GetMany - r.db.Find: %w", err)
	}

	err = r.db.Table("appointments").
		Select("*").
		Joins("JOIN schedules ON schedules.id = appointments.schedule_id").
		Where("schedules.start_at >= ? AND schedules.start_at <= ?", from, to).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("AppointmentRepo - GetMany - r.db.Count: %w", err)
	}
	return appointments, count, nil
}

// Create -.
func (r *appointmentRepositoryImpl) Create(ctx context.Context, appointment model.Appointment) (model.Appointment, error) {
	err := r.db.Create(&appointment).Error
	if err != nil {
		return model.Appointment{}, fmt.Errorf("AppointmentRepo - Create - r.db.Create: %w", err)
	}
	return appointment, nil
}

// Update -.
func (r *appointmentRepositoryImpl) Update(ctx context.Context, appointment model.Appointment) (model.Appointment, error) {
	err := r.db.Save(&appointment).Error
	if err != nil {
		return model.Appointment{}, fmt.Errorf("AppointmentRepo - Update - r.db.Save: %w", err)
	}
	return appointment, nil
}

// Delete -.
func (r *appointmentRepositoryImpl) DeleteByScheduleID(ctx context.Context, scheduleID string) error {
	err := r.db.Delete(&model.Appointment{}, "schedule_id = ?", scheduleID).Error
	if err != nil {
		return fmt.Errorf("AppointmentRepo - Delete - r.db.Delete: %w", err)
	}
	return nil
}
