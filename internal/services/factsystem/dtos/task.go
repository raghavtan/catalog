package dtos

type TaskType string

const (
	AggregateType TaskType = "aggregate"
	ExtractType   TaskType = "extract"
	ValidateType  TaskType = "validate"
)

type TaskRule string

const (
	DepsMatchRule  TaskRule = "deps_match"
	UniqueRule     TaskRule = "unique"
	RegexMatchRule TaskRule = "regex_match"
	FormulaRule    TaskRule = "formula"
)

type TaskSource string

const (
	GitHubTaskSource  TaskSource = "github"
	JSONAPITaskSource TaskSource = "jsonapi"
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

type Task struct {
	ID        string   `yaml:"id,omitempty" json:"id,omitempty"`
	Name      string   `yaml:"name,omitempty" json:"name,omitempty"`
	Type      string   `yaml:"type,omitempty" json:"type,omitempty"`
	DependsOn []string `yaml:"dependsOn,omitempty" json:"dependsOn,omitempty"`

	// Extract related fields
	Source string `yaml:"source,omitempty" json:"source,omitempty"`

	// Extract related fields for REST API calls
	URI      string    `yaml:"uri,omitempty" json:"uri,omitempty"`
	JSONPath string    `yaml:"jsonPath,omitempty" json:"jsonPath,omitempty"`
	Auth     *TaskAuth `yaml:"auth,omitempty" json:"auth,omitempty"`

	// Extract related fields for GitHub API calls
	Repo     string `yaml:"repo,omitempty" json:"repo,omitempty"`
	FilePath string `yaml:"filePath,omitempty"`

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

type TaskResult struct {
	Result string // Result of the task
}
