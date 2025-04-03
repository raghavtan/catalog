# Production Readiness Grading System

## Deployment Readiness
Ensures the component is correctly configured, adheres to best practices, and meets all necessary requirements for secure, stable, and efficient deployment in a production environment.

### Validations
- Validate that `replicas_min` and `replicas_max` are defined and not null in `app.toml` under `[service.production]` or `[service]`. To succeed validate that:
  - `replicas_min` is not equals to `replicas_max`
  - `replicas_min` is greater than or equals to 3
  - `replicas_max` is greater than to 3
  ```toml
  # Ensure replicas_min and replicas_max are properly defined
  [service.production]
  ...
  replicas_min = 3
  replicas_max = 6
  ...
  ```

[<< Back to the index](./index.md)

## Organizational Standards

Ensures the component aligns with defined organizational standards, including coding guidelines and compliance requirements.

### Validations
- Validate that the README file exists as a non empty file as either one of the following options:
  - ./README.md
  - docs/README.md
  - docs/index.md

  ```bash
  $ tree ./ | grep README
  ├── README.md
  # or
  $ tree ./docs | grep README
  ├── README.md
  # or
  $ tree ./docs | grep index
  ├── index.md

- Validate if the CI/CD pipeline runs the PaaS step invoking the `motain/onefootball-actions/paas-deploy` action
    ```yaml
    # To pass this validation the pipeline needs to invoke the motain/onefootball-actions/paas-deploy action
    - id: deploy-production
      uses: motain/onefootball-actions/paas-deploy@master
      with:
        tag: ${{ github.sha }}
        release: false
        pr-key: ${{ secrets.PAAS_CREATE_AND_MERGE_TO_HELM_CHARTS }}
    ```

- ⚠️ Validate that **ALL** the previous validations succeed

[<< Back to the index](./index.md)
