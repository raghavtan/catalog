package resources

type Component struct {
	ID            string
	Name          string
	Slug          string
	Description   string
	ConfigVersion int
	TypeID        string
	OwnerID       string
	Links         []Link
	Documents     []Document
	Labels        []string
	CustomFields  interface{}
	MetricSources map[string]*MetricSource
}

type Link struct {
	ID   string
	Name string
	Type string
	URL  string
}

type Document struct {
	ID                      string
	Title                   string
	Type                    string
	DocumentationCategoryId string
	URL                     string
}

type MetricSource struct {
	ID     string
	Name   string
	Metric string
}
