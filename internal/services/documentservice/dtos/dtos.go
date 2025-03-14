package dtos

type Document struct {
	SiteName        string    `yaml:"site_name" json:"site_name"`
	SiteDescription string    `yaml:"site_description" json:"site_description"`
	RepoURL         string    `yaml:"repo_url" json:"repo_url"`
	Nav             []NavItem `yaml:"nav" json:"nav"`
	Plugins         []string  `yaml:"plugins" json:"plugins"`
}

type NavItem struct {
	Title    string    `yaml:"-" json:"title"`
	File     string    `yaml:"-" json:"file"`
	SubItems []NavItem `yaml:"-" json:"children"`
}

func (n *NavItem) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data map[string]interface{}
	if err := unmarshal(&data); err == nil {
		for k, v := range data {
			n.Title = k
			n.unmarshalRecursive(v)
		}
		return nil
	}

	return nil
}

func (n *NavItem) unmarshalRecursive(data interface{}) {
	switch subItems := data.(type) {
	case string:
		n.File = subItems
	case []interface{}:
		for _, subItem := range subItems {
			subItemMap := subItem.(map[string]interface{})
			for subK, subV := range subItemMap {
				child := NavItem{Title: subK}
				child.unmarshalRecursive(subV)
				n.SubItems = append(n.SubItems, child)
			}
		}
	}
}
