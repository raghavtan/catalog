package dtos

import "sort"

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
		IsEqualLinks(c1.Spec.Links, c2.Spec.Links) &&
		IsEqualLabels(c1.Spec.Labels, c2.Spec.Labels) &&
		IsEqualDependsOn(c1.Spec.DependsOn, c2.Spec.DependsOn) &&
		IsEqualFields(c1.Spec.Fields, c2.Spec.Fields)
}

func SortAndRemoveDuplicateDocuments(documents []*Document) []*Document {
	if len(documents) == 0 {
		return documents
	}
	type docKey struct {
		title   string
		url     string
		docType string
	}

	uniqueDocs := make(map[docKey]*Document)
	for _, doc := range documents {
		key := docKey{
			title:   doc.Title,
			url:     doc.URL,
			docType: doc.Type,
		}
		uniqueDocs[key] = doc
	}

	result := make([]*Document, 0, len(uniqueDocs))
	for _, doc := range uniqueDocs {
		result = append(result, doc)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Title < result[j].Title
	})

	return result
}

// UniqueAndSortLinks ensures that a slice of Links contains unique elements
// and is sorted by Type, then Name, then URL.
func UniqueAndSortLinks(links []Link) []Link {
	if len(links) == 0 {
		return links
	}

	// Ensure uniqueness
	// Using a map with a composite key (Type + Name + URL) for uniqueness.
	// If links have reliable IDs from Compass that should define uniqueness,
	// this logic might need adjustment, but for mixed sources (Compass + derived),
	// content-based key is safer.
	type linkKey struct {
		linkType string
		name     string
		url      string
	}
	uniqueMap := make(map[linkKey]Link)

	for _, l := range links {
		// Normalize or ensure consistency if needed, e.g., trimming spaces, lowercasing for key
		key := linkKey{
			linkType: l.Type,
			name:     l.Name,
			url:      l.URL,
		}
		// If multiple links have the same key, the last one encountered wins.
		// If IDs are present and should be preserved, prioritize the one with an ID,
		// or decide on a merging strategy if content is identical but IDs differ (unlikely for this problem).
		// For now, simple last-one-wins for the content key.
		if existing, found := uniqueMap[key]; found {
			// If existing link has an ID and current one doesn't, keep existing.
			// Or if current one has an ID, it might be an update from Compass.
			// This part can get complex if we need to merge intelligently.
			// Simplest for now: if an ID exists, it's likely from Compass and authoritative.
			if l.ID != "" { // Current link has an ID, prefer it.
				uniqueMap[key] = l
			} else if existing.ID == "" && l.ID == "" { // Neither has ID, last one wins.
                 uniqueMap[key] = l
            }
            // If existing has ID and current does not, existing is kept (implicitly).
		} else {
			uniqueMap[key] = l
		}
	}

	result := make([]Link, 0, len(uniqueMap))
	for _, l := range uniqueMap {
		result = append(result, l)
	}

	// Sort the unique links
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type != result[j].Type {
			return result[i].Type < result[j].Type
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].URL < result[j].URL
	})

	return result
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
