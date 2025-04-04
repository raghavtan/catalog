package dtos

type TaskType string

const (
	AggregateType TaskType = "aggregate"
	ExtractType   TaskType = "extract"
	ValidateType  TaskType = "validate"
)

type TaskRule string

const (
	// Extraction rules
	JSONPathRule TaskRule = "jsonpath"
	NotEmptyRule TaskRule = "notempty"
	SearchRule   TaskRule = "search"

	// Validation rules
	DepsMatchRule  TaskRule = "deps_match"
	UniqueRule     TaskRule = "unique"
	RegexMatchRule TaskRule = "regex_match"
	FormulaRule    TaskRule = "formula"
)

type TaskSource string

const (
	GitHubTaskSource     TaskSource = "github"
	JSONAPITaskSource    TaskSource = "jsonapi"
	PrometheusTaskSource TaskSource = "prometheus"
)

type TaskMethod string

const (
	CountMethod TaskMethod = "count"
	SumMethod   TaskMethod = "sum"
	AndMethod   TaskMethod = "and"
	OrMethod    TaskMethod = "or"
)

type TaskAuth struct {
	Header   string `yaml:"header,omitempty" json:"header,omitempty"`
	TokenVar string `yaml:"tokenVar,omitempty" json:"tokenVar,omitempty"`
}

type TaskResult struct {
	Result string // Result of the task
}

type Task struct {
	ID        string   `yaml:"id,omitempty" json:"id,omitempty"`
	Name      string   `yaml:"name,omitempty" json:"name,omitempty"`
	Type      string   `yaml:"type,omitempty" json:"type,omitempty"`
	DependsOn []string `yaml:"dependsOn,omitempty" json:"dependsOn,omitempty"`

	// Extract related fields
	Source string `yaml:"source,omitempty" json:"source,omitempty"`

	// Extract related fields for REST API calls
	URI             string    `yaml:"uri,omitempty" json:"uri,omitempty"`
	JSONPath        string    `yaml:"jsonPath,omitempty" json:"jsonPath,omitempty"`
	Auth            *TaskAuth `yaml:"auth,omitempty" json:"auth,omitempty"`
	PrometheusQuery string    `yaml:"prometheusQuery,omitempty" json:"prometheusQuery,omitempty"`

	// Extract related fields for GitHub API calls
	Repo         string `yaml:"repo,omitempty" json:"repo,omitempty"`
	FilePath     string `yaml:"filePath,omitempty"`
	SearchString string `yaml:"searchString,omitempty" json:"searchString,omitempty"`

	// Validate related fields
	Rule    string `yaml:"rule,omitempty" json:"rule,omitempty"`
	Pattern string `yaml:"pattern,omitempty" json:"pattern,omitempty"`

	// Aggregate related fields
	Method string `yaml:"method,omitempty" json:"method,omitempty"`

	// Run related fields
	Result       interface{}     `yaml:"-" json:"-"`
	Dependencies []*Task         `yaml:"-" json:"-"` // List of tasks this task depends on
	DoneCh       chan TaskResult `yaml:"-" json:"-"` // Channel to signal task completion
}

func (t1 *Task) IsEqual(t2 *Task) bool {
	if t2 == nil {
		return false
	}

	return t1.ID == t2.ID &&
		t1.Name == t2.Name &&
		t1.Type == t2.Type &&
		t1.Source == t2.Source &&
		t1.URI == t2.URI &&
		t1.JSONPath == t2.JSONPath &&
		t1.Auth == t2.Auth &&
		t1.Repo == t2.Repo &&
		t1.FilePath == t2.FilePath &&
		t1.Rule == t2.Rule &&
		t1.Pattern == t2.Pattern &&
		t1.Method == t2.Method &&
		t1.Result == t2.Result &&
		t1.SearchString == t2.SearchString &&
		t1.PrometheusQuery == t2.PrometheusQuery &&
		t1.IsDependsOnEquals(t2.DependsOn)
}

func (t1 *Task) IsDependsOnEquals(dependsOn []string) bool {
	if len(t1.DependsOn) != len(dependsOn) {
		return false
	}
	for i := 0; i < len(t1.DependsOn); i++ {
		if t1.DependsOn[i] != dependsOn[i] {
			return false
		}
	}

	return true
}
