package repository

import (
	"github.com/motain/fact-collector/internal/services/configservice"
)

type RepositoryInterface interface {
	FetchData() string
}

type Repository struct {
	config configservice.ConfigServiceInterface
}

func NewRepository(
	cfg configservice.ConfigServiceInterface,
) *Repository {
	return &Repository{config: cfg}
}

func (r *Repository) FetchData() string {
	return "Data"
}
