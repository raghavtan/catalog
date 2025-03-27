package dtos

import (
	"reflect"

	fsdtos "github.com/motain/of-catalog/internal/services/factsystem/dtos"
)

// MetricDTO is a data transfer object representing a metric definition.
type MetricDTO struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name          string            `yaml:"name"`
		Labels        map[string]string `yaml:"labels"`
		ComponentType []string          `yaml:"componentType"`
		Facts         []*fsdtos.Task    `yaml:"facts"`
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
		isEqualFacts(m1.Metadata.Facts, m2.Metadata.Facts)
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

func isEqualFacts(f1, f2 []*fsdtos.Task) bool {
	if len(f1) != len(f2) {
		return false
	}

	for i, f := range f1 {
		if !f.IsEqual(f2[i]) {
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
