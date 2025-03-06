package dtos

type FactOperations struct {
	All     Facts `yaml:"all"`
	Any     Facts `yaml:"any"`
	Inspect *Fact `yaml:"inspect"`
}

func (fo *FactOperations) IsEqual(f2 FactOperations) bool {
	return fo.All.IsEqual(f2.All) &&
		fo.Any.IsEqual(f2.Any) &&
		fo.Inspect.IsEqual(f2.Inspect)
}

// FactType defines the type of fact to collect
type FactType string

const (
	FileExistsFact     FactType = "fileExists"
	FileRegexFact      FactType = "fileRegex"
	JSONPathFact       FactType = "jsonPath"
	RepoPropertiesFact FactType = "repoProperties"
	RepoSearchFact     FactType = "repoSearch"
)

type FactSource string

const (
	GitHubFactSource    FactSource = "github"
	JSONAPIFactSource   FactSource = "jsonapi"
	ComponentFactSource FactSource = "component"
)

// Fact defines a struct to handle different facts
type Fact struct {
	Source           string    `yaml:"source"`           // Source of the fact (e.g., "github")
	URI              string    `yaml:"uri"`              // URI of the fact
	ComponentName    string    `yaml:"componentName"`    // Component name
	Name             string    `yaml:"name"`             // Name of the fact
	Repo             string    `yaml:"repo"`             // Repository name (e.g., "repo")
	FactType         FactType  `yaml:"factType"`         // Type of fact to collect
	FilePath         string    `yaml:"filePath"`         // File to open/validate fact (if FactType is "fileExists", "fileRegex", or "jsonPath")
	RegexPattern     string    `yaml:"regexPattern"`     // Regex to match file content or response (if FactType is "fileRegex")
	JSONPath         string    `yaml:"jsonPath"`         // JSONPath to navigate file content or json response (if FactType is "jsonPath")
	RepoProperty     string    `yaml:"repoProperty"`     // Property to explore in the repo (if FactType is "repoProperties")
	ReposSearchQuery string    `yaml:"reposSearchQuery"` // Query to search for repositories (if FactType is "repoSearch")
	ExpectedValue    string    `yaml:"expectedValue"`    // Expected value (matched against value of  "repoProperty" or "jsonPath")
	ExpectedFormula  string    `yaml:"expectedFormula"`  // Expected formula (matched against value of  "repoProperty" or "jsonPath")
	Auth             *FactAuth `yaml:"auth"`             // Auth to use to access the fact when using a URI require authentication
}

type Facts []*Fact

func (facts Facts) IsEqual(f2 []*Fact) bool {
	if len(facts) != len(f2) {
		return false
	}

	for i, fact := range facts {
		if !fact.IsEqual(f2[i]) {
			return false
		}
	}

	return true
}

func (f *Fact) IsEqual(f2 *Fact) bool {
	if f2 == nil {
		return false
	}

	return f.Source == f2.Source &&
		f.URI == f2.URI &&
		f.ComponentName == f2.ComponentName &&
		f.Name == f2.Name &&
		f.Repo == f2.Repo &&
		f.FactType == f2.FactType &&
		f.FilePath == f2.FilePath &&
		f.RegexPattern == f2.RegexPattern &&
		f.JSONPath == f2.JSONPath &&
		f.RepoProperty == f2.RepoProperty &&
		f.ReposSearchQuery == f2.ReposSearchQuery &&
		f.ExpectedValue == f2.ExpectedValue &&
		f.ExpectedFormula == f2.ExpectedFormula &&
		f.Auth.IsEqual(f2.Auth)
}

type FactAuth struct {
	Header           string `yaml:"header"`
	TokenEnvVariable string `yaml:"tokenEnvVariable"`
}

func (a *FactAuth) IsEqual(a2 *FactAuth) bool {
	if a2 == nil {
		return false
	}

	return a.Header == a2.Header &&
		a.TokenEnvVariable == a2.TokenEnvVariable
}
