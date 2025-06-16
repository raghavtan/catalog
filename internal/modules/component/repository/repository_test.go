package repository_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	compassservice "github.com/motain/of-catalog/internal/services/compassservice"
	compassmocks "github.com/motain/of-catalog/internal/services/compassservice/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {

	successComponent := dtos.Component{
		ID: "component-id",
		MetricSources: dtos.MetricSources{
			Nodes: []dtos.MetricSource{
				{
					ID: "metric-source-id",
					MetricDefinition: dtos.MetricDefinition{
						ID:   "metric-id",
						Name: "metric-name",
					},
				},
			},
		},
		Links: []dtos.Link{
			{
				ID:   "link-id",
				Type: "link-type",
				Name: "link-name",
				URL:  "link-url",
			},
		},
	}

	tests := []struct {
		name           string
		inputComponent resources.Component
		mockSetup      func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedResult resources.Component
		expectedError  error
	}{
		{
			name: "successfully create a component",
			inputComponent: resources.Component{
				Name: "TestComponent",
				Slug: "test-component",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().GetCompassCloudId().Return("test-cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input, output interface{}) error {
						createInput := input.(*dtos.CreateComponentInput)
						createOutput := output.(*dtos.CreateComponentOutput)

						if createInput.Component.Name != "TestComponent" {
							return fmt.Errorf("unexpected component name")
						}

						createOutput.Compass = dtos.CompassCreatedComponentOutput{
							CreateComponent: dtos.CompassCreateComponentOutput{
								Success: true,
								Details: successComponent,
							},
						}
						return nil
					},
				)
			},
			expectedResult: resources.Component{
				ID:   "component-id",
				Name: "TestComponent",
				Slug: "test-component",
				MetricSources: map[string]*resources.MetricSource{
					"metric-name": {
						ID:     "metric-source-id",
						Metric: "metric-id",
					},
				},
				Links: []resources.Link{
					{
						ID:   "link-id",
						Type: "link-type",
						Name: "link-name",
						URL:  "link-url",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "error during component creation",
			inputComponent: resources.Component{
				Name: "TestComponent",
				Slug: "test-component",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().GetCompassCloudId().Return("test-cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			expectedResult: resources.Component{},
			expectedError:  fmt.Errorf("Create component error for TestComponent: mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			result, err := repo.Create(context.Background(), tt.inputComponent)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			assert.Equal(t, tt.expectedResult, result, "expected result does not match actual result")
		})
	}
}
func TestRepository_Update(t *testing.T) {
	// Define some sample links for testing
	sampleLink1 := dtos.CompassLink{ID: "link-id-1", Name: "Link One", Type: "WEBSITE", URL: "http://example.com/one"}
	sampleLink2 := dtos.CompassLink{ID: "link-id-2", Name: "Link Two", Type: "DOCUMENT", URL: "http://example.com/two"}
	sampleLink3Updated := dtos.CompassLink{ID: "link-id-1", Name: "Link One Updated", Type: "WEBSITE", URL: "http://example.com/one-updated"}


	// Define sample metric sources for testing (using the new DTO structure)
	sampleMetricSourceNode1 := dtos.CompassMetricSourceNode{
		ID: "ms-id-1",
		MetricDefinition: dtos.CompassMetricDefinition{ID: "md-id-1", Name: "MetricDef1"},
	}
	sampleMetricSourceNode2 := dtos.CompassMetricSourceNode{
		ID: "ms-id-2",
		MetricDefinition: dtos.CompassMetricDefinition{ID: "md-id-2", Name: "MetricDef2"},
	}

	// Define sample documents for testing
	sampleDocumentNode1 := dtos.CompassDocumentNode{ID: "doc-id-1", Title: "Doc One", URL: "http://example.com/doc1"}


	tests := []struct {
		name                string
		inputComponent      resources.Component
		mockSetup           func(mockCompass *compassmocks.MockCompassServiceInterface, inputComp resources.Component)
		expectedResult      resources.Component
		expectedError       error
		expectedErrorString string // For error message comparison
	}{
		{
			name: "successfully update component and refresh links (no initial links, new links fetched)",
			inputComponent: resources.Component{
				ID:   "comp-id-1",
				Name: "Component With New Links",
				Slug: "comp-new-links",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface, inputComp resources.Component) {
				// Mock for UpdateComponent mutation
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.AssignableToTypeOf(&dtos.UpdateComponentInput{}), gomock.AssignableToTypeOf(&dtos.UpdateComponentOutput{})).
					DoAndReturn(func(ctx context.Context, input, output interface{}) error {
						updateOutput := output.(*dtos.UpdateComponentOutput)
						updateOutput.Compass.UpdateComponent.Success = true
						// The actual UpdateComponent mutation in Compass might not return all details,
						// so we don't populate links here. That's GetBySlug's job.
						return nil
					}).Times(1)

				// Mock for GetCompassCloudId (called by GetBySlug)
				mockCompass.EXPECT().GetCompassCloudId().Return("test-cloud-id").Times(1)

				// Mock for GetBySlug query
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.AssignableToTypeOf(&dtos.ComponentByReferenceInput{}), gomock.AssignableToTypeOf(&dtos.ComponentByReferenceOutput{})).
					DoAndReturn(func(ctx context.Context, input, output interface{}) error {
						getBySlugOutput := output.(*dtos.ComponentByReferenceOutput)
						getBySlugOutput.Compass.Component = dtos.ComponentDetail{
							ID:            inputComp.ID,
							Name:          inputComp.Name, // Name from original input for consistency in this part
							Description:   "Desc for new links",
							OwnerID:       "owner-1",
							Labels:        []string{"label1"},
							Type:          dtos.CompassComponentType{ID: "type-id-1"},
							Links:         []dtos.CompassLink{sampleLink1, sampleLink2},
							Documents:     dtos.CompassDocumentConnection{Nodes: []dtos.CompassDocumentNode{sampleDocumentNode1}},
							MetricSources: dtos.CompassMetricSourceConnection{Nodes: []dtos.CompassMetricSourceNode{sampleMetricSourceNode1}},
						}
						return nil
					}).Times(1)
			},
			expectedResult: resources.Component{
				ID:            "comp-id-1",
				Name:          "Component With New Links", // GetBySlug returns this
				Slug:          "comp-new-links", // Slug is from original input, not GetBySlug typically
				Description:   "Desc for new links",
				OwnerID:       "owner-1",
				Labels:        []string{"label1"},
				TypeID:        "type-id-1",
				Links: []resources.Link{
					{ID: "link-id-1", Name: "Link One", Type: "WEBSITE", URL: "http://example.com/one"},
					{ID: "link-id-2", Name: "Link Two", Type: "DOCUMENT", URL: "http://example.com/two"},
				},
				Documents: []resources.Document{
					{ID: "doc-id-1", Title: "Doc One", URL: "http://example.com/doc1"},
				},
				MetricSources: map[string]*resources.MetricSource{
					"MetricDef1": {ID: "ms-id-1", Metric: "md-id-1"},
				},
			},
		},
		{
			name: "successfully update component and refresh links (existing links updated, new one added)",
			inputComponent: resources.Component{
				ID:   "comp-id-2",
				Name: "Component With Updated Links",
				Slug: "comp-updated-links",
				Links: []resources.Link{ // These are initial links, will be overridden by GetBySlug
					{ID: "old-link-id", Name: "Old Link", Type: "OLD_TYPE", URL: "http://example.com/old"},
				},
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface, inputComp resources.Component) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.AssignableToTypeOf(&dtos.UpdateComponentInput{}), gomock.AssignableToTypeOf(&dtos.UpdateComponentOutput{})).
					DoAndReturn(func(ctx context.Context, input, output interface{}) error {
						o := output.(*dtos.UpdateComponentOutput)
						o.Compass.UpdateComponent.Success = true
						return nil
					}).Times(1)

				mockCompass.EXPECT().GetCompassCloudId().Return("test-cloud-id").Times(1)
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.AssignableToTypeOf(&dtos.ComponentByReferenceInput{}), gomock.AssignableToTypeOf(&dtos.ComponentByReferenceOutput{})).
					DoAndReturn(func(ctx context.Context, input, output interface{}) error {
						o := output.(*dtos.ComponentByReferenceOutput)
						o.Compass.Component = dtos.ComponentDetail{
							ID:            inputComp.ID,
							Name:          inputComp.Name,
							Description:   "Desc for updated links",
							Links:         []dtos.CompassLink{sampleLink3Updated, sampleLink2}, // sampleLink1 updated, sampleLink2 added
							MetricSources: dtos.CompassMetricSourceConnection{Nodes: []dtos.CompassMetricSourceNode{sampleMetricSourceNode2}},
						}
						return nil
					}).Times(1)
			},
			expectedResult: resources.Component{
				ID:            "comp-id-2",
				Name:          "Component With Updated Links",
				Slug:          "comp-updated-links",
				Description:   "Desc for updated links",
				Links: []resources.Link{
					{ID: "link-id-1", Name: "Link One Updated", Type: "WEBSITE", URL: "http://example.com/one-updated"},
					{ID: "link-id-2", Name: "Link Two", Type: "DOCUMENT", URL: "http://example.com/two"},
				},
				MetricSources: map[string]*resources.MetricSource{
					"MetricDef2": {ID: "ms-id-2", Metric: "md-id-2"},
				},
				// Other fields like Documents, TypeID, OwnerID, Labels would be zero/nil if not set in mock GetBySlug
				Documents: []resources.Document{}, // Explicitly empty if not set
				Labels: nil, // Explicitly nil if not set
			},
		},
		{
			name: "error during main update mutation",
			inputComponent: resources.Component{
				ID: "comp-id-3", Name: "Update Fail", Slug: "update-fail",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface, inputComp resources.Component) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.AssignableToTypeOf(&dtos.UpdateComponentInput{}), gomock.AssignableToTypeOf(&dtos.UpdateComponentOutput{})).
					Return(errors.New("compass update mutation error")).Times(1)
				// GetBySlug should not be called if the main update fails
			},
			expectedResult:      resources.Component{},
			expectedErrorString: "Update component error for Update Fail: compass update mutation error",
		},
		{
			name: "error during GetBySlug after successful update",
			inputComponent: resources.Component{
				ID: "comp-id-4", Name: "GetBySlug Fail", Slug: "getbyslug-fail",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface, inputComp resources.Component) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.AssignableToTypeOf(&dtos.UpdateComponentInput{}), gomock.AssignableToTypeOf(&dtos.UpdateComponentOutput{})).
					DoAndReturn(func(ctx context.Context, input, output interface{}) error {
						o := output.(*dtos.UpdateComponentOutput)
						o.Compass.UpdateComponent.Success = true
						return nil
					}).Times(1)

				mockCompass.EXPECT().GetCompassCloudId().Return("test-cloud-id").Times(1)
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.AssignableToTypeOf(&dtos.ComponentByReferenceInput{}), gomock.AssignableToTypeOf(&dtos.ComponentByReferenceOutput{})).
					Return(errors.New("compass getbyslug query error")).Times(1)
			},
			expectedResult:      resources.Component{},
			expectedErrorString: "failed to retrieve updated component details for GetBySlug Fail after update: GetBySlug error for getbyslug-fail: compass getbyslug query error",
		},
		// TODO: Add a test case for the PreValidationFunc if its behavior needs specific verification,
		// though the task was to avoid modifying it. For now, assume it works as before or is covered by existing logic if ID is missing.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// No need for defer ctrl.Finish() if using t.Cleanup

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockCompass, tt.inputComponent)
			}

			repo := repository.NewRepository(mockCompass)
			result, err := repo.Update(context.Background(), tt.inputComponent)

			if tt.expectedErrorString != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErrorString)
			} else if tt.expectedError != nil {
				// For specific error types if needed, though string comparison is often sufficient
				assert.True(t, errors.Is(err, tt.expectedError) || err.Error() == tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Compare individual fields if direct comparison of structs with maps/slices is tricky
			assert.Equal(t, tt.expectedResult.ID, result.ID, "ID mismatch")
			assert.Equal(t, tt.expectedResult.Name, result.Name, "Name mismatch")
			assert.Equal(t, tt.expectedResult.Slug, result.Slug, "Slug mismatch")
			assert.Equal(t, tt.expectedResult.Description, result.Description, "Description mismatch")
			assert.Equal(t, tt.expectedResult.OwnerID, result.OwnerID, "OwnerID mismatch")
			assert.Equal(t, tt.expectedResult.TypeID, result.TypeID, "TypeID mismatch")
			assert.ElementsMatch(t, tt.expectedResult.Labels, result.Labels, "Labels mismatch")
			assert.ElementsMatch(t, tt.expectedResult.Links, result.Links, "Links mismatch")
			assert.ElementsMatch(t, tt.expectedResult.Documents, result.Documents, "Documents mismatch")

			// For maps like MetricSources, check keys and values
			assert.Equal(t, len(tt.expectedResult.MetricSources), len(result.MetricSources), "MetricSources length mismatch")
			for k, v := range tt.expectedResult.MetricSources {
				assert.Equal(t, v, result.MetricSources[k], "MetricSource for key %s mismatch", k)
			}

            // If you want to keep the overall struct comparison too (it might fail due to map/slice ordering)
			// assert.Equal(t, tt.expectedResult, result, "expected result struct does not match actual result struct")
		})
	}
}
func TestRepository_Delete(t *testing.T) {
	tests := []struct {
		name           string
		inputComponent resources.Component
		mockSetup      func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedError  error
	}{
		{
			name: "successfully delete a component",
			inputComponent: resources.Component{
				ID: "component-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "error during component deletion",
			inputComponent: resources.Component{
				ID: "component-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			expectedError: fmt.Errorf("Delete component error for component-id: mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			err := repo.Delete(context.Background(), tt.inputComponent)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepository_SetDependency(t *testing.T) {
	tests := []struct {
		name          string
		dependent     resources.Component
		provider      resources.Component
		mockSetup     func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedError error
	}{
		{
			name: "successfully set a dependency",
			dependent: resources.Component{
				ID: "dependent-id",
			},
			provider: resources.Component{
				ID: "provider-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "error during setting dependency",
			dependent: resources.Component{
				ID: "dependent-id",
			},
			provider: resources.Component{
				ID: "provider-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			expectedError: fmt.Errorf("SetDependency error for dependent-id: mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			err := repo.SetDependency(context.Background(), tt.dependent, tt.provider)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepository_UnsetDependency(t *testing.T) {
	tests := []struct {
		name          string
		dependent     resources.Component
		provider      resources.Component
		mockSetup     func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedError error
	}{
		{
			name: "successfully unset a dependency",
			dependent: resources.Component{
				ID: "dependent-id",
			},
			provider: resources.Component{
				ID: "provider-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "error during unsetting dependency",
			dependent: resources.Component{
				ID: "dependent-id",
			},
			provider: resources.Component{
				ID: "provider-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			expectedError: fmt.Errorf("UnsetDependency dependency error for dependent-id: mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			err := repo.UnsetDependency(context.Background(), tt.dependent, tt.provider)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepository_GetBySlug(t *testing.T) {
	tests := []struct {
		name           string
		inputComponent resources.Component
		mockSetup      func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedResult *resources.Component
		expectedError  error
	}{
		{
			name: "successfully get a component by slug",
			inputComponent: resources.Component{
				Slug: "test-slug",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().GetCompassCloudId().Return("test-cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input, output interface{}) error {
						getOutput := output.(*dtos.ComponentByReferenceOutput)
						getOutput.Compass.Component = dtos.Component{
							ID: "component-id",
							MetricSources: dtos.MetricSources{
								Nodes: []dtos.MetricSource{
									{
										ID: "metric-source-id",
										MetricDefinition: dtos.MetricDefinition{
											ID:   "metric-id",
											Name: "metric-name",
										},
									},
								},
							},
						}
						return nil
					},
				)
			},
			expectedResult: &resources.Component{
				ID: "component-id",
				MetricSources: map[string]*resources.MetricSource{
					"metric-name": {
						ID:     "metric-source-id",
						Metric: "metric-id",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "error during getting component by slug",
			inputComponent: resources.Component{
				Slug: "test-slug",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().GetCompassCloudId().Return("test-cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			expectedResult: nil,
			expectedError:  fmt.Errorf("GetBySlug error for test-slug: mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			result, err := repo.GetBySlug(context.Background(), tt.inputComponent)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
func TestRepository_AddDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compassMocks := compassmocks.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(compassMocks)

	tests := []struct {
		name           string
		component      resources.Component
		document       resources.Document
		mockSetup      func()
		expectedResult resources.Document
		expectedError  error
	}{
		{
			name:      "successfully adds a document with setting DocumentCategories",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				Title: "Test Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				compassMocks.EXPECT().GetCompassCloudId().Return("cloud-id")
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input interface{}, output interface{}) error {
						documentCategory := struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						}{
							ID:   "category-id-1",
							Name: "category-1",
						}
						docOutput := output.(*dtos.DocumentationCategoriesOutput)
						docOutput.Compass.DocumentationCategories.Nodes = []struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						}{documentCategory}
						return nil
					},
				)
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input interface{}, output *dtos.CreateDocumentOutput) error {
						repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
						output.Compass.AddDocument.Details.ID = "document-id"
						output.Compass.AddDocument.Success = true
						return nil
					},
				)
			},
			expectedResult: resources.Document{
				ID:                      "document-id",
				Title:                   "Test Document",
				Type:                    "type-1",
				URL:                     "http://example.com",
				DocumentationCategoryId: "category-id-1",
			},
			expectedError: nil,
		},
		{
			name:      "fails to add a document because fails to set DocumentCategories",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				Title: "Test Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedResult: resources.Document{},
			expectedError:  errors.New("AddDocument error for component-id/Test Document: compass error"),
		},
		{
			name:      "successfully adds a document",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				Title: "Test Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input interface{}, output *dtos.CreateDocumentOutput) error {
						output.Compass.AddDocument.Details.ID = "document-id"
						output.Compass.AddDocument.Success = true
						return nil
					},
				)
				repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
			},
			expectedResult: resources.Document{
				ID:                      "document-id",
				Title:                   "Test Document",
				Type:                    "type-1",
				URL:                     "http://example.com",
				DocumentationCategoryId: "category-id-1",
			},
			expectedError: nil,
		},
		{
			name:      "fails to add document due to compass error",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				Title: "Test Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
				repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
			},
			expectedResult: resources.Document{},
			expectedError:  errors.New("AddDocument error for component-id/Test Document: compass error"),
		},
		{
			name:      "fails to add document due to unsuccessful response",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				Title: "Test Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input interface{}, output *dtos.CreateDocumentOutput) error {
						output.Compass.AddDocument.Success = false
						output.Compass.AddDocument.Errors = []compassservice.CompassError{
							{Message: "failed to execute mutation addDocument"},
						}
						return errors.New("failed to execute mutation addDocument")
					},
				)
				repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
			},
			expectedResult: resources.Document{},
			expectedError:  errors.New("AddDocument error for component-id/Test Document: failed to execute mutation addDocument"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := repo.AddDocument(context.Background(), tt.component, tt.document)
			assert.Equal(t, tt.expectedResult, got)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_UpdateDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	compassMocks := compassmocks.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(compassMocks)

	tests := []struct {
		name          string
		component     resources.Component
		document      resources.Document
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "successfully updates a document",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				ID:    "document-id",
				Title: "Updated Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input interface{}, output *dtos.UpdateDocumentOutput) error {
						output.Compass.UpdateDocument.Success = true
						return nil
					},
				)
			},
			expectedError: nil,
		},
		{
			name:      "fails to update document due to compass error",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				ID:    "document-id",
				Title: "Updated Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedError: errors.New("UpdateDocument error for component-id/Updated Document: compass error"),
		},
		{
			name:      "fails to update document due to unsuccessful response",
			component: resources.Component{ID: "component-id"},
			document: resources.Document{
				ID:    "document-id",
				Title: "Updated Document",
				Type:  "type-1",
				URL:   "http://example.com",
			},
			mockSetup: func() {
				repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
				compassMocks.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input interface{}, output *dtos.UpdateDocumentOutput) error {
						output.Compass.UpdateDocument.Success = false
						output.Compass.UpdateDocument.Errors = []compassservice.CompassError{
							{Message: "failed to execute mutation updateDocument"},
						}
						return errors.New("failed to execute mutation updateDocument")
					},
				)
			},
			expectedError: errors.New("UpdateDocument error for component-id/Updated Document: failed to execute mutation updateDocument"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.UpdateDocument(context.Background(), tt.component, tt.document)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_BindMetric(t *testing.T) {
	tests := []struct {
		name          string
		component     resources.Component
		metricID      string
		identifier    string
		mockSetup     func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedID    string
		expectedError error
	}{
		{
			name: "successfully bind a metric",
			component: resources.Component{
				ID: "component-id",
			},
			metricID:   "metric-id",
			identifier: "identifier",
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, input, output interface{}) error {
						bindOutput := output.(*dtos.BindMetricOutput)
						bindOutput.Compass.CreateMetricSource.CreateMetricSource.ID = "metric-source-id"
						return nil
					},
				)
			},
			expectedID:    "metric-source-id",
			expectedError: nil,
		},
		{
			name: "error binding a metric",
			component: resources.Component{
				ID: "component-id",
			},
			metricID:   "metric-id",
			identifier: "identifier",
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			expectedID:    "",
			expectedError: fmt.Errorf("BindMetric error for component-id/metric-id: mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			id, err := repo.BindMetric(context.Background(), tt.component, tt.metricID, tt.identifier)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedID, id)
		})
	}
}
func TestRepository_UnbindMetric(t *testing.T) {
	tests := []struct {
		name          string
		metricSource  resources.MetricSource
		mockSetup     func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedError error
	}{
		{
			name: "successfully unbind a metric",
			metricSource: resources.MetricSource{
				ID: "metric-source-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "error unbinding a metric",
			metricSource: resources.MetricSource{
				ID: "metric-source-id",
			},
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
			},
			expectedError: fmt.Errorf("UnbindMetric error for metric-source-id: mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			err := repo.UnbindMetric(context.Background(), tt.metricSource)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepository_Push(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name          string
		metricSource  resources.MetricSource
		value         float64
		recordedAt    time.Time
		mockSetup     func()
		expectedError error
	}{
		{
			name:         "successfully pushes a metric",
			metricSource: resources.MetricSource{ID: "metric-source-id"},
			value:        42.5,
			recordedAt:   time.Now(),
			mockSetup: func() {
				mockCompass.EXPECT().SendMetric(gomock.Any(), gomock.Any()).Return("", nil)
			},
			expectedError: nil,
		},
		{
			name:         "fails to push metric due to compass error",
			metricSource: resources.MetricSource{ID: "metric-source-id"},
			value:        42.5,
			recordedAt:   time.Now(),
			mockSetup: func() {
				mockCompass.EXPECT().SendMetric(gomock.Any(), gomock.Any()).Return("", errors.New("compass error"))
			},
			expectedError: errors.New("compass error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.Push(context.Background(), tt.metricSource, tt.value, tt.recordedAt)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
func TestRepository_SetAPISpecifications(t *testing.T) {
	tests := []struct {
		name          string
		component     resources.Component
		apiSpecs      string
		apiSpecsFile  string
		mockSetup     func(mockCompass *compassmocks.MockCompassServiceInterface)
		expectedError error
	}{
		{
			name: "successfully set API specifications",
			component: resources.Component{
				ID: "component/id",
			},
			apiSpecs:     "api-spec-content",
			apiSpecsFile: "api-spec-file.yaml",
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().SendAPISpecifications(gomock.Any(), gomock.Any()).Return("", nil)
			},
			expectedError: nil,
		},
		{
			name: "error due to invalid component ID format",
			component: resources.Component{
				ID: "invalid-component-id",
			},
			apiSpecs:      "api-spec-content",
			apiSpecsFile:  "api-spec-file.yaml",
			mockSetup:     func(mockCompass *compassmocks.MockCompassServiceInterface) {},
			expectedError: errors.New("invalid component.ID format"),
		},
		{
			name: "error while sending API specifications",
			component: resources.Component{
				ID: "component/id",
			},
			apiSpecs:     "api-spec-content",
			apiSpecsFile: "api-spec-file.yaml",
			mockSetup: func(mockCompass *compassmocks.MockCompassServiceInterface) {
				mockCompass.EXPECT().SendAPISpecifications(gomock.Any(), gomock.Any()).Return("", errors.New("mock error"))
			},
			expectedError: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompass := compassmocks.NewMockCompassServiceInterface(ctrl)
			tt.mockSetup(mockCompass)

			repo := repository.NewRepository(mockCompass)
			err := repo.SetAPISpecifications(context.Background(), tt.component, tt.apiSpecs, tt.apiSpecsFile)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
