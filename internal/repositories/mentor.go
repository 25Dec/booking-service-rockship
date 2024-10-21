package repositories

import (
	"booking-service/internal/model"
	"booking-service/pkg/utils/errs"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mentorRepositoryImpl struct {
	db *gorm.DB
}

func NewMentorRepositoryImpl(db *gorm.DB) MentorRepository {
	return &mentorRepositoryImpl{db: db}
}

func (r *mentorRepositoryImpl) GetByID(ctx context.Context, id string) (model.Mentor, error) {
	var mentor model.Mentor
	err := r.db.First(&mentor, "id = ?", id).Error
	if err != nil {
		return model.Mentor{}, fmt.Errorf("MentorRepo - GetByID - r.db.First: %w", err)
	}

	if mentor.Disabled {
		return model.Mentor{}, errs.BadRequestError{
			Message: "Mentor is disabled",
		}
	}

	return mentor, nil
}

func (r *mentorRepositoryImpl) GetMany(ctx context.Context) ([]model.Mentor, error) {
	var mentors []model.Mentor
	err := r.db.Find(&mentors, "disabled = ?", false).Error
	if err != nil {
		return nil, fmt.Errorf("MentorRepo - GetAll - r.db.Find: %w", err)
	}
	return mentors, nil
}

func (r *mentorRepositoryImpl) Create(ctx context.Context, mentor model.Mentor) (model.Mentor, error) {
	mentor.ID = uuid.New().String()
	mentor.CreatedAt = time.Now().Unix()
	err := r.db.Create(&mentor).Error
	if err != nil {
		return model.Mentor{}, fmt.Errorf("MentorRepo - Create - r.db.Create: %w", err)
	}
	return mentor, nil
}

func (r *mentorRepositoryImpl) Update(ctx context.Context, mentor model.Mentor) (model.Mentor, error) {
	err := r.db.Save(&mentor).Error
	if err != nil {
		return model.Mentor{}, fmt.Errorf("MentorRepo - Update - r.db.Save: %w", err)
	}
	return mentor, nil
}

func (r *mentorRepositoryImpl) Delete(ctx context.Context, id string) error {
	var mentor model.Mentor
	err := r.db.First(&mentor, "id = ?", id).Error
	if err != nil {
		return errs.NotFoundError{
			Message: fmt.Sprintf("Mentor with ID %s not found", id),
		}
	}
	if mentor.Disabled {
		return errs.BadRequestError{
			Message: fmt.Sprintf("Mentor with ID %s already disabled", id),
		}
	}
	mentor.Disabled = true
	_, err = r.Update(ctx, mentor)
	if err != nil {
		return fmt.Errorf("MentorRepo - Delete - r.Update: %w", err)
	}
	return nil
}
