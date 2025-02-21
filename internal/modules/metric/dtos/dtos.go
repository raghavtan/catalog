package dtos

// MetricDTO is a data transfer object representing a metric definition.
type MetricDTO struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name          string            `yaml:"name"`
		Labels        map[string]string `yaml:"labels"`
		ComponentType []string          `yaml:"componentType"`
	} `yaml:"metadata"`
	Spec struct {
		ID          *string `yaml:"id"`
		Name        string  `yaml:"name"`
		Description string  `yaml:"description"`
		Format      struct {
			Unit string `yaml:"unit"`
		} `yaml:"format"`
	} `yaml:"spec"`
}

// MetricSourceDTO is a data transfer object representing a metric source definition.
type MetricSourceDTO struct {
	APIVersion string                  `yaml:"apiVersion"`
	Kind       string                  `yaml:"kind"`
	Metadata   MetricSourceMetadataDTO `yaml:"metadata"`
	Spec       MetricSourceSpecDTO     `yaml:"spec"`
}

type MetricSourceMetadataDTO struct {
	Name          string   `yaml:"name"`
	ComponentType []string `yaml:"componentType"`
	Status        string   `yaml:"status"`
}

type MetricSourceSpecDTO struct {
	ID        *string `yaml:"id"`
	Name      string  `yaml:"name"`
	Metric    string  `yaml:"metric"`
	Component string  `yaml:"component"`
}
