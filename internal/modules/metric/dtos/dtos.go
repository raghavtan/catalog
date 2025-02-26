package dtos

import "reflect"

// MetricDTO is a data transfer object representing a metric definition.
type MetricDTO struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name          string            `yaml:"name"`
		Labels        map[string]string `yaml:"labels"`
		ComponentType []string          `yaml:"componentType"`
		Facts         FactOperations    `yaml:"facts"`
	} `yaml:"metadata"`
	Spec MetricSpec `yaml:"spec"`
}

func GetMetricUniqueKey(m *MetricDTO) string {
	return m.Spec.Name
}

func SetMetricID(m *MetricDTO, id string) {
	m.Spec.ID = id
}

func GetMetricID(m *MetricDTO) string {
	return m.Spec.ID
}

func IsEqualMetric(m1, m2 *MetricDTO) bool {
	return m1.Spec.Name == m2.Spec.Name &&
		m1.Spec.Description == m2.Spec.Description &&
		reflect.DeepEqual(m1.Spec.Format, m2.Spec.Format) &&
		m1.Metadata.Name == m2.Metadata.Name &&
		reflect.DeepEqual(m1.Metadata.Labels, m2.Metadata.Labels) &&
		reflect.DeepEqual(m1.Metadata.ComponentType, m2.Metadata.ComponentType) &&
		reflect.DeepEqual(m1.Metadata.Facts, m2.Metadata.Facts)

}

type FactOperations struct {
	All     []*Fact `yaml:"all"`
	Any     []*Fact `yaml:"any"`
	Inspect *Fact   `yaml:"inspect"`
}

// FactType defines the type of fact to collect
type FactType string

const (
	FileExistsFact     FactType = "fileExists"
	FileRegexFact      FactType = "fileRegex"
	FileJSONPathFact   FactType = "fileJsonPath"
	RepoPropertiesFact FactType = "repoProperties"
)

// Fact defines a struct to handle different facts
type Fact struct {
	Source          string   `yaml:"source"`          // Source of the fact (e.g., "github")
	URI             string   `yaml:"uri"`             // URI of the fact
	Name            string   `yaml:"name"`            // Name of the fact
	Repo            string   `yaml:"repo"`            // Repository name (e.g., "repo")
	FactType        FactType `yaml:"factType"`        // Type of fact to collect
	FilePath        string   `yaml:"filePath"`        // File to open/validate fact (if FactType is "fileExists", "fileRegex", or "fileJsonPath")
	RegexPattern    string   `yaml:"regexPattern"`    // Regex to match file content or response (if FactType is "fileRegex")
	JSONPath        string   `yaml:"jsonPath"`        // JSONPath to navigate file content or json response (if FactType is "fileJsonPath")
	RepoProperty    string   `yaml:"repoProperty"`    // Property to explore in the repo (if FactType is "repoProperties")
	ExpectedValue   string   `yaml:"expectedValue"`   // Expected value (matched against value of  "repoProperty" or "fileJsonPath")
	ExpectedFormula string   `yaml:"expectedFormula"` // Expected formula (matched against value of  "repoProperty" or "fileJsonPath")
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
	return m.Spec.Name
}

type MetricSourceMetadataDTO struct {
	Name          string         `yaml:"name"`
	ComponentType []string       `yaml:"componentType"`
	Status        string         `yaml:"status"`
	Facts         FactOperations `yaml:"facts"`
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
