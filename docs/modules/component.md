# Component Module

This document provides an overview of the module functionality and command options, serving as a guide for using the component module effectively.

The component module is responsible for managing service catalogue components. A component represents a service, cloud resource, or application, and is defined using [YAML files](../component-definition.md).

## Overview

- **Purpose:**
  The module manages the lifecycle of components by handling their configurations, synchronizing state with a remote IDP, and pairing them with metrics.

- **Resource Definitions:**
  - Components are defined as YAML files.
  - Configuration files follow the naming convention:
  `component(.*).yaml` (The file name prefix identifies the kind of resources.)
  - The module filters resources by matching the `Kind` [property](../component-definition.md#1-api-version--kind).

- **State Management:**
  The module creates resources on the remote IDP and saves the enriched definition (with an ID from the IDP) in a local state file.
  The overall flow is:

```
############             ###########             #############
#  CONFIG  #  ========>  #   IDP   #  ========>  #   STATE   #
############             ###########             #############
```


- **File Conventions:**
  - **Configuration Files:** Can be centralized or spread among multiple files.
  - **State File:** Holds all resource definitions in one single file per Kind. Filenames are lowercase while Kind names are in Pascal Case, and resources in the state file are sorted alphabetically by `Metadata.Name`.

## Commands

The component module exposes three primary commands: **Apply**, **Bind**, and **Compute**.

### Apply

The `apply` command synchronizes configuration files with the state file by detecting drifts between them.

- **Workflow:**
1. Load all resource definitions from both configuration and state.
2. Match resources by `Metadata.Name` and determine one of four scenarios:
   - **New Resource:** Exists in configuration but not in state.
     → Create in the remote IDP, retrieve the identifier, and store the enriched definition in the state file.
     - If the reource already exists in the remote IDP the pattern is to retrieve such entity and to enrich definition in the state with the retrieved identifier.
     - If any step fails, the resource definition is not flushed in the state file.
   - **Unchanged Resource:** Exists in both configuration and state and are identical.
     → Rewrite configuration into the state file without remote action.
   - **Modified Resource:** Exists in both but differ.
     → Refresh the resource on the remote IDP and update the state file.
   - **Deleted Resource:** Found in state but missing in configuration.
     → Delete the resource from the remote IDP and remove it from the state file.
      - If the resource is missing on the remote IDP, the error is ignored and the state is updated.

### Components relations
Componetns defines a list of components on which they depends on, this can be used to create a dependency graph. Dependecies are recalculated for each run for resources that are new, updated, or did not change. The only scenario where these are not refreshed is when a compoennt is removed from the configuration.

### Documenation
The documentation list does not come with the component definition but is build  fetched from the component repository looking for an `mkdocs` file.
Similarly to what happens with the dependencies, the documentation list get refreshed with every apply for compoenents that were not removed.

### API Specification
The api specification does not come with the component definition but is build  fetched from the component repository looking for `openapi` or `swagger` files.
Similarly to what happens with the dependencies and documentation list, the api specification get refreshed with every apply for compoenents that were not removed.

The relationships between components, documentation, and the API specification are the reasons for overwriting the state file for components that haven't changed. Since these details are recalculated with each run and are not specified in the configuration, changes outside the scope of the configuration can occur, resulting in the state being refreshed.

- **Command Options:**

```
-c, --component           string  Name of the component
-l, --configRootLocation  string  Root location of the config
-h, --help                        Help for apply
-r, --recursive                   Apply changes recursively
```

- The **configRootLocation** is required and can be either a full or relative path.
- Use the **recursive** flag if configuration files are stored in subfolders.
- To apply changes to a specific component, pass the `--component` flag with the component's name.
  If no matching resource is found, the command exits with a failing status code of 1.

### Bind

The `bind` command is used to pair metrics with components. Remote IDPs needs to match metrics with components to store data.

- **Workflow:**
1. Parse both component and metric definitions.
2. Match metrics to components.
3. If needed, create a resource in the remote IDP and store its identifier in the state.
4. Encapsulate the metric definition into the component for fast retrieval during computation (similar to NoSQL database denormalization).

- **Dynamic Placeholders:**
Metrics may include dynamic placeholders (e.g., `${Spec.Name}`) that are replaced by the corresponding component values (like `Component.Spec.Name`).
This replacement is performed during bind rather than at compute time to reduce processing overhead during metrics computation, although this may increase disk and memory usage.

### Compute

The `compute` command processes metrics for a component, computing facts and pushing values to the remote IDP.

Read the documentation for more information regarding the [fact system](../fact-system/overview.md).

- **Command Options:**

```
-a, --all                Compute all metrics for the component
-c, --component  string  Name of the component
-h, --help               Help for compute
-m, --metric     string  Name of the metric
```
- **Usage Scenarios:**
- **Compute a Single Metric:**
  ```bash
  compute --component simple-service --metric organizational-standards
  ```
- **Compute All Metrics:**
  ```bash
  compute --component simple-service --all
  ```


## GitHub Workflow
To compute all the metrics for a component run the GitHub workflow [ComputeComponentMetrics](https://github.com/motain/of-catalog/actions/workflows/compute-component-metrics.yaml)


## Batch trigger apply
To trigger the refresh of all the services run:
``` bash
$ ./apply-all-individually.sh
```

This is the same as run:

```bash
 go run ./cmd/root.go component apply -l ./config/components
```

The first option is safer though as it apply adn flush to the state one component at the time. In the second option, if the process fails, it can happen that the state is not persisted.

## Batch trigger compute
To trigger the refresh of all the services run:
``` bash
$ ./compute-all.sh
```
Similarly to the `apply-all-individually.sh` script, this script compute all the metric for all the components.
