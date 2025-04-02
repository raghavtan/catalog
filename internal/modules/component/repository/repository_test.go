package repository_test

import (
	context "context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	compassserviceError "github.com/motain/of-catalog/internal/services/compassservice"
	compassservice "github.com/motain/of-catalog/internal/services/compassservice/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

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
		component      resources.Component
		mockSetup      func()
		expectedResult resources.Component
		expectedError  error
	}{
		{
			name:      "successfully creates a component",
			component: resources.Component{Slug: "test-slug"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateComponentOutput) error {
						output.Compass = dtos.CompassCreatedComponentOutput{
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
				Slug: "test-slug",
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
			name:      "fails to create component due to compass error",
			component: resources.Component{Name: "test-slug"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateComponentOutput) error {
						return errors.New("compass error")
					},
				)
			},
			expectedResult: resources.Component{},
			expectedError:  errors.New("Create component error for test-slug: compass error"),
		},
		{
			name:      "fails to create component due to unsuccessful response",
			component: resources.Component{},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateComponentOutput) error {
						output.Compass = dtos.CompassCreatedComponentOutput{
							CreateComponent: dtos.CompassCreateComponentOutput{
								Success: false,
								Errors: []compassserviceError.CompassError{
									{Message: "error message"},
								},
							},
						}

						return nil
					},
				)
			},
			expectedResult: resources.Component{},
			expectedError:  errors.New("Create component error for : failed to execute mutation createComponent: [error message]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := repo.Create(context.Background(), tt.component)
			assert.Equal(t, tt.expectedResult, got)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name           string
		component      resources.Component
		mockSetup      func()
		expectedResult resources.Component
		expectedError  error
	}{
		{
			name: "successfully updates a component",
			component: resources.Component{
				ID:   "component-id",
				Slug: "test-slug",
			},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateComponentOutput) error {
						output.Compass.UpdateComponent.Success = true

						return nil
					},
				)
			},
			expectedResult: resources.Component{
				ID:   "component-id",
				Slug: "test-slug",
			},
			expectedError: nil,
		},
		{
			name: "successfully updates a component when component ID does not match remote one",
			component: resources.Component{
				ID:   "component-id",
				Slug: "test-slug",
			},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateComponentOutput) error {
						output.Compass.UpdateComponent.Success = false
						output.Compass.UpdateComponent.Errors = []compassserviceError.CompassError{
							{Message: "component not found"},
						}

						return nil
					},
				)
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.ComponentByReferenceOutput) error {
						output.Compass.Component = dtos.Component{
							ID: "component-id-1",
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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateComponentOutput) error {
						output.Compass.UpdateComponent.Success = true

						return nil
					},
				)
			},
			expectedResult: resources.Component{
				ID:   "component-id-1",
				Slug: "test-slug",
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
			name: "fails to update component due to compass error",
			component: resources.Component{
				ID:   "component-id",
				Slug: "test-slug",
			},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedResult: resources.Component{},
			expectedError:  errors.New("Update component error for : compass error"),
		},
		{
			name: "fails to update component due to unsuccessful response",
			component: resources.Component{
				ID:   "component-id",
				Slug: "test-slug",
			},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateComponentOutput) error {
						output.Compass.UpdateComponent.Success = false
						output.Compass.UpdateComponent.Errors = []compassserviceError.CompassError{
							{Message: "failed to run operation"},
						}

						return nil
					},
				)
			},
			expectedResult: resources.Component{},
			expectedError:  errors.New("Update component error for : failed to execute mutation updateComponent: [failed to run operation]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			component, err := repo.Update(context.Background(), tt.component)
			assert.Equal(t, tt.expectedResult, component)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name          string
		component     resources.Component
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "successfully deletes a component",
			component: resources.Component{ID: "component-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DeleteComponentOutput) error {
						output.Compass.DeleteComponent.Success = true
						return nil
					},
				)
			},
			expectedError: nil,
		},
		{
			name:      "fails to delete component due to compass error",
			component: resources.Component{ID: "component-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedError: errors.New("Delete component error for component-id: compass error"),
		},
		{
			name:      "fails to delete component due to unsuccessful response",
			component: resources.Component{ID: "component-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DeleteComponentOutput) error {
						output.Compass.DeleteComponent.Success = false
						output.Compass.DeleteComponent.Errors = []compassserviceError.CompassError{
							{Message: "failed to execute mutation deleteComponent"},
						}
						return nil
					},
				)
			},
			expectedError: errors.New("Delete component error for component-id: failed to execute mutation deleteComponent: [failed to execute mutation deleteComponent]"),
		},
		{
			name:      "component not found, no error returned",
			component: resources.Component{ID: "component-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DeleteComponentOutput) error {
						output.Compass.DeleteComponent.Errors = []compassserviceError.CompassError{
							{Message: "not found"},
						}
						return nil
					},
				)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.Delete(context.Background(), tt.component)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
func TestRepository_SetDependency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name          string
		dependent     resources.Component
		provider      resources.Component
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "successfully sets a dependency",
			dependent: resources.Component{ID: "dependent-id"},
			provider:  resources.Component{ID: "provider-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateDependencyOutput) error {
						output.Compass.CreateDependency.Success = true
						return nil
					},
				)
			},
			expectedError: nil,
		},
		{
			name:      "fails to set dependency due to compass error",
			dependent: resources.Component{ID: "dependent-id"},
			provider:  resources.Component{ID: "provider-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedError: errors.New("SetDependency error for dependent-id: compass error"),
		},
		{
			name:      "fails to set dependency due to unsuccessful response",
			dependent: resources.Component{ID: "dependent-id"},
			provider:  resources.Component{ID: "provider-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateDependencyOutput) error {
						output.Compass.CreateDependency.Success = false
						output.Compass.CreateDependency.Errors = []compassserviceError.CompassError{
							{Message: "failed to execute mutation createRelationship"},
						}
						return nil
					},
				)
			},
			expectedError: errors.New("SetDependency error for dependent-id: failed to execute mutation createRelationship: [failed to execute mutation createRelationship]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.SetDependency(context.Background(), tt.dependent, tt.provider)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_UnsetDependency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name          string
		dependent     resources.Component
		provider      resources.Component
		mockSetup     func()
		expectedError error
	}{
		{
			name:      "successfully unsets a dependency",
			dependent: resources.Component{ID: "dependent-id"},
			provider:  resources.Component{ID: "provider-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DeleteDependencyOutput) error {
						output.Compass.DeleteDependency.Success = true
						return nil
					},
				)
			},
			expectedError: nil,
		},
		{
			name:      "fails to unset dependency due to compass error",
			dependent: resources.Component{ID: "dependent-id"},
			provider:  resources.Component{ID: "provider-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedError: errors.New("UnsetDependency dependency error for dependent-id: compass error"),
		},
		{
			name:      "fails to unset dependency due to unsuccessful response",
			dependent: resources.Component{ID: "dependent-id"},
			provider:  resources.Component{ID: "provider-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DeleteDependencyOutput) error {
						output.Compass.DeleteDependency.Success = false
						output.Compass.DeleteDependency.Errors = []compassserviceError.CompassError{
							{Message: "failed to execute mutation deleteRelationship"},
						}
						return nil
					},
				)
			},
			expectedError: errors.New("UnsetDependency dependency error for dependent-id: failed to execute mutation deleteRelationship: [failed to execute mutation deleteRelationship]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.UnsetDependency(context.Background(), tt.dependent, tt.provider)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_GetBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name           string
		component      resources.Component
		mockSetup      func()
		expectedResult *resources.Component
		expectedError  error
	}{
		{
			name:      "successfully retrieves a component by slug",
			component: resources.Component{Slug: "test-slug"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.ComponentByReferenceOutput) error {
						output.Compass.Component = dtos.Component{
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
			name:      "fails to retrieve component due to compass error",
			component: resources.Component{Slug: "test-slug"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("GetBySlug error for test-slug: compass error"),
		},
		{
			name:      "fails to retrieve component due to unsuccessful response",
			component: resources.Component{Slug: "test-slug"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.ComponentByReferenceOutput) error {
						output.Compass.Component = dtos.Component{}
						return nil
					},
				)
			},
			expectedResult: nil,
			expectedError:  errors.New("GetBySlug error for test-slug: failed to execute query getComponentBySlug: []"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := repo.GetBySlug(context.Background(), tt.component)
			assert.Equal(t, tt.expectedResult, got)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_AddDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

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
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DocumentationCategoriesOutput) error {
						documentCategory := struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						}{
							ID:   "category-id-1",
							Name: "category-1",
						}
						output.Compass.DocumentationCategories.Nodes = []struct {
							ID   string `json:"id"`
							Name string `json:"name"`
						}{documentCategory}
						return nil
					},
				)
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateDocumentOutput) error {
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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateDocumentOutput) error {
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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateDocumentOutput) error {
						output.Compass.AddDocument.Success = false
						output.Compass.AddDocument.Errors = []compassserviceError.CompassError{
							{Message: "failed to execute mutation addDocument"},
						}
						return nil
					},
				)
				repo.DocumentCategories = map[string]string{"type-1": "category-id-1"}
			},
			expectedResult: resources.Document{},
			expectedError:  errors.New("AddDocument error for component-id/Test Document: failed to execute mutation addDocument: [failed to execute mutation addDocument]"),
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

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateDocumentOutput) error {
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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
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
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateDocumentOutput) error {
						output.Compass.UpdateDocument.Success = false
						output.Compass.UpdateDocument.Errors = []compassserviceError.CompassError{
							{Message: "failed to execute mutation updateDocument"},
						}
						return nil
					},
				)
			},
			expectedError: errors.New("UpdateDocument error for component-id/Updated Document: failed to execute mutation updateDocument: [failed to execute mutation updateDocument]"),
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name           string
		component      resources.Component
		metricID       string
		identifier     string
		mockSetup      func()
		expectedResult string
		expectedError  error
	}{
		{
			name:       "successfully binds a metric",
			component:  resources.Component{ID: "component-id"},
			metricID:   "metric-id",
			identifier: "identifier",
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.BindMetricOutput) error {
						output.Compass.CreateMetricSource.CreateMetricSource.ID = "metric-source-id"
						output.Compass.CreateMetricSource.Success = true
						return nil
					},
				)
			},
			expectedResult: "metric-source-id",
			expectedError:  nil,
		},
		{
			name:       "fails to bind metric due to compass error",
			component:  resources.Component{ID: "component-id"},
			metricID:   "metric-id",
			identifier: "identifier",
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedResult: "",
			expectedError:  errors.New("BindMetric error for component-id/metric-id: compass error"),
		},
		{
			name:       "fails to bind metric due to unsuccessful response",
			component:  resources.Component{ID: "component-id"},
			metricID:   "metric-id",
			identifier: "identifier",
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.BindMetricOutput) error {
						output.Compass.CreateMetricSource.Success = false
						output.Compass.CreateMetricSource.Errors = []compassserviceError.CompassError{
							{Message: "failed to execute mutation createMetricSource"},
						}
						return nil
					},
				)
			},
			expectedResult: "",
			expectedError:  errors.New("BindMetric error for component-id/metric-id: failed to execute mutation createMetricSource: [failed to execute mutation createMetricSource]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := repo.BindMetric(context.Background(), tt.component, tt.metricID, tt.identifier)
			assert.Equal(t, tt.expectedResult, got)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_UnbindMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name          string
		metricSource  resources.MetricSource
		mockSetup     func()
		expectedError error
	}{
		{
			name:         "successfully unbinds a metric",
			metricSource: resources.MetricSource{ID: "metric-source-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UnbindMetricOutput) error {
						output.Compass.DeleteMetricSource.Success = true
						return nil
					},
				)
			},
			expectedError: nil,
		},
		{
			name:         "fails to unbind metric due to compass error",
			metricSource: resources.MetricSource{ID: "metric-source-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			expectedError: errors.New("UnbindMetric error for metric-source-id: compass error"),
		},
		{
			name:         "fails to unbind metric due to unsuccessful response",
			metricSource: resources.MetricSource{ID: "metric-source-id"},
			mockSetup: func() {
				mockCompass.EXPECT().RunWithDTOs(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UnbindMetricOutput) error {
						output.Compass.DeleteMetricSource.Success = false
						output.Compass.DeleteMetricSource.Errors = []compassserviceError.CompassError{
							{Message: "failed to run operation"},
						}
						return nil
					},
				)
			},
			expectedError: errors.New("UnbindMetric error for metric-source-id: failed to execute mutation deleteMetricSource: [failed to run operation]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.UnbindMetric(context.Background(), tt.metricSource)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_Push(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
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
