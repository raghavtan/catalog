package repository

//go:generate mockgen -destination=./mock_repository.go -package=repository github.com/motain/fact-collector/internal/modules/component/repository RepositoryInterface

import (
	"context"
	"errors"
	"log"

	"github.com/motain/fact-collector/internal/modules/component/resources"
	"github.com/motain/fact-collector/internal/services/compassservice"
)

type RepositoryInterface interface {
	Create(ctx context.Context, component resources.Component) (string, error)
	Update(ctx context.Context, component resources.Component) error
	Delete(ctx context.Context, id string) error
	GetBySlug(slug string) (*resources.Component, error)
}

type Repository struct {
	compass compassservice.CompassServiceInterface
}

func NewRepository(
	compass compassservice.CompassServiceInterface,
) *Repository {
	return &Repository{compass: compass}
}

func (r *Repository) Create(ctx context.Context, component resources.Component) (string, error) {
	query := `
		mutation createComponent ($cloudId: ID!, $componentDetails: CreateCompassComponentInput!) {
			compass {
				createComponent(cloudId: $cloudId, input: $componentDetails) {
					success
					componentDetails {
						id
					}
					errors {
						message
					}
				}
			}
		}`

	links := make([]map[string]string, 0)
	for _, link := range component.Links {
		links = append(links, map[string]string{
			"type": link.Type,
			"name": link.Name,
			"url":  link.URL,
		})
	}

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"componentDetails": map[string]interface{}{
			"name":        component.Name,
			"slug":        component.Slug,
			"description": component.Description,
			"typeId":      component.TypeID,
			"links":       links,
			"labels":      component.Labels,
		},
	}

	if component.OwnerID != "" {
		variables["componentDetails"].(map[string]interface{})["ownerId"] = component.OwnerID
	}

	var response struct {
		Compass struct {
			CreateComponent struct {
				Success          bool                          `json:"success"`
				Errors           []compassservice.CompassError `json:"errors"`
				ComponentDetails struct {
					ID string `json:"id"`
				} `json:"componentDetails"`
			} `json:"createComponent"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create component: %v", err)
		return "", err
	}

	if compassservice.HasAlreadyExistsError(response.Compass.CreateComponent.Errors) {
		remoteComponent, err := r.GetBySlug(component.Slug)
		if err != nil {
			return "", err
		}

		component.ID = remoteComponent.ID
		updateError := r.Update(ctx, component)

		return *remoteComponent.ID, updateError
	} else {
		if !response.Compass.CreateComponent.Success {
			return "", errors.New("failed to create component")
		}
	}

	return response.Compass.CreateComponent.ComponentDetails.ID, nil
}

func (r *Repository) Update(ctx context.Context, component resources.Component) error {
	query := `
		mutation updateComponent ($componentDetails: UpdateCompassComponentInput!) {
			compass {
				updateComponent(input: $componentDetails) {
					success
					errors {
						message
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"componentDetails": map[string]interface{}{
			"id":          component.ID,
			"name":        component.Name,
			"slug":        component.Slug,
			"description": component.Description,
		},
	}

	if component.OwnerID != "" {
		variables["componentDetails"].(map[string]interface{})["ownerId"] = component.OwnerID
	}
	var response struct {
		Compass struct {
			UpdateComponentDefinition struct {
				Success bool `json:"success"`
			} `json:"updateComponent"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to update component: %v", err)
		return err
	}

	if !response.Compass.UpdateComponentDefinition.Success {
		return errors.New("failed to update component")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		mutation deleteComponent($id: ID!) {
			compass {
				deleteComponent(input: {id: $id}) {
					deletedComponentId
					errors {
						message
					}
					success
				}
			}
		}`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Compass struct {
			DeleteComponent struct {
				Success bool `json:"success"`
			} `json:"deleteComponent"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create component: %v", err)
		return err
	}

	if !response.Compass.DeleteComponent.Success {
		return errors.New("failed to delete component")
	}

	return nil
}

func (r *Repository) GetBySlug(slug string) (*resources.Component, error) {
	query := `
		query getComponentBySlug($cloudId: ID!, $slug: String!) {
			compass {
				componentByReference(reference: {slug: {slug: $slug, cloudId: $cloudId}}) {
					... on CompassComponent {
						id
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"slug":    slug,
	}

	var response struct {
		Compass struct {
			ComponentByReference struct {
				ID string `json:"id"`
			} `json:"componentByReference"`
		} `json:"compass"`
	}

	if err := r.compass.Run(context.Background(), query, variables, &response); err != nil {
		log.Printf("Failed to get component by slug: %v", err)
		return nil, err
	}

	component := resources.Component{
		ID: &response.Compass.ComponentByReference.ID,
	}

	return &component, nil
}
