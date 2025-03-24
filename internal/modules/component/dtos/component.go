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
	conf.Spec.OwnerID = state.Spec.OwnerID
}

func IsEqualLinks(l1, l2 []Link) bool {
	if len(l1) != len(l2) {
		return false
	}

	linkMap := make(map[Link]bool)
	for _, link := range l1 {
		linkMap[link] = true
	}

	// CHAT_CHANNEL is populated when applying. We must not consider diff on it
	for _, link := range l2 {
		if link.Type != "CHAT_CHANNEL" && !linkMap[link] {
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

func IsEqualDependsOn(d1, d2 []string) bool {
	if len(d1) != len(d2) {
		return false
	}

	for i, label := range d1 {
		if label != d2[i] {
			return false
		}
	}
	return true
}

func IsEqualFields(f1, f2 map[string]interface{}) bool {
	if len(f1) != len(f2) {
		return false
	}

	for k, v := range f1 {
		if f2[k] != v {
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
		// IsEqualLinks(c1.Spec.Links, c2.Spec.Links) &&
		IsEqualLabels(c1.Spec.Labels, c2.Spec.Labels) &&
		IsEqualDependsOn(c1.Spec.DependsOn, c2.Spec.DependsOn) &&
		IsEqualFields(c1.Spec.Fields, c2.Spec.Fields)
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
	DependsOn     []string                    `yaml:"dependsOn" json:"dependsOn"`
	Fields        map[string]interface{}      `yaml:"fields" json:"fields"`
	Links         []Link                      `yaml:"links" json:"links"`
	Documents     []*Document                 `yaml:"documents" json:"documents"`
	Labels        []string                    `yaml:"labels" json:"labels"`
	MetricSources map[string]*MetricSourceDTO `yaml:"metricSources" json:"metricSources"`
	Tribe         string                      `yaml:"tribe" json:"tribe"`
	Squad         string                      `yaml:"squad" json:"squad"`
}

type Link struct {
	ID   string `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
	URL  string `yaml:"url" json:"url"`
}

type Document struct {
	ID                      string `yaml:"id" json:"id"`
	Title                   string `yaml:"title" json:"title"`
	Type                    string `yaml:"type" json:"type"`
	DocumentationCategoryId string `yaml:"documentationCategoryId" json:"documentationCategoryId"`
	URL                     string `yaml:"url" json:"url"`
}
