package handler

import "github.com/motain/fact-collector/internal/modules/sample/repository"

type Handler struct {
	repo repository.RepositoryInterface
}

func NewHandler(repo repository.RepositoryInterface) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle() string {
	return h.repo.FetchData()
}
