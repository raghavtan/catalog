package dtos

import fsdtos "github.com/motain/of-catalog/internal/services/factsystem/dtos"

type MetricSourceStatus string

type MetricSourceDTO struct {
	ID     string         `yaml:"id"`
	Name   string         `yaml:"name"`
	Metric string         `yaml:"metric"`
	Facts  []*fsdtos.Task `yaml:"facts"`
}

func GetMetricSourceUniqueKey(m *MetricSourceDTO) string {
	return m.Name
}
