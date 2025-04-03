# Resiliency Grading System

## High Availability
Ensure the component is properly configured to support High Availability, which keeps a system or service running with minimal downtime. This helps maintain continuous operation, even during failures or unexpected issues.

### Validations
- Validate that `target_cpu_utilization_percentage` is defined and not null in `app.toml` under `[service.production]` or `[service]`. The validation should pass only if its value is greater than or equal to 20.
  ```toml
  # Ensure target_cpu_utilization_percentage is properly defined
  [service.production]
  ...
  target_cpu_utilization_percentage = 60
  ...
  ```

- Validate that `target_memory_utilization_percentage` is defined and not null in `app.toml` under `[service.production]` or `[service]`. The validation should pass only if its value is greater than or equal to 20.
  ```toml
  # Ensure target_memory_utilization_percentage is properly defined
  [service.production]
  ...
  target_memory_utilization_percentage = 60
  ...
  ```

- ⚠️ Validate that **AT LEAST ONE** of the previous validations succeed

[<< Back to the index](./index.md)

## Adaptive Systems

Ensure the component is properly configured to support Adaptive Systems, which can automatically adjust to changing conditions. This helps maintain optimal performance, efficiency, and reliability without the need for manual intervention.

### Validations
- Validate that `replicas_min` and `replicas_max` are defined and not null in `app.toml` under `[service.production]` or `[service]`. To succeed validate that:
  - replicas_min is less than replicas_max
  - replicas_min is greater than or equals to 3
  ```toml
  # Ensure replicas_min and replicas_max are properly defined
  [service.production]
  ...
  replicas_min = 3
  replicas_max = 6
  ...
  ```

[<< Back to the index](./index.md)
