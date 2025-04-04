# Fact System

This document provides an overview the fact system functionality, serving as a guide for using the component module effectively.

To understand why to facts let's first review the why and the what are metrics and scorecards.

## Why Metrics Matter in an Internal Developer Portal

### Technical Perspective
Metrics in an Internal Developer Portal provide quantifiable data points that help teams objectively measure various aspects of software development and operations. They:

- Enable data-driven decision making instead of relying on intuition
- Provide historical context for observing trends and patterns over time
- Facilitate standardization across different teams and projects
- Create a foundation for automation and programmatic actions
- Support implementation of SLOs (Service Level Objectives) and KPIs
- Allow for benchmarking against industry standards or internal targets

### Strategic Benefits
Beyond the technical implementation, metrics deliver compelling value by:

- Making invisible work visible, highlighting efforts that might otherwise go unnoticed
- Creating shared understanding of priorities across engineering and business teams
- Driving cultural improvement by focusing attention on what matters
- Reducing friction in cross-team collaboration with common measurement frameworks
- Empowering teams with self-service insights rather than relying on centralized reporting
- Accelerating learning cycles by providing faster feedback on initiatives

Scorecards aggregate these metrics into meaningful collections that tell a coherent story about your software ecosystem's health, performance, and compliance posture.

It's important to define metrics that matter for each team and organization. Metric might at some point become obsolete as teams and processes evolves they might not be insightful.
For example, applying automation or fitness functions a metric can be forced to reach a specific standard for all services.

## Metric source

Metrics can represents many things, some examples are:
 - avarage cpu utilazation, mesured in %;
 - configuration correctness (boolean);
 - number of critical vulenerabilities (no unit).

It seems that potentially metrics can have differenet types, indeed this is a wrong approach, the goal of metrics is to report the a number for which the changes can be observed over time. Instead of visualizing the correctness of the configuration we could count the number of missconfiguration. It will be responsibility of the scorecards to match the result against a desired treshhold.

## The fact system
Given the huge amount of possibile sources and operation that can be applied to the findings, of-catalog allow to define operations as pipelines so that it's not necessary to write code for each metric.
The pipeline is itself the fact system and it's composed by three types of facts:
- Extractors: used to collect data
- Validators: used to validate extracted data
- Aggregators: used to aggregate extracted data, validations results, and aggregation results.

### Processor
At the heart of the fact system is the **processor**, which is responsible for coordinating the handling of facts by directing each one to its appropriate handler.

Since facts may depend on the results of other facts, the processor ensures they are executed in the correct order—waiting for dependencies to complete and return their results before proceeding.

Facts are processed concurrently using lightweight threads. However, failures are ignored by design, which can result in missing outputs from dependencies. This may trigger a cascading failure throughout the pipeline, potentially leading to no metric values being returned at all.

Facts are generic objects, but certain properties are specific to components within the fact system. The processor primarily relies on the following:

  - `id`: Uniquely identifies each fact.
  - `type`: Determines which handler should process the fact.

### Extractors
The goal of extractors is to fetch data from remote sources. These sources are defined in the property `source` and include:

- **GitHub**: Fetches data from GitHub repositories.
- **JSON API**: Fetches data from generic hosts that return JSON responses.
- **Prometheus**: Fetches data from AWS AMP.

Each source hander accept specific rules and configuration that are used to handle the request to the remote service.

### GitHub Source

The GitHub source handles the following properties:

- `repo`: Repository to use for queries.
- `filePath`: Files to fetch.
- `jsonPath`: JSON path to apply to results.
- `searchString`: String to search in the repository.
- `rule`: Rule to apply.

**Rule behaviors for this source:**

- **jsonpath**: Applies the JSON path defined in the `jsonPath` property.
- **notempty**: Validates that the response is not empty, returning a boolean.
- **search**: Searches for the given string in the repository.
- **no rule**: If no rule is specified, returns the raw content.

---

### JSON API Source

The JSON API source handles the following properties:

- `uri`: The URI to query.
- `jsonPath`: JSON path to apply to results.
- `rule`: Rule to apply.
- `auth`:
  - `header`: Header to send for authorizing the request.
  - `tokenVar`: Environment variable name used to retrieve the token for the header.
    *Note: If the remote source expects a prefix such as `Bearer` or `Basic`, it must be included in the value of `tokenVar`.*

**Rule behaviors for this source:**

- **jsonpath**: Applies the JSON path defined in the `jsonPath` property.
- **notempty**: Validates that the response is not empty, returning a boolean.
- **no rule**: If no rule is specified, returns the raw content.

---

### Prometheus Source

The Prometheus source handles the following properties:

- `uri`: The URI to query.
- `jsonPath`: JSON path to apply to results.
- `rule`: Rule to apply.
- `prometheusQuery`: Query to run against the Prometheus server.

**Rule behaviors for this source:**

- **jsonpath**: Applies the JSON path defined in the `jsonPath` property.
- **notempty**: Validates that the response is not empty, returning a boolean.
- **no rule**: If no rule is specified, returns the raw content.

## Validator

The **validator** takes the result from a dependency and applies a specific rule to it. Validator facts return `true` or `false` depending on whether the validation succeeds.

#### Properties

Validators support the following properties:

- `rule`: Defines the rule to apply to the result of the previous task.
- `pattern`: A value or expression used by some rules (e.g., for regex or formulas).

#### Allowed Rules

- **deps_match**:
  Expects the previous task to return a list of primitive items (e.g., strings, integers, floats).
  Returns `true` if all lists have equal length, `false` otherwise.

- **unique**:
  Expects a list of primitive items from the previous task.
  Returns `true` if all values are unique (no repetitions), `false` otherwise.

- **regex_match**:
  Interprets `pattern` as a regular expression.
  Matches the result of the previous task against this pattern.
  Returns `true` if it matches, `false` otherwise.

- **formula**:
  Interprets `pattern` as a partial algebraic expression in the format `<operator> <value>` (e.g., `> 5`).
  The previous task must return a number, which is combined with the pattern (e.g., `4 > 5`).
  Returns `true` if the expression evaluates to true, `false` otherwise.


## Aggregator

An **aggregator** combines the results from one or more dependencies by applying a specified method.

### Properties

Aggregators support the following property:

- `method`: Defines the method used to aggregate the results of the previous tasks.

### Allowed Methods

- **sum**:
  Expects the previous task to return a list of primitive items (e.g., strings, integers, floats).
  Returns the **count** of the items in the list.

- **count**:
  Expects the previous task to return a list of numbers (`int` or `float`).
  Returns the **sum** of all items in the list.

- **and**:
  Expects the previous task to return a list of booleans.
  Returns `true` if **all** items in the list are `true`, `false` otherwise.

- **or**:
  Expects the previous task to return a list of booleans.
  Returns `true` if **at least one** item in the list is `true`, `false` otherwise.

## The fact results
Fact can have dependencies. Let's dive into how the results are handled.

### Executor
## Executor

The **executor** accepts either one dependency or none.

- If **no dependency** is provided, the executor simply applies its own predefined rule.
- If a **dependency** is present, the executor processes its result, which can be:
  - **A list**: Treated as a list of strings. Each item triggers a separate request (or extraction).
  - **A single item**: Treated as a string. Triggers exactly one request (or extraction).

### Dynamic URI and File Path Substitution

When an executor depends on the result of a previous step, certain properties are expected to include placeholders that will be replaced at runtime with values from the previous fact's result:

- For the **JSON API source**, the `uri` property must contain a placeholder.
- For the **GitHub source**, the `filePath` property must contain a placeholder.

For clarity let's review an example:
```yaml
    # Extract SLOs from Honeycomb
    - id: fetch-slos
      name: Fetch SLOs
      type: extract
      source: jsonapi
      uri: https://api.eu1.honeycomb.io/1/slos/${Metadata.Name}
      jsonPath: .[].id
      rule: "jsonpath"
      auth:
        header: X-Honeycomb-Team
        tokenVar: HONEYCOMB_API_KEY
    # Extract Alerts for each SLO from Honeycomb
    - id: fetch-alerts-for-slos
      name: Fetch alerts for SLOs
      type: extract
      dependsOn:
        - fetch-slos
      source: jsonapi
      uri: https://api.eu1.honeycomb.io/1/burn_alerts/${Metadata.Name}?slo_id=:slo_id
      jsonPath: .[].id
      rule: "jsonpath"
      auth:
        header: X-Honeycomb-Team
        tokenVar: HONEYCOMB_API_KEY
```

The fact `fetch-alerts-for-slos` depends on the fact `fetch-slos`, which returns a list of SLO IDs.

For each item in the list, the extractor performs a request by replacing the placeholder `:slo_id` with the individual ID.

> **Note:** The results of these multiple requests are **flattened** into a single, simple list.

---

### Validator

The **validator** accepts **one and only one** dependency.
As described earlier, the result of this dependency can be either a **list** or a **single item**.

---

### Aggregator

The **aggregator** requires **at least one dependency** and will fail if none are provided.
It is designed to operate **only on list results** and will fail if the input is not a list.

For more details on how results are handled, refer to the [Aggregator section](#aggregator).
