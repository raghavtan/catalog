# Bind Command Overview

The `bind` subcommand of `metric` is used to associate a metric definitions with component definitions. It retrieves the identifier for each definition and creates a corresponding `MetricSource` one. This definition is stored in a state file (metricsource.yaml), which contains all necessary details for tracking the relationship between metrics and components.

## Identifier Format

The identifier for each `MetricSource` is generated as:
```
<metric-identifier>-<component-identifier>
```

## Fact Processing

- The command extracts fact definitions from the metric and processes them to generate `MetricSource` facts.
- If a metric fact contains dynamic placeholders (e.g., `${<json.path>}`), they are replaced with the corresponding values from the component definition.
- This substitution applies specifically to the `repo` and `expectedValue` properties.

## Example

If a metric fact includes `${spec.name}`, it will be replaced with the corresponding value from `Component.spec.name` in the `MetricSource`.

This ensures that metrics are dynamically adapted to each component they are bound to.

[<- back to index](./index.md)
