# Cost Optimization Grading System

# Allocation Efficiency

Verifies if the component has defined resource allocations in an efficient way.

### Validations
- Validate `cpu_requests` and `memory_requests` are defined and not null in `app.toml` under `[service]` or `[service.production]`
  ```toml
  # Ensure cpu_requests and memory_requests are defined in [service] or [service.production]
  [service.production]
  ...
  cpu_requests = "256Mi"
  memory_requests = "512Mi"
  ...
  ```

- Validate `cpu_limits` is not defined in `app.toml` under `[service]` nor `[service.production]`
  ```toml
  # Ensure the cpu_limits is not set neither in [service] nor in [service.production]
  [service.production]
  ...
  # cpu_limits ="512Mi"
  ...
  ```

- Validate `memory_limits` is defined and not null in `app.toml` under `[service]` or `[service.production]`
  ```toml
  # Ensure memory_limits is defined in [service] or [service.production]
  [service.production]
  ...
  memory_limits = "256Mi"
  ...
  ```

- ⚠️ Validate that **ALL** the previous validations succeed

[<< Back to the index](./index.md)
