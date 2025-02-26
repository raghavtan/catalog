package resources

// MetricFormat represents the format of a metric, including its unit.
type MetricFormat struct {
	Unit string // Unit is the unit of measurement for the metric.
}

// Metric represents a metric with its details.
type Metric struct {
	ID          string       // ID is the unique identifier of the metric.
	Name        string       // Name is the name of the metric.
	Description string       // Description provides details about the metric.
	Format      MetricFormat // Format specifies the format of the metric.
}
