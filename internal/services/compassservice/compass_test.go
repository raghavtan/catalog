package compassservice_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/motain/of-catalog/internal/services/compassservice/dtos"
	mocks "github.com/motain/of-catalog/internal/services/compassservice/mocks"
	configservicemocks "github.com/motain/of-catalog/internal/services/configservice/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCompassService_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGqlClient := mocks.NewMockGraphQLClientInterface(ctrl)
	mockHttpClient := mocks.NewMockHTTPClientInterface(ctrl)
	mockConfigService := configservicemocks.NewMockConfigServiceInterface(ctrl)

	type args struct {
		ctx       context.Context
		query     string
		variables map[string]interface{}
		response  interface{}
	}
	tests := []struct {
		name          string
		args          args
		mockSetup     func()
		expectedError error
	}{
		{
			name: "successfully executes query",
			args: args{
				ctx:   context.Background(),
				query: "query { test }",
				variables: map[string]interface{}{
					"key": "value",
				},
				response: &struct {
					Data string `json:"data"`
				}{},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockGqlClient.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "fails to execute query",
			args: args{
				ctx:   context.Background(),
				query: "query { test }",
				variables: map[string]interface{}{
					"key": "value",
				},
				response: &struct {
					Data string `json:"data"`
				}{},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockGqlClient.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("query execution failed"))
			},
			expectedError: errors.New("query execution failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			service := compassservice.NewCompassService(mockConfigService, mockGqlClient, mockHttpClient)

			err := service.Run(tt.args.ctx, tt.args.query, tt.args.variables, tt.args.response)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestCompassService_SendMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGqlClient := mocks.NewMockGraphQLClientInterface(ctrl)
	mockHttpClient := mocks.NewMockHTTPClientInterface(ctrl)
	mockConfigService := configservicemocks.NewMockConfigServiceInterface(ctrl)

	type args struct {
		body map[string]string
	}
	tests := []struct {
		name          string
		args          args
		mockSetup     func()
		expectedResp  string
		expectedError error
	}{
		{
			name: "successfully sends metric",
			args: args{
				body: map[string]string{
					"metric": "value",
				},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "/gateway/api/compass/v1/metrics", req.URL.Path)
					assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
					assert.Equal(t, "application/json", req.Header.Get("Accept"))

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`success`)),
					}, nil
				})
			},
			expectedResp:  "success",
			expectedError: nil,
		},
		{
			name: "fails to send metric due to HTTP error",
			args: args{
				body: map[string]string{
					"metric": "value",
				},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("http error"))
			},
			expectedResp:  "",
			expectedError: errors.New("failed to send request: http error"),
		},
		{
			name: "fails to send metric due to non-200 response",
			args: args{
				body: map[string]string{
					"metric": "value",
				},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockHttpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader(`"error"`)),
				}, nil)
			},
			expectedResp:  "",
			expectedError: errors.New("response body: \"error\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			service := compassservice.NewCompassService(mockConfigService, mockGqlClient, mockHttpClient)

			resp, err := service.SendMetric(context.TODO(), tt.args.body)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}

func TestCompassService_SendAPISpecifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGqlClient := mocks.NewMockGraphQLClientInterface(ctrl)
	mockHttpClient := mocks.NewMockHTTPClientInterface(ctrl)
	mockConfigService := configservicemocks.NewMockConfigServiceInterface(ctrl)

	type args struct {
		input dtos.APISpecificationsInput
	}
	tests := []struct {
		name          string
		args          args
		mockSetup     func()
		expectedResp  string
		expectedError error
	}{
		{
			name: "successfully sends API specifications",
			args: args{
				input: dtos.APISpecificationsInput{
					ComponentID: "component123",
					FileName:    "api-spec.yaml",
					ApiSpecs:    "openapi: 3.0.0",
				},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockHttpClient.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, "/gateway/api/compass/v1/component/component123/api_specs", req.URL.Path)
					assert.Equal(t, "application/json", req.Header.Get("Accept"))
					assert.Contains(t, req.Header.Get("Content-Type"), "multipart/form-data")

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`success`)),
					}, nil
				})
			},
			expectedResp:  "success",
			expectedError: nil,
		},
		{
			name: "fails to send API specifications due to HTTP error",
			args: args{
				input: dtos.APISpecificationsInput{
					ComponentID: "component123",
					FileName:    "api-spec.yaml",
					ApiSpecs:    "openapi: 3.0.0",
				},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockHttpClient.EXPECT().Do(gomock.Any()).Return(nil, errors.New("http error"))
			},
			expectedResp:  "",
			expectedError: errors.New("failed to send request: http error"),
		},
		{
			name: "fails to send API specifications due to non-200 response",
			args: args{
				input: dtos.APISpecificationsInput{
					ComponentID: "component123",
					FileName:    "api-spec.yaml",
					ApiSpecs:    "openapi: 3.0.0",
				},
			},
			mockSetup: func() {
				mockConfigService.EXPECT().GetCompassToken().Return("mock-token")
				mockConfigService.EXPECT().GetCompassCloudId().Return("mock-cloud-id")

				mockHttpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader(`"error"`)),
				}, nil)
			},
			expectedResp:  "",
			expectedError: errors.New("response body: \"error\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			service := compassservice.NewCompassService(mockConfigService, mockGqlClient, mockHttpClient)

			resp, err := service.SendAPISpecifications(context.TODO(), tt.args.input)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
