package compassservice

import (
	"context"
	"log"

	"github.com/machinebox/graphql"
	"github.com/motain/fact-collector/internal/services/configservice"
)

type CompassServiceInterface interface {
	Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error
	GetCompassCloudId() string
}

type CompassService struct {
	client  GraphQLClientInterface
	token   string
	cloudId string
}

func NewCompassService(config configservice.ConfigServiceInterface, client GraphQLClientInterface) *CompassService {
	return &CompassService{
		client:  client,
		token:   config.GetCompassToken(),
		cloudId: config.GetCompassCloudId(),
	}
}

func (c *CompassService) Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	req := graphql.NewRequest(query)
	for key, value := range variables {
		req.Var(key, value)
	}

	req.Header.Set("Authorization", "Basic "+c.token)

	if err := c.client.Run(ctx, req, response); err != nil {
		log.Printf("Failed to execute query: %v", err)
		return err
	}

	return nil
}

func (c *CompassService) GetCompassCloudId() string {
	return c.cloudId
}
