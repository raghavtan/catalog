package ownerservice

import (
	"fmt"
	"strings"

	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"

	oforg "github.com/motain/of-org"
)

type OwnerServiceInterface interface {
	GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error)
}

type OwnerService struct{}

func NewOwnerService() *OwnerService {
	return &OwnerService{}
}

func (os *OwnerService) GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error) {
	squadDetails, err := GetSquadDetails(squad)
	if err != nil {
		return nil, err
	}

	if squadDetails.Tribe != tribe {
		return nil, fmt.Errorf("squad '%s' belongs to tribe '%s', not '%s'", squad, squadDetails.Tribe, tribe)
	}

	owner := &dtos.Owner{
		OwnerID:       squadDetails.JiraTeamID,
		SlackChannels: make(map[string]string),
		Projects:      make(map[string]string),
		DisplayName:   squad,
	}

	if squadDetails.SlackURL != "" {
		owner.SlackChannels[squadDetails.SlackTitle] = squadDetails.SlackURL
	}

	if squadDetails.JiraProjectURL != "" {
		owner.Projects[squadDetails.JiraProjectName] = squadDetails.JiraProjectURL
	}

	return owner, nil
}

type SquadDetails struct {
	JiraTeamID      string
	SlackURL        string
	SlackTitle      string
	JiraProjectURL  string
	JiraProjectName string
	Tribe           string
}

func GetSquadDetails(squadName string) (SquadDetails, error) {
	squadName = strings.TrimSpace(squadName)

	squad, exists := oforg.Squads[squadName]
	if !exists {
		return SquadDetails{}, fmt.Errorf("squad '%s' not found", squadName)
	}

	details := SquadDetails{
		Tribe: squad.Spec.Parent,
	}

	if squad.Metadata.Annotations != nil {
		details.JiraTeamID = squad.Metadata.Annotations["jiraTeamID"]
	}

	for _, link := range squad.Metadata.Links {
		switch strings.ToLower(link.Type) {
		case "slack":
			details.SlackURL = link.URL
			details.SlackTitle = link.Title
		case "project":
			if strings.ToLower(link.Icon) == "jira" {
				details.JiraProjectURL = link.URL
				details.JiraProjectName = link.Title
			}
		}
	}

	return details, nil
}
