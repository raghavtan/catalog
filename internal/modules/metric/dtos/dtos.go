package dtos

import (
	"reflect"

	fsdtos "github.com/motain/fact-collector/internal/services/factsystem/dtos"
)

// MetricDTO is a data transfer object representing a metric definition.
type MetricDTO struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name          string                `yaml:"name"`
		Labels        map[string]string     `yaml:"labels"`
		ComponentType []string              `yaml:"componentType"`
		Facts         fsdtos.FactOperations `yaml:"facts"`
	} `yaml:"metadata"`
	Spec MetricSpec `yaml:"spec"`
}

func GetMetricUniqueKey(m *MetricDTO) string {
	return m.Spec.Name
}

func FromStateToConfig(state *MetricDTO, conf *MetricDTO) {
	conf.Spec.ID = state.Spec.ID
}

func IsEqualMetric(m1, m2 *MetricDTO) bool {
	return m1.Spec.Name == m2.Spec.Name &&
		m1.Spec.Description == m2.Spec.Description &&
		reflect.DeepEqual(m1.Spec.Format, m2.Spec.Format) &&
		m1.Metadata.Name == m2.Metadata.Name &&
		isEqualLabels(m1.Metadata.Labels, m2.Metadata.Labels) &&
		isEqualComponentTypes(m1.Metadata.ComponentType, m2.Metadata.ComponentType) &&
		m1.Metadata.Facts.IsEqual(m2.Metadata.Facts)
}

func isEqualLabels(l1, l2 map[string]string) bool {
	if len(l1) != len(l2) {
		return false
	}

	for k, v := range l1 {
		if l2[k] != v {
			return false
		}
	}

	return true
}

func isEqualComponentTypes(c1, c2 []string) bool {
	if len(c1) != len(c2) {
		return false
	}

	for i, c := range c1 {
		if c2[i] != c {
			return false
		}
	}

	return true
}

type MetricSpec struct {
	ID          string           `yaml:"id"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Format      MetricSpecFormat `yaml:"format"`
}

type MetricSpecFormat struct {
	Unit string `yaml:"unit"`
}

type MetricSourceStatus string

const (
	ActiveMetricSourceStatus   MetricSourceStatus = "active"
	InactiveMetricSourceStatus MetricSourceStatus = "inactvie"
)

// MetricSourceDTO is a data transfer object representing a metric source definition.
type MetricSourceDTO struct {
	APIVersion string                  `yaml:"apiVersion"`
	Kind       string                  `yaml:"kind"`
	Metadata   MetricSourceMetadataDTO `yaml:"metadata"`
	Spec       MetricSourceSpecDTO     `yaml:"spec"`
}

func GetMetricSourceUniqueKey(m *MetricSourceDTO) string {
	return m.Metadata.Name
}

type MetricSourceMetadataDTO struct {
	Name           string                `yaml:"name"`
	ComponentTypes []string              `yaml:"componentTypes"`
	Status         string                `yaml:"status"`
	Facts          fsdtos.FactOperations `yaml:"facts"`
}

func IsActiveMetricSources(metricSource *MetricSourceDTO) bool {
	return MetricSourceStatus(metricSource.Metadata.Status) == ActiveMetricSourceStatus
}

func IsInactiveMetricSources(metricSource *MetricSourceDTO) bool {
	return MetricSourceStatus(metricSource.Metadata.Status) == InactiveMetricSourceStatus
}

type MetricSourceSpecDTO struct {
	ID        *string `yaml:"id"`
	Name      string  `yaml:"name"`
	Metric    string  `yaml:"metric"`
	Component string  `yaml:"component"`
}
