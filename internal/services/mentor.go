package services

import (
	"booking-service/internal/model"

	"booking-service/internal/repositories"
	"context"
	"fmt"
)

type mentorServiceImpl struct {
	repo repositories.MentorRepository
}

func NewMentorService(r repositories.MentorRepository) MentorService {
	return &mentorServiceImpl{
		repo: r,
	}
}

func (s *mentorServiceImpl) GetMentorByID(ctx context.Context, id string) (model.Mentor, error) {
	mentor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Mentor{}, err
	}

	return mentor, nil
}

func (s *mentorServiceImpl) GetMentors(ctx context.Context) ([]model.Mentor, error) {
	mentors, err := s.repo.GetMany(ctx)
	if err != nil {
		return nil, fmt.Errorf("MentorServiceImpl - GetMentors - s.repo.GetAll: %w", err)
	}

	return mentors, nil
}

func (s *mentorServiceImpl) CreatedMentor(ctx context.Context, mentor model.Mentor) (model.Mentor, error) {
	mentor, err := s.repo.Create(ctx, mentor)
	if err != nil {
		return model.Mentor{}, fmt.Errorf("MentorServiceImpl - CreatedMentor - s.repo.Create: %w", err)
	}

	return mentor, nil
}

func (s *mentorServiceImpl) UpdateMentor(ctx context.Context, mentor model.Mentor) (model.Mentor, error) {
	mentor, err := s.repo.Update(ctx, mentor)
	if err != nil {
		return model.Mentor{}, fmt.Errorf("MentorServiceImpl - UpdateMentor - s.repo.Update: %w", err)
	}

	return mentor, nil
}

func (s *mentorServiceImpl) DeleteMentor(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
