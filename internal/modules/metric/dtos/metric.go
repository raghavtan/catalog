package dtos

// MetricDTO is a data transfer object representing a metric definition.
type MetricDTO struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name"`
		Labels map[string]string `yaml:"labels"`
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
