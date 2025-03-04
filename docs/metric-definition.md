# Metric Catalog Overview

## YAML Structure Explanation

Each metric definition follows a standard structure with the following properties:

```yaml
---
apiVersion: of-catalog/v1alpha1
kind: Metric
metadata:
  name: <string>
  labels:
    grading-system: production-readiness
  componentType: <list> # service, cloud-resource
  facts:
    all: # Boolean conditions (all must be true)
      - source: <string>
        uri: <string>
        name: <string>
        repo: <string>
        branch: <string>
        factType: <string>
        filePath: <string>
        regexPattern: <string>
        jsonPath: <string>
        repoProperty: <string>
        reposSearchQuery: <string>
        expectedValue: <string>
        expectedFormula: <string>
        auth:
          header: <string>
          tokenEnvVariable: <string>
    any: # Boolean conditions (at least one must be true)
      - source: <string>
        uri: <string>
        name: <string>
        repo: <string>
        branch: <string>
        factType: <string>
        filePath: <string>
        regexPattern: <string>
        jsonPath: <string>
        repoProperty: <string>
        reposSearchQuery: <string>
        expectedValue: <string>
        expectedFormula: <string>
        auth:
          header: <string>
          tokenEnvVariable: <string>
    inspect:
      - source: <string>
        uri: <string>
        name: <string>
        repo: <string>
        factType: <string>
        jsonPath: <string>
        auth:
          header: <string>
          tokenEnvVariable: <string>
```

### 1. API Version & Kind

```yaml
apiVersion: of-catalog/v1alpha1
kind: Metric
```

- **Type**: `apiVersion` (string), `kind` (string)
- **Description**: Defines the API version and resource type (Metric).

### 2. Metadata

```yaml
metadata:
  name: metric-test-1
  labels:
  componentTypes:
  - service
  - cloud-resource
```

- **Type**: `metadata` (object)
  - **name** (string) - Unique identifier of the metric.
  - **labels** (object) - Tags for classification (e.g., grading-system).
  - **componentType** (list) - The type of component being evaluated.

### 3. Facts (Evaluation Criteria)

```yaml
facts:
  all:
    - source: "github"
      factType: fileExists
      filePath: app.toml
  any:
    - source: "jsonapi"
      uri: "https://pokeapi.co/api/v2/pokemon/ditto"
      factType: jsonPath
      jsonPath: weight
      expectedFormula: "> 50"
```

- **Type**: `facts` (object)
  - **all** (array) - List of conditions that must all be met.
  - **any** (array) - List of conditions where at least one must be met.

### 4. Specification (spec)

```yaml
spec:
  name: metric-test-1
  description: example metric-1 used to test cli
  format:
  unit: msg/sec
```

- **Type**: `spec` (object)
  - **name** (string) - Metric display name. **Follow kebaba notation**.
  - **description** (string) - Brief description of the metric.
  - **format.unit** (object) - Defines measurement unit. Metrics will always be sent as float, and this field only provides contextual information about the expected unit of measurement.

## Usage:

### Facts Configuration

The facts section defines conditions and inspections for validating data from different sources (github or jsonapi). It consists of three parts:

**`all` (AND conditions)**

A list of facts that must all be true for validation to pass.
Evaluated using a lazy AND (if one condition fails, evaluation stops).

**`any` (OR conditions)**

A list of facts where at least one must be true for validation to pass.
Evaluated using a lazy OR (stops when the first true condition is found).

**Evaluation Order**

If both all and any are present:
- all is evaluated first.
- If all is false, any is not evaluated.
- If all is true, any is evaluated, and the final result is a lazy AND between the two groups.

**`inspect` (Data Extraction Only)**

Defines a single fact used to extract data from a remote source but not for validation.

**Supported Sources**

GitHub (source: github)
- Uses the repo field (repository name without owner).
- The uri field is ignored.

JSON API (source: jsonapi)
- The uri field contains the URL of the JSON API endpoint.
- The repo field is ignored.

### Fact Types
| Fact Type        | Source          | Description                                                   |
|------------------|-----------------|---------------------------------------------------------------|
| fileExists       | GitHub          | Checks if a file exists in the repository.                    |
| fileRegex        | GitHub          | Checks if a file contains a regex pattern.                    |
| jsonPath     | GitHub, JSON API| Checks if a file or API response contains a JSON path.        |
| repoProperties   | GitHub          | Checks if a repository has a specific property.               |
| reposSearchQuery | GitHub          | Checks if a string exists in a repository.                    |
| expectedValue    | GitHub, JSON API| Compares a JSON path's value with an expected value.          |
| expectedFormula  | GitHub, JSON API| Compares a JSON path's value using an expected formula.       |


---

**JSONAPI Type**

The JsonAPI Type enables making `GET` requests to endpoints that return a JSON response.
In many cases, when gathering or verifying information about a specific component, the upstream server requires authentication. This collector allows you to configure authentication, with the only supported method being via a custom request header.
To set up authentication when defining the fact, specify the following options:
```yaml
auth:
  header: Authorization
  tokenEnvVariable: MY_SECRET_API_KEY
```

Additionally, ensure an environment variable is set with the name specified in `tokenEnvVariable`. The collector will retrieve this variable's value and include it in the request header.

---

**Placeholders**

In repo and expectedValue, placeholders can be used to dynamically reference values from the component the metric is bound to.
Format: `${<component-json-path>}`
The placeholder is replaced with the value of the specified JSON path.

---

### `metadata.componentTypes`
The `componentTypes` field enables dynamic binding of the metric to components. The metric will be linked to all components whose `componentType` matches any value listed in this field.

### `spec.expectedFormula`
The `expectedFormula` field allows for dynamic comparisons of extracted values using mathematical expressions. When a value is extracted (e.g., from `jsonPath` or `repoProperty`), it is appended as a prefix to the expectedFormula string and then evaluated as a logical expression.

The extracted value (e.g., 45) is prepended to the formula (e.g., "> 50"), forming the expression:

```math
45 > 50
```

The function returns true if the condition holds or false otherwise.

Example Scenarios:

    expectedFormula: "> 50" → Passes if extracted value is greater than 50.
    expectedFormula: "!= 0" → Passes if extracted value is not 0.
    expectedFormula: "<= 100" → Passes if extracted value is less than or equal to 100.

This allows for flexible metric validation beyond exact matches (expectedValue).

Allowed values for expectedFormula are: `>=`, `<=`, `>`, `<`, `==`, `!=`.

Note that foe equality the correct formula is `==` and that `=` will result in an execution error.

[<- back to index](./index.md)
