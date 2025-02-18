package compassservice

import (
	"context"

	"github.com/machinebox/graphql"
	"github.com/motain/fact-collector/internal/services/configservice"
)

type GraphQLClientInterface interface {
	Run(ctx context.Context, req *graphql.Request, resp interface{}) error
}

func NewGraphQLClient(config configservice.ConfigServiceInterface) GraphQLClientInterface {
	client := graphql.NewClient(config.GetCompassHost())
	// Keep this until we properly implement logging
	// client.Log = func(s string) { log.Println(s) }
	return client
}
