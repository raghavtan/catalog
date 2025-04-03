# Security Grading System

## Security as Pipeline
Ensure We have early detection security alerts at CI/CD stage before going live with production deployment.

### Validations
  - Validate if the CI/CD pipeline runs the security step invoking the `motain/onefootball-actions/security` action
    ```yaml
    # To pass this validation the pipeline needs to invoke the motain/onefootball-actions/security action
    - id: security-checks
      uses: motain/onefootball-actions/security@master
      with:
        token: ${{ github.token }}
        path: "."
        image-url: ${{ steps.fetch-release-candidate.outputs.image-url }}
    ```
[<< Back to the index](./index.md)

## Vulnerability Management

Resource (in focus) should be able to enforce condition to stay at latest patch when severity/bugs are found and also remediation options in place when threats are discovered. This will focus specifically on state of the software codebase.

### Validations
- Validate that the are no Critical vulnerabilities for the namespace matching the component name.
To fix any vulenerability follow the instruction at ... <We need to validate where to find these information>
  ```promql
  # To double-check this value run this query in the Explore of Grafana and ensure there are no vulnerabiities.
  sum(trivy_image_vulnerabilities{namespace="<component-name>", severity="Critical" })
  ```

[<< Back to the index](./index.md)
