package user

import (
	"context"
	"fmt"

	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

// Service implements user-facing use cases.
type Service struct {
	users domainuser.Repository
}

// New creates a user service.
func New(users domainuser.Repository) *Service {
	return &Service{users: users}
}

// GetByID returns a user without exposing persistence details.
func (s *Service) GetByID(ctx context.Context, id string) (domainuser.User, error) {
	result, err := s.users.FindByID(ctx, id)
	if err != nil {
		return domainuser.User{}, fmt.Errorf("get user: %w", err)
	}
	return result, nil
}
