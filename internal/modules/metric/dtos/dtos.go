package dtos

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
	Spec struct {
		ID          *string `yaml:"id"`
		Name        string  `yaml:"name"`
		Description string  `yaml:"description"`
		Format      struct {
			Unit string `yaml:"unit"`
		} `yaml:"format"`
	} `yaml:"spec"`
}

type FactOperations struct {
	All    []Fact `yaml:"all"`
	Any    []Fact `yaml:"any"`
	Report []Fact `yaml:"report"`
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

// MetricSourceDTO is a data transfer object representing a metric source definition.
type MetricSourceDTO struct {
	APIVersion string                  `yaml:"apiVersion"`
	Kind       string                  `yaml:"kind"`
	Metadata   MetricSourceMetadataDTO `yaml:"metadata"`
	Spec       MetricSourceSpecDTO     `yaml:"spec"`
}

type MetricSourceMetadataDTO struct {
	Name          string         `yaml:"name"`
	ComponentType []string       `yaml:"componentType"`
	Status        string         `yaml:"status"`
	Facts         FactOperations `yaml:"facts"`
}

type MetricSourceSpecDTO struct {
	ID        *string `yaml:"id"`
	Name      string  `yaml:"name"`
	Metric    string  `yaml:"metric"`
	Component string  `yaml:"component"`
}
