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
				Success          bool `json:"success"`
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

	if !response.Compass.CreateComponent.Success {
		return "", errors.New("failed to create component")
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
