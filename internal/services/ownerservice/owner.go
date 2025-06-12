package ownerservice

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
	"gopkg.in/yaml.v3"
)

type SquadDetails struct {
	JiraTeamID      string
	SlackURL        string
	SlackTitle      string
	JiraProjectURL  string
	JiraProjectName string
	Tribe           string
}

// Group follows the backstage API definition (copied from oforg to avoid dependency)
type Group struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string            `yaml:"name"`
		Description string            `yaml:"description"`
		Annotations map[string]string `yaml:"annotations"`
		Links       []struct {
			Title string `yaml:"title"`
			URL   string `yaml:"url"`
			Icon  string `yaml:"icon"`
			Type  string `yaml:"type"`
		} `yaml:"links"`
	} `yaml:"metadata"`
	Spec struct {
		Profile struct {
			DisplayName string `yaml:"displayName"`
		} `yaml:"profile"`
		Type     string   `yaml:"type"`
		Parent   string   `yaml:"parent"`
		Children []string `yaml:"children"`
	} `yaml:"spec"`
}

type OwnerServiceInterface interface {
	GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error)
}

type OwnerService struct {
	squads map[string]Group
}

// NewOwnerService creates a new OwnerService with direct YAML parsing
func NewOwnerService() *OwnerService {
	squads, err := loadSquadsFromYAML()
	if err != nil {
		squads = make(map[string]Group)
	}

	return &OwnerService{
		squads: squads,
	}
}

// NewOwnerServiceWithSquads allows injecting squads for testing
func NewOwnerServiceWithSquads(squads map[string]Group) *OwnerService {
	return &OwnerService{
		squads: squads,
	}
}

func (os *OwnerService) GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error) {
	squadDetails, err := os.getSquadDetails(squad)
	if err != nil {
		return nil, err
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

// loadSquadsFromYAML reads and parses the YAML file directly
// loadSquadsFromYAML reads and parses the YAML file directly
func loadSquadsFromYAML() (map[string]Group, error) {
	// Try multiple possible paths for the YAML file
	possiblePaths := []string{
		"vendor/github.com/motain/of-org/main.yaml",
		"../vendor/github.com/motain/of-org/main.yaml",
		"../../vendor/github.com/motain/of-org/main.yaml",
		"go.mod", // Will find go.mod and construct path from there
	}

	var yamlPath string
	var found bool

	// Find the YAML file
	for _, path := range possiblePaths {
		if path == "go.mod" {
			// Special case: find go.mod and construct path
			if goModPath, err := findGoMod(); err == nil {
				yamlPath = filepath.Join(filepath.Dir(goModPath), "vendor/github.com/motain/of-org/main.yaml")
				if _, err := os.Stat(yamlPath); err == nil {
					found = true
					break
				}
			}
		} else {
			if _, err := os.Stat(path); err == nil {
				yamlPath = path
				found = true
				break
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("could not find main.yaml in any expected location")
	}

	fmt.Printf("Loading squads from: %s\n", yamlPath)

	// Read the YAML file
	content, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Parse YAML documents (using same logic as oforg)
	dec := yaml.NewDecoder(bytes.NewReader(content))

	squads := make(map[string]Group)

	// FIXED: Create new variable inside the loop to avoid reference issues
	for {
		group := Group{} // ← FIXED: Declare variable inside loop
		if err := dec.Decode(&group); err != nil {
			break // End of documents or error
		}

		// Only store squads (not tribes or areas)
		if group.Spec.Type == "squad" {
			squads[group.Metadata.Name] = group
		}
	}

	return squads, nil
}

// findGoMod walks up the directory tree to find go.mod
func findGoMod() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return goModPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found")
}

// getSquadDetails retrieves detailed information for a given squad
func (os *OwnerService) getSquadDetails(squadName string) (SquadDetails, error) {
	squadName = strings.TrimSpace(squadName)

	squad, exists := os.squads[squadName]
	if !exists {
		for name := range os.squads {
			fmt.Printf("'%s' ", name)
		}
		fmt.Printf("\n")
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
