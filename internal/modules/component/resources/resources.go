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
	Labels        []string
	CustomFields  interface{}
	MetricSources map[string]*MetricSource
}

type Link struct {
	Name string
	Type string
	URL  string
}

type MetricSource struct {
	ID     string
	Name   string
	Metric string
}
