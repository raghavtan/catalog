package repository_test

import (
	context "context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	resources "github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateComponent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx       context.Context
		component resources.Component
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		want      string
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				component: resources.Component{
					Name:        "test-component",
					Slug:        "test-slug",
					Description: "test-description",
					TypeID:      "type-id",
					Links: []resources.Link{
						{
							Type: "type",
							Name: "name",
							URL:  "url",
						},
					},
					Labels:  []string{"label1", "label2"},
					OwnerID: "owner-id",
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							CreateComponent struct {
								Success          bool `json:"success"`
								ComponentDetails struct {
									ID string `json:"id"`
								} `json:"componentDetails"`
							} `json:"createComponent"`
						} `json:"compass"`
					})
					resp.Compass.CreateComponent.Success = true
					resp.Compass.CreateComponent.ComponentDetails.ID = "component-id"
				}).Return(nil).Times(1)
			},
			want:    "component-id",
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx: context.TODO(),
				component: resources.Component{
					Name:        "test-component",
					Slug:        "test-slug",
					Description: "test-description",
					TypeID:      "type-id",
					Links: []resources.Link{
						{
							Type: "type",
							Name: "name",
							URL:  "url",
						},
					},
					Labels:  []string{"label1", "label2"},
					OwnerID: "owner-id",
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to create component")).Times(1)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repository.Create(tt.args.ctx, tt.args.component)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	componentID := "component-id"

	type args struct {
		ctx       context.Context
		component resources.Component
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				component: resources.Component{
					ID:          &componentID,
					Name:        "test-component",
					Slug:        "test-slug",
					Description: "test-description",
					OwnerID:     "owner-id",
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							UpdateComponentDefinition struct {
								Success bool `json:"success"`
							} `json:"updateComponent"`
						} `json:"compass"`
					})
					resp.Compass.UpdateComponentDefinition.Success = true
				}).Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx: context.TODO(),
				component: resources.Component{
					ID:          &componentID,
					Name:        "test-component",
					Slug:        "test-slug",
					Description: "test-description",
					OwnerID:     "owner-id",
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to update component")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "update not successful",
			args: args{
				ctx: context.TODO(),
				component: resources.Component{
					ID:          &componentID,
					Name:        "test-component",
					Slug:        "test-slug",
					Description: "test-description",
					OwnerID:     "owner-id",
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							UpdateComponentDefinition struct {
								Success bool `json:"success"`
							} `json:"updateComponent"`
						} `json:"compass"`
					})
					resp.Compass.UpdateComponentDefinition.Success = false
				}).Return(nil).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repository.Update(tt.args.ctx, tt.args.component)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				id:  "component-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							DeleteComponent struct {
								Success bool `json:"success"`
							} `json:"deleteComponent"`
						} `json:"compass"`
					})
					resp.Compass.DeleteComponent.Success = true
				}).Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx: context.TODO(),
				id:  "component-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to delete component")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repository.Delete(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
