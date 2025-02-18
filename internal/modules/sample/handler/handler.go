package handler

import (
	"context"
	"fmt"

	"github.com/motain/fact-collector/internal/services/compassservice"
	"github.com/motain/fact-collector/internal/services/githubservice"
)

type Handler struct {
	github  githubservice.GitHubRepositoriesServiceInterface
	compass compassservice.CompassServiceInterface
}

func NewHandler(
	gh githubservice.GitHubRepositoriesServiceInterface,
	compass compassservice.CompassServiceInterface,
) *Handler {
	return &Handler{github: gh, compass: compass}
}

func (h *Handler) Handle() string {
	t, fileErr := h.github.GetFileContent("idp-wrapper", "app.toml")
	if fileErr != nil {
		panic(fileErr)
	}
	fmt.Printf("File Content: %s\n", t)

	q := `query getCloudId {
		tenantContexts(hostNames: ["onefootball.atlassian.net"]) {
		  cloudId
		}
	}`

	var resp struct {
		TenantContexts []struct {
			CloudId string `json:"cloudId"`
		} `json:"tenantContexts"`
	}

	err := h.compass.Run(context.Background(), q, nil, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Cloud ID: %+v\n", resp)

	return ""
}
