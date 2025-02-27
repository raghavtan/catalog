# Scorecard Definition Guide

## YAML Structure Explanation

Each scorecard definition follows a structured format with the following properties:

```yaml
---
apiVersion: of-catalog/v1alpha1
kind: Scorecard
metadata:
  name: <string>
spec:
  name: <string>
  description: <string>
  ownerId: <string> # or null
  state: <string> # DRAFT - PUBLISHED
  componentTypeIds: <list> # SERVICE - APPLICATION - WEBSITE - CLOUD_RESOURCE
  importance: <string> # REQUIRED, OPTIONAL
  scoringStrategyType: <string> # WEIGHT_BASED - POINT_BASED
  criteria:
  - hasMetricValue:
    weight: <integer>
    name: <string>
    metricName: <string>
    comparatorValue: <number>
    comparator: <string> # EQUALS - GREATER_THAN - LESS_THAN - GREATER_THAN_OR_EQUAL_TO - LESS_THAN_OR_EQUAL_TO
```

### 1. API Version & Kind

```yaml
apiVersion: of-catalog/v1alpha1
kind: Scorecard
```

- **Type**: `apiVersion` (string), `kind` (string)
- **Description**: Defines the API version and resource type (Scorecard).

### 2. Metadata

```yaml
metadata:
  name: grading-system-x
```

- **Type**: `metadata` (object)
  - `name` (string) - Unique identifier of the scorecard definition.

### 3. Specification (spec)

```yaml
spec:
  name: "grading-system-x"
  description: "This is a test scorecard"
  ownerId: null
  state: "PUBLISHED"
  componentTypeIds: ["SERVICE"]
  importance: REQUIRED
  scoringStrategyType: "WEIGHT_BASED"
```

- **Type**: `spec` (object)
  - `name` (string) - Display name of the scorecard.
  - `description` (string) - Brief description of the scorecard.
  - `ownerId` (string | null) - Owner identifier.
  - `state` (string) - Indicates whether the scorecard is a draft or published.
  - `componentTypeIds` (list) - Defines applicable component types.
  - `importance` (string) - Specifies whether the scorecard is required or optional.
  - `scoringStrategyType` (string) - Defines the scoring method (Weight-based or Point-based).

### 4. Criteria (Evaluation Metrics)

```yaml
criteria:
  - hasMetricValue:
    weight: 60
    name: "msg-per-sec"
    metricName: "metric-test-1"
    comparatorValue: 15
    comparator: "GREATER_THAN"
  - hasMetricValue:
    weight: 40
    name: "foo-bar-baz"
    metricName: "metric-test-2"
    comparatorValue: 100
    comparator: "GREATER_THAN"
```

- **Type**: `criteria` (array)
  - `hasMetricValue` (object) - Defines a metric evaluation rule.
  - `weight` (integer) - Defines the weight of the metric in scoring.
  - `name` (string) - Identifier for the metric condition.
  - `metricName` (string) - References a metric from the catalog.
  - `comparatorValue` (number) - Defines the expected value for comparison.
  - `comparator` (string) - Defines how the extracted value is compared to `comparatorValue`.

### 5. Evaluation Process

The evaluation process follows these rules:

- Each metric condition is checked against its comparator.
- The final score is computed based on the scoring strategy:
  - **Weight-based**: Each metric's contribution is proportional to its weight.
  - **Point-based**: Each metric contributes a fixed point value.

### Allowed Comparators

| Comparator                 | Description              |
|----------------------------|--------------------------|
| `EQUALS`                   | Exact match              |
| `GREATER_THAN`             | Value must be greater    |
| `LESS_THAN`                | Value must be smaller    |
| `GREATER_THAN_OR_EQUAL_TO` | Value must be greater or equal |
| `LESS_THAN_OR_EQUAL_TO`    | Value must be smaller or equal |

This guide provides a structured way to define scorecards and ensures consistency in evaluating components.

[<- back to index](./index.md)
