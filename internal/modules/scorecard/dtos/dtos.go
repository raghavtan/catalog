package dtos

type ScorecardDTO struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

func GetScorecardUniqueKey(c *ScorecardDTO) string {
	return c.Spec.Name
}

func FromStateToConfig(state *ScorecardDTO, conf *ScorecardDTO) {
	conf.Spec.ID = state.Spec.ID
}

func IsComponentTypeIDsEqual(t1, t2 []string) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i, componentTypeID := range t1 {
		if componentTypeID != t2[i] {
			return false
		}
	}
	return true
}

func IsCriteriaEqual(c1, c2 []*Criterion) bool {
	if len(c1) != len(c2) {
		return false
	}

	for i, criterion := range c1 {
		if !IsCriterionEqual(criterion, c2[i]) {
			return false
		}
	}
	return true
}

func IsCriterionEqual(s1, s2 *Criterion) bool {
	return s1.HasMetricValue.Name == s2.HasMetricValue.Name &&
		s1.HasMetricValue.Weight == s2.HasMetricValue.Weight &&
		s1.HasMetricValue.MetricName == s2.HasMetricValue.MetricName &&
		s1.HasMetricValue.MetricDefinitionId == s2.HasMetricValue.MetricDefinitionId &&
		s1.HasMetricValue.ComparatorValue == s2.HasMetricValue.ComparatorValue &&
		s1.HasMetricValue.Comparator == s2.HasMetricValue.Comparator
}

func IsScoreCardEqual(s1, s2 *ScorecardDTO) bool {
	return s1.Spec.Name == s2.Spec.Name &&
		s1.Spec.Description == s2.Spec.Description &&
		s1.Spec.OwnerID == s2.Spec.OwnerID &&
		s1.Spec.State == s2.Spec.State &&
		IsComponentTypeIDsEqual(s1.Spec.ComponentTypeIDs, s2.Spec.ComponentTypeIDs) &&
		s1.Spec.Importance == s2.Spec.Importance &&
		s1.Spec.ScoringStrategyType == s2.Spec.ScoringStrategyType &&
		IsCriteriaEqual(s1.Spec.Criteria, s2.Spec.Criteria)
}

type Metadata struct {
	Name string `yaml:"name"`
}

type Spec struct {
	ID                  *string      `yaml:"id"`
	Name                string       `yaml:"name"`
	Description         string       `yaml:"description"`
	OwnerID             string       `yaml:"ownerId"`
	State               string       `yaml:"state"`
	ComponentTypeIDs    []string     `yaml:"componentTypeIds"`
	Importance          string       `yaml:"importance"`
	ScoringStrategyType string       `yaml:"scoringStrategyType"`
	Criteria            []*Criterion `yaml:"criteria"`
}

type Criterion struct {
	HasMetricValue MetricValue `yaml:"hasMetricValue"`
}

func FromStateCriteriaToConfig(state *Criterion, conf *Criterion) {
	conf.HasMetricValue.ID = state.HasMetricValue.ID
}

type MetricValue struct {
	ID                 string `yaml:"id"`
	Weight             int    `yaml:"weight"`
	Name               string `yaml:"name"`
	MetricName         string `yaml:"metricName"`
	MetricDefinitionId string `yaml:"metricDefinitionId"`
	ComparatorValue    int    `yaml:"comparatorValue"`
	Comparator         string `yaml:"comparator"`
}
