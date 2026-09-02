package habit

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("habit not found")
var ErrValidation = errors.New("name is required")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID int) ([]Habit, error) {
	return s.repo.GetAll(ctx, userID)
}

func (s *Service) Create(ctx context.Context, userID int, input Habit) (Habit, error) {
	if input.Name == "" {
		return Habit{}, ErrValidation
	}
	input.UserID = userID
	return s.repo.Create(ctx, input)
}

func (s *Service) Get(ctx context.Context, userID int, id int) (Habit, error) {
	return s.repo.Get(ctx, id, userID)
}

func (s *Service) Update(ctx context.Context, userID int, id int, input Habit) (Habit, error) {
	if input.Name == "" {
		return Habit{}, ErrValidation
	}
	return s.repo.Update(ctx, id, userID, input)
}

func (s *Service) Delete(ctx context.Context, userID int, id int) error {
	return s.repo.Delete(ctx, id, userID)
}
