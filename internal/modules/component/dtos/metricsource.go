package dtos

import fsdtos "github.com/motain/fact-collector/internal/services/factsystem/dtos"

type MetricSourceStatus string

type MetricSourceDTO struct {
	ID     string                `yaml:"id"`
	Name   string                `yaml:"name"`
	Metric string                `yaml:"metric"`
	Facts  fsdtos.FactOperations `yaml:"facts"`
}

func GetMetricSourceUniqueKey(m *MetricSourceDTO) string {
	return m.Name
}
