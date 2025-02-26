package resources

type Scorecard struct {
	ID                  *string
	Name                string
	Description         string
	OwnerID             string
	State               string
	ComponentTypeIDs    []string
	Importance          string
	ScoringStrategyType string
	Criteria            []*Criterion `yaml:"criteria"`
}

type Criterion struct {
	HasMetricValue MetricValue `yaml:"hasMetricValue"`
}

type MetricValue struct {
	ID                 string
	Weight             int
	Name               string
	MetricDefinitionId string
	MetricName         string
	ComparatorValue    int
	Comparator         string
}
