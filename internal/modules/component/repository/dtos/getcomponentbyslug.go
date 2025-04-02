package dtos

import compassdtos "github.com/motain/of-catalog/internal/services/compassservice/dtos"

/*************
 * INPUT DTO *
 *************/
type ComponentByReferenceInput struct {
	compassdtos.InputDTO
	CompassCloudID string
	Slug           string
}

func (dto *ComponentByReferenceInput) GetQuery() string {
	return `
		query getComponentBySlug($cloudId: ID!, $slug: String!) {
			compass {
				componentByReference(reference: {slug: {slug: $slug, cloudId: $cloudId}}) {
					... on CompassComponent {
						id
						metricSources {
							... on CompassComponentMetricSourcesConnection {
								nodes {
									id,
									metricDefinition {
										name
									}
								}
							}
						}
					}
				}
			}
		}`
}

func (dto *ComponentByReferenceInput) SetVariables() map[string]interface{} {
	return map[string]interface{}{
		"cloudId": dto.CompassCloudID,
		"slug":    dto.Slug,
	}
}

/**************
 * OUTPUT DTO *
 **************/

type ComponentByReferenceOutput struct {
	Compass struct {
		Component Component `json:"componentByReference"`
	} `json:"compass"`
}

func (dto *ComponentByReferenceOutput) IsSuccessful() bool {
	return dto.Compass.Component.ID != ""
}

func (dto *ComponentByReferenceOutput) GetErrors() []string {
	return nil
}
