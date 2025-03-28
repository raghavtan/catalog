package prometheusservice

import (
	"time"

	"github.com/prometheus/common/model"

	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)
type PrometheusServiceInterface interface {
	Query(query string) (v1.Value, error)
}
type PrometheusService struct {
	Client v1.API
}

func NewPrometheusService(config configservice.ConfigServiceInterface) *PrometheusService {
	client, _ := api.NewClient(api.Config{Address: config.GetPrometheusURL()})
	return &PrometheusService{Client: v1.NewAPI(client)}
	return p.Client.Query(query, time.Now()) // Ensure the return type matches model.Value
func (p *PrometheusService) Query(query string) (v1.Value, error) {
	return p.Client.Query(query, time.Now())
}
