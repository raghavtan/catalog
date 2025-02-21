package resources

type Component struct {
	ID            *string
	Name          string
	Slug          string
	Description   string
	ConfigVersion int
	TypeID        string
	OwnerID       string
	Links         []Link
	Labels        []string
	CustomFields  interface{}
}

type Link struct {
	Name string
	Type string
	URL  string
}
