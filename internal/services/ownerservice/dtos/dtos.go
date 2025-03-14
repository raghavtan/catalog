package dtos

type GroupList []*Group

type Group struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Links       []Link              `yaml:"links"`
	Annotations AnnotationsMetadata `yaml:"annotations"`
}

type Link struct {
	URL   string `yaml:"url"`
	Title string `yaml:"title"`
	Type  string `yaml:"type"`
	Icon  string `yaml:"icon"`
}

type AnnotationsMetadata struct {
	JiraTeamID string `yaml:"jiraTeamID"`
}

type Spec struct {
	ID       string   `yaml:"id"`
	Profile  Profile  `yaml:"profile"`
	Type     string   `yaml:"type"`
	Parent   string   `yaml:"parent"`
	Children []string `yaml:"children"`
}

type Profile struct {
	DisplayName string `yaml:"displayName"`
}

type Owner struct {
	OwnerID       string            `yaml:"ownerID"`
	SlackChannels map[string]string `yaml:"slackChannel"`
	Projects      map[string]string `yaml:"projects"`
	DisplayName   string            `yaml:"displayName"`
}
