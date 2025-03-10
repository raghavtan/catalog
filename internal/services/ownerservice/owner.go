package ownerservice

import (
	"bytes"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
)

var groups dtos.GroupList

type OwnerServiceInterface interface {
	GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error)
}

type OwnerService struct {
	gitHubService githubservice.GitHubServiceInterface
}

func NewOwnerService(gitHubService githubservice.GitHubServiceInterface) *OwnerService {
	return &OwnerService{
		gitHubService: gitHubService,
	}
}

func (os *OwnerService) GetOwnerByTribeAndSquad(tribe, squad string) (*dtos.Owner, error) {
	groups, extractErr := os.extractData()
	if extractErr != nil {
		return nil, extractErr
	}

	for _, group := range groups {
		if os.matchesTribeAndSquad(group, tribe, squad) {
			return os.mapGroupToOwner(group), nil
		}
	}

	return nil, fmt.Errorf("no matching group found")
}

func (os *OwnerService) extractData() (dtos.GroupList, error) {
	// Cacbe the groups to avoid multiple requests.
	// The cache is valid for one execution of the command.
	if groups != nil {
		return groups, nil
	}

	fileContent, fileErr := os.gitHubService.GetFileContent("of-org", "main.yaml")
	if fileErr != nil {
		return nil, fileErr
	}

	var results []*dtos.Group
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(fileContent)))
	for {
		var result dtos.Group
		decodeErr := decoder.Decode(&result)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				break
			}
			return nil, decodeErr
		}
		results = append(results, &result)
	}

	groups = results
	return results, nil
}

func (os *OwnerService) matchesTribeAndSquad(group *dtos.Group, tribe, squad string) bool {
	if group.Spec.Type != "squad" {
		return false
	}
	if group.Metadata.Name != squad {
		return false
	}
	if group.Spec.Parent != tribe {
		return false
	}

	return true
}

func (os *OwnerService) mapGroupToOwner(group *dtos.Group) *dtos.Owner {
	slackChannel := ""
	for _, link := range group.Metadata.Links {
		if link.Title == "Slack" {
			slackChannel = link.URL
		}
	}

	return &dtos.Owner{
		CompassID:    group.Spec.ID,
		Email:        group.Spec.Email,
		SlackChannel: slackChannel,
		DisplayName:  group.Spec.Profile.DisplayName,
	}
}
