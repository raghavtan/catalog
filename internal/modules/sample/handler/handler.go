package handler

import (
	"fmt"

	"github.com/motain/fact-collector/internal/services/githubservice"
)

type Handler struct {
	github githubservice.GitHubRepositoriesServiceInterface
}

func NewHandler(gh githubservice.GitHubRepositoriesServiceInterface) *Handler {
	return &Handler{github: gh}
}

func (h *Handler) Handle() string {
	t, fileErr := h.github.GetFileContent("idp-wrapper", "app.toml")
	if fileErr != nil {
		panic(fileErr)
	}
	fmt.Printf("File Content: %s\n", t)

	return ""
}
