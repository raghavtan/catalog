# GitHub

## GitHub Authentication

To interact with GitHub's APIs, a valid API token is required. This token can be provided in one of two ways:

1. **Environment Variable Injection**
    Recommended for containerized environments.
    Set the `GITHUB_TOKEN` environment variable.

2. **Retrieval from Local Keyring**
    On macOS, the token is retrieved from the Keychain.
    When running locally, it's recommended to set the `GITHUB_USER` environment variable.
    This helps access the keyring and retrieve a (hopefully) valid token.

[<- back to index](./index.md)
