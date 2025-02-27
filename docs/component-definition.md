# Component Catalog Overview

## YAML Structure Explanation
Each component definition follows a standard structure with the following properties:

```yaml
---
apiVersion: of-catalog/v1alpha1
kind: Component
metadata:
    name: <string>
    componentType: <string> # service, cloud-resource, application
dpendsOn: [] # TBD
spec: # Mirroring Compass.yaml file
  name: <string> # matching metadata.name
  description: <string>
  typeId: <string> # SERVICE, CLOUD_RESOURCE, APPLICATION
  ownerId: <string>
  fields:
    lifecycle: <string>
    tier: <number> # 1, 2, 3, 4
  links: <list>
    - name: <string>
      type: <string> # CHAT_CHANNEL, DASHBOARD, DOCUMENT, REPOSITORY, ON_CALL, OTHER_LINK
      url: <string>
  labels: [<string>]
```

### 1. API Version & Kind

```yaml
apiVersion: of-catalog/v1alpha1
kind: Component
```

- **Type:** `apiVersion` (string), `kind` (string)
- **Description:** Defines the API version and resource type (Component).

### 2. Metadata

```yaml
metadata:
  name: comp-test-1
  componentType: service
```

- **Type:** `metadata` (object)
  - `name` (string) - Unique identifier of the component.
  - `componentType` (string) - The type of component (e.g., service).

### 3. Dependencies - TBD

```yaml
dependsOn: []
```

- **Type:** `dependsOn` (array)
- **Description:** Lists other components this service depends on (empty by default).

### 4. Specification (spec)

```yaml
spec:
  name: comp-test-1
  description: null
  configVersion: 1
  typeId: SERVICE
  ownerId: null
```

- **Type:** `spec` (object)
  - `name` (string) - Component display name.
  - `description` (string, nullable) - Brief description.
  - `configVersion` (integer) - Schema version.
  - `typeId` (string) - Type of component (SERVICE).
  - `ownerId` (string, nullable) - Identifies the owner.

### 5. Fields

```yaml
fields:
  lifecycle: Active
  tier: 4
```

- **Type:** `fields` (object)
  - `lifecycle` (string) - Current lifecycle status (e.g., Active).
  - `tier` (integer) - Tier level indicating priority or importance.

### 6. Links

```yaml
links:
  - name: Repository
    type: REPOSITORY
    url: "https://github.com/motain/comp-test-1"
  - name: Monitoring page
    type: DASHBOARD
    url: "https://grafana.mgm.onefootball.com"
```

- **Type:** `links` (array of objects)
  - `name` (string) - Display name of the link.
  - `type` (string) - Type of the link (e.g., REPOSITORY, DASHBOARD).
  - `url` (string) - URL of the resource.

### 7. Labels

```yaml
labels: ["aws", "service"]
```

- **Type:** `labels` (array of strings)
- **Description:** Keywords or tags associated with the component.

[<- back to index](./index.md)
