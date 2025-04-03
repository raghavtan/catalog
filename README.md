# Onefootball Catalog (OFC)

FC is a Go-based command-line interface (CLI) designed to manage the grading system, featuring an integrated scraper that processes data points either in batches or individually for each component, populating scorecards. It adheres to the Unix principle of idempotency and modularity.

For a detailed documenation about features and definitions refers the [docs](./docs/index.md).

## Installation

**Prerequisites**

Ensure you have:

 - Go 1.20+ installed ([Download Go](https://go.dev/dl/))

**Clone the Repository**

```bash
git clone https://github.com/motain/of-catalog.git
cd of-catalog
```

**Install Dependencies**

```bash
go mod tidy
```
This will download and install any missing dependencies listed in go.mod.

**Adding New Dependencies**

```bash
go get <dependency>
```
or
```bash
go get <dependency>@<version>
```

**Verifying Dependencies**

To ensure all dependencies are correct and match their checksums, run:

```bash
go mod verify
```

**Achieving Decoupling and Maintainability with DI & IoC**

To build scalable and maintainable software, we aim for loose coupling, where components interact with minimal dependencies on each other. One way to achieve this is through Dependency Injection (DI), a design pattern that shifts the responsibility of creating and managing dependencies from within a class to an external source. Instead of hardcoding dependencies, they are "injected" from the outside, making the system more modular, testable, and easier to modify.

Inversion of Control (IoC) takes this idea further by reversing the traditional flow of control. Instead of a class managing its dependencies, an external framework or container takes over, handling object creation and lifecycle. DI is a key implementation of IoC, ensuring that dependencies are managed efficiently and reducing tight coupling between components.

Tools like _Wire_, a compile-time dependency injection framework for Go, help automate this process by generating code to wire dependencies together, simplifying configuration and improving efficiency.

For more about how to use wire in this project refer the the [wire page](./docs/wire.md).

## Configuration

Usage

Before using it locally copy the content of the `of-catalog-env-file` Bitwarden note into a `.env` file in the root of the project.
You will need to adjust the `GITHUB_USER` entry to match your GitHub user. If you never used the `gh` cli install it and follow the [quickstart](https://docs.github.com/en/github-cli/github-cli/quickstart) to set up your local environemnt.

Example:

```bash
go run  go run ./cmd/root.go
```

## Running Tests

**Unit Tests**

```bash
go test ./internal/...
```

**Functional (Black Box) Tests**

```bash
go test ./tests/...
```

## Project Architecture Overview

The project is structured around a root command, which initializes and organizes subcommands. Each subcommand represents a specific module, encapsulating all related logic and functionality.
Module Structure

A module consists of the following components:

- Commands – The controller layer, responsible for triggering the DI framework, validating input, and calling the appropriate handler.
- Handlers – The logic layer, orchestrating service calls to retrieve or store data from the source of truth (repository), which in our case is- Compass.
- Repositories – Abstract interactions with the source of truth, ensuring data consistency and separation of concerns.
- Resources – Core domain objects representing business entities.
- DTOs (Data Transfer Objects) – Used for reading and writing definitions, ensuring structured data exchange.
- Services – Abstractions for external resources commonly used across modules.
- Utils – Collections of lightweight, self-contained functions that do not interact with third-party services and are simple enough to not require- mocking in tests.

Module Encapsulation & Dependencies

All services and functions within a module should be used only internally. The only resources that other modules may access are DTOs and utils, ensuring a clean and modular architecture with well-defined boundaries.


## Contribution

This repository follows a trunk-based development workflow. To contribute:

1. Clone the repository.

    ```
    git clone https://github.com/motain/of-catalog.git
    cd go-cli-app
    ```

2. Make your changes and commit them.

  ```bash
  git commit -sam "feat(module-name): feature" -m 'short feature description'
  ```

  or

  ```bash
  git commit -sam "fix(module-name): fix" -m 'short fix description'
  ```

  or other commit types.

3. Push directly to the main branch (if allowed) or merge via a PR.

    ```bash
    git push origin main
    ```

If using a feature branch, merge it back into main following internal guidelines.

## License

This project is for internal use only. Unauthorized copying, distribution, or modification of this codebase is strictly prohibited.
