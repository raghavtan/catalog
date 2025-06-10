package ownerservice

import (
	"fmt"
	oforg "github.com/motain/of-org"
	"strings"

	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
)

type SquadDetails struct {
	JiraTeamID      string
	SlackURL        string
	SlackTitle      string
	JiraProjectURL  string
	JiraProjectName string
	Tribe           string
}

type OwnerServiceInterface interface {
	GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error)
}

type OwnerService struct {
	getSquadDetails func(string) (SquadDetails, error)
}

// NewOwnerService creates a new OwnerService with default dependencies
func NewOwnerService() *OwnerService {
	return &OwnerService{
		getSquadDetails: GetSquadDetails,
	}
}

// NewOwnerServiceWithDependencies allows injecting dependencies for testing
func NewOwnerServiceWithDependencies(getSquadDetails func(string) (SquadDetails, error)) *OwnerService {
	return &OwnerService{
		getSquadDetails: getSquadDetails,
	}
}

func (os *OwnerService) GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error) {
	squadDetails, err := os.getSquadDetails(squad)
	if err != nil {
		return nil, err
	}

	// Verify tribe matches
	if squadDetails.Tribe != tribe {
		return nil, fmt.Errorf("squad '%s' belongs to tribe '%s', not '%s'", squad, squadDetails.Tribe, tribe)
	}

	owner := &dtos.Owner{
		OwnerID:       squadDetails.JiraTeamID,
		SlackChannels: make(map[string]string),
		Projects:      make(map[string]string),
		DisplayName:   squad, // Using squad name as display name
	}

	// Add Slack channel if available
	if squadDetails.SlackURL != "" {
		owner.SlackChannels[squadDetails.SlackTitle] = squadDetails.SlackURL
	}

	// Add Jira project if available
	if squadDetails.JiraProjectURL != "" {
		owner.Projects[squadDetails.JiraProjectName] = squadDetails.JiraProjectURL
	}

	return owner, nil
}

// GetSquadDetails retrieves detailed information for a given squad
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
