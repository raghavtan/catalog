# Definitions
- [Component](./component-definition.md)
- [Metric](./metric-definition.md)
- [Scorecard](./scorecard-definition.md)

# Commands
- [Bind](./bind-command.md)

# Environment Variables
# Definitions
- [Component](./component-definition.md)
- [Metric](./metric-definition.md)
- [Scorecard](./scorecard-definition.md)

# Commands
- [Bind](./bind-command.md)

# Generic
- [GitHub](./github.md)

# Environment Variables
Environment variables are fetched in the order:
1. `.env` file: if a `.env` file is found;
2. `environment` as second and option;
3. if not found, in some cases the program can apply default values.

Onefootball catalog search for the following variables:
- **GITHUB_ORG**: The GitHub organization name (default: `motain`).
- **GITHUB_TOKEN**: The GitHub token used sent with each api call to GitHub.
- **GITHUB_USER**: The GitHub username associated with the authentication token used for API calls.
- **COMPASS_TOKEN**: The authentication token for performing CRUD operations in Compass.
- **COMPASS_HOST**: The Compass host domain (without protocol).
- **COMPASS_CLOUD_ID**: A unique identifier for the Compass organization.

[<- back to index](./../README.md)
