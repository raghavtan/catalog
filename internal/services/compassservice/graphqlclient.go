package compassservice

//go:generate mockgen -destination=./mocks/mock_graphql_client.go -package=compassservice github.com/motain/of-catalog/internal/services/compassservice GraphQLClientInterface

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
	"github.com/motain/of-catalog/internal/services/configservice"
)

type GraphQLClientInterface interface {
	Run(ctx context.Context, req *graphql.Request, resp interface{}) error
}

func NewGraphQLClient(config configservice.ConfigServiceInterface) GraphQLClientInterface {
	gqlUri := fmt.Sprintf("https://%s%s", config.GetCompassHost(), "/gateway/api/graphql")
	client := graphql.NewClient(gqlUri)

	// Keep this until we properly implement logging
	client.Log = func(s string) { log.Println(s) }
	return client
}
