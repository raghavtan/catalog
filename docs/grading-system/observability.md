# Observability Grading System

## Instrumentation Check

Verifies that the component is set up to record relevant telemetry (metrics or traces).

### Validations
- Validate that `OTEL_SERVICE_NAME` is defined and not null in `app.toml` under `[envs]`, `[env.production]`, `[service.envs]` or `[service.production.envs]`. It should be a string that begins with the component's name.
  ```toml
  # Ensure OTEL_SERVICE_NAME is defined as in the following example
  [service.production.envs]
  ...
  # OTEL_SERVICE_NAME = <component-name>([opt]<my-environemnt>)
  OTEL_SERVICE_NAME = clips-service_production
  ...
  ```

- Validate that `OTEL_RESOURCE_ATTRIBUTES` is defined and not null in `app.toml` under `[envs]`, `[env.production]`, `[service.envs]`, or `[service.production.envs]`. The value must include at least one of `of.sample_rate` or `of.error_sample_rate`.
  ```toml
  # Ensure OTEL_RESOURCE_ATTRIBUTES is defined as in the following examples
  [service.production.envs]
  ...
  OTEL_RESOURCE_ATTRIBUTES = "of.tail_sampling=refinery,of.sample_rate=100000,of.error_sample_rate=1000"
  # or
  OTEL_RESOURCE_ATTRIBUTES =  "of.tail_sampling=refinery,of.sample_rate=100000"
  # or
  OTEL_RESOURCE_ATTRIBUTES =  "of.tail_sampling=refinery,of.error_sample_rate=1000"
  ...
  ```

- ⚠️ Validate that **ALL** the previous validations succeed

[<< Back to the index](./index.md)

## Observability Documentation

Checks if at least one critical performance or reliability metric is actively monitored and triggers an alert (e.g., SLO or metric threshold).

### Validations
- Validate that the observability.md file exists in the docs folder
  ```bash
  # Ensure the observability.md file exists under docs.
  $ tree ./docs
  ./docs
  ├──...
  ├── observability.md
  ├── ...
  ```

  > ⚠️ This validation will change in the future to validate the file content
  > Be an early adopter by proactively taking action, include a section for each SLO.
  >
  > # Observability documentation
  > ## SLO One
  > ### Alerts
  > The SLO triggers an alarm when it reaches a critical treshold notifying the slack channel Y
  > ###
  > Notification target
  > #my-slack-channel
  > This SLO monitor the performance of the golden path X which is critical to guarantee ...
  >
  > ...
  >

[<< Back to the index](./index.md)

## Critical Alerts SLO Check

### Validations
- For each SLO defined in Honeycomb's dataset that matches the component name, extract the associated alerts and validate that the number of alerts is equal to the number of SLOs.
  ```hcl
  # Ensure you define the the SLIs, SLOS and Triggers using Terraform.
  # Ensure the trigger is correctly defined
  resource "honeycombio_trigger" "a_trigger_for_a_failing_slo" {
    alert_type  = "on_change"
    dataset     = local.service_name
    description = "This alert triggers if a SLO fails."
    disabled    = "false"
    frequency   = "900" // - 15 mins
    name        = "[${local.service_name}] Failed to ingest a premium video"
    query_id    = honeycombio_query.a_query.id

    recipient {
      target = "#a-slack-channel"
      type   = "slack"
    }

    threshold {
      exceeded_limit = "1"
      op             = ">"
      value          = "0"
    }
  }
  ```
[<< Back to the index](./index.md)

## Alert Routing and Notifications

Checks if at least one critical performance or reliability metric is actively monitored and triggers an alert (e.g., SLO or metric threshold).

### Validations
- For each SLO defined in Honeycomb's dataset that matches the component name, extract the associated alerts, for each alert validate that the recipients targets are not empty strings.
  ```hcl
  # Ensure you define the the SLIs, SLOS and Triggers using Terraform.
  # Ensure the trigger is correctly defined
  resource "honeycombio_trigger" "a_trigger_for_a_failing_slo" {
    ...

    recipient {
      target = "#a-slack-channel"
      type   = "slack"
    }

    ..

  }
  ```

[<< Back to the index](./index.md)
