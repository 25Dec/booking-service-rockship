package repositories

import (
	"booking-service/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type scheduleRepositoryImpl struct {
	db *gorm.DB
}

func NewScheduleRepoImpl(db *gorm.DB) ScheduleRepository {
	return &scheduleRepositoryImpl{db: db}
}

func (r *scheduleRepositoryImpl) GetManyByLearnerID(ctx context.Context, learnerID string) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Table("schedules").
		Select("schedules.*").
		Joins("JOIN appointments ON schedules.id = appointments.schedule_id").
		Where("appointments.learner_id = ?", learnerID).
		Find(&schedules).Error
	if err != nil {
		return nil, fmt.Errorf("ScheduleRepo - GetManyByLearnerID - r.db.Find: %w", err)
	}
	return schedules, nil
}

func (r *scheduleRepositoryImpl) GetAvailableSchedulesAtTime(ctx context.Context, scheduleAt string) ([]model.Schedule, error) {
	var schedules []model.Schedule

	err := r.db.
		Table("schedules").
		Select("schedules.*").
		Joins("JOIN mentors ON mentors.id = schedules.mentor_id").
		Joins("LEFT JOIN appointments ON schedules.id = appointments.schedule_id").
		Where("appointments.schedule_id IS NULL").
		Where("mentors.disabled = ?", false).
		Where("schedules.start_at = ?", scheduleAt).
		Find(&schedules).Error
	if err != nil {
		return nil, fmt.Errorf("ScheduleRepo - GetAvailableSchedulesAtTime - r.db.Find: %w", err)
	}

	return schedules, nil
}

func (r *scheduleRepositoryImpl) GetAvailableSchedules(ctx context.Context, from string, to string) ([]model.Schedule, error) {
	var schedules []model.Schedule

	err := r.db.
		Table("schedules").
		Select("schedules.*").
		Joins("JOIN mentors ON mentors.id = schedules.mentor_id").
		Joins("LEFT JOIN appointments ON schedules.id = appointments.schedule_id").
		Where("appointments.schedule_id IS NULL").
		Where("mentors.disabled = ?", false).
		Where("schedules.start_at >= ? AND schedules.start_at <= ?", from, to).
		Find(&schedules).Error

	if err != nil {
		return nil, fmt.Errorf("ScheduleRepo - GetAvailableSchedules - r.db.Find: %w", err)
	}
	return schedules, nil
}

func (r *scheduleRepositoryImpl) GetByID(ctx context.Context, id string) (model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.First(&schedule, "id = ?", id).Error
	if err != nil {
		return model.Schedule{}, fmt.Errorf("ScheduleRepo - GetByID - r.db.First: %w", err)
	}
	return schedule, nil
}

func (r *scheduleRepositoryImpl) GetAllByTime(ctx context.Context, scheduleAt string) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Where("start_at IN (?)", scheduleAt).Find(&schedules).Error
	if err != nil {
		return nil, fmt.Errorf("ScheduleRepo - GetAllByTime - r.db.Find: %w", err)
	}
	return schedules, nil
}

func (r *scheduleRepositoryImpl) GetMany(ctx context.Context, from string, to string, limit int, offset int) ([]model.Schedule, int64, error) {
	var schedules []model.Schedule
	var count int64
	err := r.db.Where("start_at >= ? AND start_at <= ?", from, to).Limit(limit).Offset(offset).Find(&schedules).Error
	if err != nil {
		return nil, 0, fmt.Errorf("ScheduleRepo - GetMany - r.db.Find: %w", err)
	}
	err = r.db.Where("start_at >= ? AND start_at <= ?", from, to).Model(&model.Schedule{}).Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("ScheduleRepo - GetMany - r.db.Count: %w", err)
	}
	return schedules, count, nil
}

func (r *scheduleRepositoryImpl) GetManyByMentorID(ctx context.Context, mentorID string) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Where("mentor_id = ?", mentorID).Find(&schedules).Error
	if err != nil {
		return nil, fmt.Errorf("ScheduleRepo - GetManyByMentorID - r.db.Find: %w", err)
	}
	return schedules, nil
}

func (r *scheduleRepositoryImpl) Create(ctx context.Context, schedule model.Schedule) (model.Schedule, error) {
	schedule.ID = uuid.New().String()
	schedule.CreatedAt = time.Now().Unix()
	err := r.db.Create(&schedule).Error
	if err != nil {
		return model.Schedule{}, fmt.Errorf("ScheduleRepo - Create - r.db.Create: %w", err)
	}
	return schedule, nil
}

func (r *scheduleRepositoryImpl) Update(ctx context.Context, schedule model.Schedule) (model.Schedule, error) {
	err := r.db.Save(&schedule).Error
	if err != nil {
		return model.Schedule{}, fmt.Errorf("ScheduleRepo - Update - r.db.Save: %w", err)
	}
	return schedule, nil
}

func (r *scheduleRepositoryImpl) Delete(ctx context.Context, id string) error {
	err := r.db.Delete(&model.Schedule{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("ScheduleRepo - Delete - r.db.Delete: %w", err)
	}
	return nil
}
