package dtos

type ComponentDTO struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name          string `yaml:"name"`
	ComponentType string `yaml:"componentType"`
}

type Spec struct {
	ID            *string  `yaml:"id"`
	Name          string   `yaml:"name"`
	Slug          string   `yaml:"slug"`
	Description   string   `yaml:"description"`
	ConfigVersion int      `yaml:"configVersion"`
	TypeID        string   `yaml:"typeId"`
	OwnerID       string   `yaml:"ownerId"`
	Links         []Link   `yaml:"links"`
	Labels        []string `yaml:"labels"`
}

type Link struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	URL  string `yaml:"url"`
}
