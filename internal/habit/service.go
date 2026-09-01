package habit

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("habit not found")
var ErrValidation = errors.New("name is required")

// placeholderUserID is a temporary hardcoded value used until authentication
// (Phase 5) exists and we can derive the real user from the request.
// TODO(auth): replace this with the authenticated user's ID.
const placeholderUserID = 1

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]Habit, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) Create(ctx context.Context, input Habit) (Habit, error) {
	if input.Name == "" {
		return Habit{}, ErrValidation
	}
	input.UserID = placeholderUserID
	return s.repo.Create(ctx, input)
}

func (s *Service) Get(ctx context.Context, id int) (Habit, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int, input Habit) (Habit, error) {
	if input.Name == "" {
		return Habit{}, ErrValidation
	}
	return s.repo.Update(ctx, id, input)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
