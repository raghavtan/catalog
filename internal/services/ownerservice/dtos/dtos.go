package dtos

type GroupList []*Group

type Group struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Links       []Link `yaml:"links"`
}

type Link struct {
	URL   string `yaml:"url"`
	Title string `yaml:"title"`
	Icon  string `yaml:"icon"`
}

type Spec struct {
	ID       string   `yaml:"id"`
	Email    string   `yaml:"email"`
	Profile  Profile  `yaml:"profile"`
	Type     string   `yaml:"type"`
	Parent   string   `yaml:"parent"`
	Children []string `yaml:"children"`
}

type Profile struct {
	DisplayName string `yaml:"displayName"`
}

type Owner struct {
	CompassID    string `yaml:"compassID"`
	Email        string `yaml:"email"`
	SlackChannel string `yaml:"slackChannel"`
	DisplayName  string `yaml:"displayName"`
}
