package dtos

type ComponentDTO struct {
	APIVersion string   `yaml:"apiVersion" json:"apiVersion"`
	Kind       string   `yaml:"kind" json:"kind"`
	Metadata   Metadata `yaml:"metadata" json:"metadata"`
	Spec       Spec     `yaml:"spec" json:"spec"`
}

func GetComponentUniqueKey(c *ComponentDTO) string {
	return c.Spec.Name
}

func FromStateToConfig(state *ComponentDTO, conf *ComponentDTO) {
	conf.Spec.ID = state.Spec.ID
	conf.Spec.MetricSources = state.Spec.MetricSources
}

func IsEqualLinks(l1, l2 []Link) bool {
	for i, link := range l1 {
		if link.Name != l2[i].Name {
			return false
		}

		if link.Type != l2[i].Type {
			return false
		}

		if link.URL != l2[i].URL {
			return false
		}
	}
	return true
}

func IsEqualLabels(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}

	for i, label := range l1 {
		if label != l2[i] {
			return false
		}
	}
	return true
}

func IsEqualComponent(c1, c2 *ComponentDTO) bool {
	return c1.Spec.Name == c2.Spec.Name &&
		c1.Spec.Description == c2.Spec.Description &&
		c1.Spec.ConfigVersion == c2.Spec.ConfigVersion &&
		c1.Spec.TypeID == c2.Spec.TypeID &&
		c1.Spec.OwnerID == c2.Spec.OwnerID &&
		IsEqualLinks(c1.Spec.Links, c2.Spec.Links) &&
		IsEqualLabels(c1.Spec.Labels, c2.Spec.Labels)
}

type Metadata struct {
	Name          string `yaml:"name" jsonyaml:"name"`
	ComponentType string `yaml:"componentType" jsonyaml:"componentType"`
}

type Spec struct {
	ID            string                      `yaml:"id" json:"id"`
	Name          string                      `yaml:"name" json:"name"`
	Slug          string                      `yaml:"slug" json:"slug"`
	Description   string                      `yaml:"description" json:"description"`
	ConfigVersion int                         `yaml:"configVersion" json:"configVersion"`
	TypeID        string                      `yaml:"typeId" json:"typeId"`
	OwnerID       string                      `yaml:"ownerId" json:"ownerId"`
	Links         []Link                      `yaml:"links" json:"links"`
	Labels        []string                    `yaml:"labels" json:"labels"`
	MetricSources map[string]*MetricSourceDTO `yaml:"metricSources" json:"metricSources"`
}

type Link struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
	URL  string `yaml:"url" json:"url"`
}
