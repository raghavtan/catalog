# Security Grading System

## Security as Pipeline
Ensure We have early detection security alerts at CI/CD stage before going live with production deployment.

**Validations**

- Validate if the CI/CD pipeline runs the security step invoking the `motain/onefootball-actions/security` action

    ```yaml
    # To pass this validation the pipeline needs to invoke the motain/onefootball-actions/security action
    # Use security Module: https://github.com/motain/onefootball-actions/tree/master/security
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

**Validations**

- Validate that the are no Critical vulnerabilities for the namespace matching the component name.

  ```promql
  # To double-check this value run this query in the Explore of Grafana and ensure there are no vulnerabiities.
  sum(trivy_image_vulnerabilities{namespace="<component-name>", severity="Critical" })
  ```

### Howto - Fix Guidelines


To know more about the other security issues and fix vulnerabilities follow the instruction below:

**[Visit Dashboard : Grafana - Trivy Vulnerabilty](https://grafana.mgm.onefootball.com/d/ycwPj724k/trivy-operator-dashboard?orgId=1&from=now-3h&to=now&timezone=browser&var-DS_PROMETHEUS=P0F161AC36DE6FE17&var-namespace=cloud-runtime-bot&refresh=5m&tab=transformations&viewPanel=panel-33)**

**Steps:**
- Select the `namespace` for the services
- Review the list defined which contains all the details of the package impacted. Please upgrade to `fixed_version` to eliminate and patch the vulnerability. 

<img width="1858" alt="Screenshot 2025-06-13 at 16 24 46" src="https://github.com/user-attachments/assets/dcc45196-9b00-401f-bfbe-79c7b3a363f5" />

--- 
 
**(Optional) Additonal CRDs available in Kubernetes for more details:**

Add references on solving software vulnerability detected by the security tools
  ```bash
  # login to aws production in your terminal
  kubectx # chose prod cluster
  kubectl get vulns -n <component-name> -ojson | jq '[.items[].report.vulnerabilities | unique_by(.vulnerabilityID) | .[]]'
  ```

[<< Back to the index](./index.md)
