# Metric Module

This document provides an overview of the module functionality and command options, serving as a guide for using the metric module effectively.

The metric module is responsible for managing service catalogue metrics. A metric is a computed measurement derived from a component's data, used to track and evaluate its state and performance, and is defined using [YAML files](../metric-definition.md).


## Overview

- **Purpose:**
  The module manages the lifecycle of metrics by handling their configurations, synchronizing state with a remote IDP, and ensuring that computed values are properly applied.

- **Resource Definitions:**
  - Metrics are defined as YAML files.
  - Configuration files follow the naming convention:
  `metric(.*).yaml`
  - The file name prefix identifies the kind of resources.
  - The module filters resources by matching the `Kind` property.

- **State Management:**
  The module creates resources on the remote IDP and saves the enriched definition (including the ID retrieved from the IDP) in a local state file.
  The overall flow is similar to the component module:

```
############             ###########             #############
#  CONFIG  #  ========>  #   IDP   #  ========>  #   STATE   #
############             ###########             #############
```

- **File Conventions:**
- **Configuration Files:** Can be centralized or spread among multiple files.
- **State File:** Holds all resource definitions in one single file per Kind. Filenames are lowercase while Kind names are in Pascal Case, and resources in the state file are sorted alphabetically by `Metadata.Name`.

## Facts
Metrics define a list of facts following the fact system syntax. For more information refer to the fact system documentation.

Read the documentation for more information regarding the [fact system](../fact-system/overview.md).

## Dynamic placeholders
When defining metrics it's possible to specify dynamic placeholder that are evenutally processed and replaced by the [component bind command](./component.md#bind).
Dynamic placeholders are processed for the Fact Properties:
- URI
- Repo
- Pattern
- PrometheusQuery

## Command

The metric module exposes a single command: **Apply**.

### Apply

The `apply` command synchronizes configuration files with the state file by detecting differences between them.

- **Workflow:**
1. Load all resource definitions from both configuration and state.
2. Match metrics by `Metadata.Name` and determine one of four scenarios:
   - **New Resource:** Exists in configuration but not in state.
     → Create in the remote IDP, retrieve the identifier, and store the enriched definition in the state file.
   - **Unchanged Resource:** Exists in both configuration and state and are identical.
     → Rewrite the configuration into the state file without remote action.
   - **Modified Resource:** Exists in both but differ.
     → Refresh the resource on the remote IDP and update the state file.
   - **Deleted Resource:** Found in state but missing in configuration.
     → Delete the resource from the remote IDP and remove it from the state file.
      - If the resource is missing on the remote IDP, the error is ignored and the state is updated.

- **Command Options:**
```
  -l, --configRootLocation string   Root location of the config
  -h, --help                        help for apply
  -r, --recursive                   Apply changes recursively
```

- The **configRootLocation** is required and can be either a full or relative path.
- Use the **recursive** flag if configuration files are stored in subfolders.
