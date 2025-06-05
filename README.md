# Onefootball Catalog (OFC)

FC is a Go-based command-line interface (CLI) designed to manage the grading system, featuring an integrated scraper that processes data points either in batches or individually for each component, populating scorecards. It adheres to the Unix principle of idempotency and modularity.

For a detailed documentation about features and definitions refers the [docs](./docs/index.md).

## Table of Contents

- [Usage](#usage) ⭐ **Most Important**
   - [Building the Application](#building-the-application)
   - [Managing Dependencies](#managing-dependencies)
   - [Code Quality](#code-quality)
   - [Applying Changes](#applying-changes)
- [Development Workflow](#development-workflow)
- [Project Architecture](#project-architecture-overview)
- [Testing](#running-tests)

<details>
<summary><strong>Prerequisites</strong> (Click to expand)</summary>

Ensure you have the following installed:

- **Go 1.20+** ([Download Go](https://go.dev/dl/))
- **Git** for version control
- **GitHub CLI (`gh`)** for authentication ([Installation Guide](https://github.com/cli/cli#installation))

</details>

<details>
<summary><strong>Local Setup</strong> (Click to expand)</summary>

### 1. Clone and Setup Repository

```bash
# Clone the repository
git clone https://github.com/motain/of-catalog.git
cd of-catalog

# Install dependencies
make vendor

# Generate required code
make generate
make wire-all
```

### 2. Environment Configuration

Before using the application locally:

1. Copy the content from the `of-catalog-env-file` Bitwarden note
2. Create a `.env` file in the project root
3. Paste the content and adjust the `GITHUB_USER` entry to match your GitHub username

If you haven't used the GitHub CLI before:
```bash
# Install gh CLI and authenticate
gh auth login
```

Follow the [GitHub CLI quickstart](https://docs.github.com/en/github-cli/github-cli/quickstart) for detailed setup instructions.

</details>

<details>
<summary><strong>Configuration</strong> (Click to expand)</summary>

The application uses configuration files located in:
- `./config/grading-system/` - Metrics definitions
- `./config/scorecard/` - Scorecard definitions
- Component configurations are managed through individual scripts

</details>

## Usage

### Building the Application

```bash
# Build Linux binary
make build
```

### Applying Changes

> ⚠️ **Important**: Always follow the correct order when applying changes to avoid dependency issues.

#### Order of Operations

1. **Metrics** (must be applied first)
2. **Scorecards** (depends on metrics)
3. **Components** (depends on scorecards)
4. **Binding** (must be done last)

#### Individual Operations

```bash
# 1. Apply metrics changes
make create-metrics

# 2. Apply scorecard changes  
make create-scorecards

# 3. Apply component changes
make create-components

# 4. Bind components to grading system (ALWAYS LAST)
make bind-components
```

#### Full Setup (New Environment or Clean Atlassian Compass Setup)

For a complete fresh setup or major changes:

```bash
# Clean state if needed (DESTRUCTIVE - requires confirmation)
make clean-state

# Create everything in the correct order
make create-all
```

The `create-all` target will:
1. Prompt you about prerequisites
2. Apply metrics → scorecards → components → bindings in the correct order
3. Provide progress feedback throughout the process

#### Making Changes

**For Metrics:**
1. Edit files in `./config/grading-system/`
2. Run `make create-metrics`

**For Scorecards:**
1. Edit files in `./config/scorecard/`
2. Run `make create-scorecards`
3. If metrics were also changed, run both in order:
   ```bash
   make create-metrics
   make create-scorecards
   ```

**For Components:**
1. Update component configurations
2. Run `make create-components`
3. Always rebind after component changes:
   ```bash
   make create-components
   make bind-components
   ```

**For Major Changes:**
If you've modified multiple types (metrics, scorecards, and components):
```bash
make create-all
```

## Development Workflow

### Managing Dependencies

```bash
# Update all dependencies and vendor them
make update-deps

# Only vendor existing dependencies
make vendor
```

### Code Quality

```bash
# Run linter
make lint

# Run all tests with coverage
make test

# Run tests for specific component
make stest C=path/to/component

# Generate coverage report
make test/coverage
```

### Available Commands

Run `make help` to see all available targets with descriptions.

### Testing Your Changes

```bash
# Run all tests
make test

# Test specific component
make stest C=internal/component-name

# Check code quality
make lint
```

### Code Generation

```bash
# Generate code using go generate
make generate

# Run wire for dependency injection
make wire-all
```

## Project Architecture Overview

The project is structured around a root command, which initializes and organizes subcommands. Each subcommand represents a specific module, encapsulating all related logic and functionality.

### Module Structure

A module consists of the following components:

- **Commands** – The controller layer, responsible for triggering the DI framework, validating input, and calling the appropriate handler.
- **Handlers** – The logic layer, orchestrating service calls to retrieve or store data from the source of truth (repository), which in our case is Compass.
- **Repositories** – Abstract interactions with the source of truth, ensuring data consistency and separation of concerns.
- **Resources** – Core domain objects representing business entities.
- **DTOs (Data Transfer Objects)** – Used for reading and writing definitions, ensuring structured data exchange.
- **Services** – Abstractions for external resources commonly used across modules.
- **Utils** – Collections of lightweight, self-contained functions that do not interact with third-party services and are simple enough to not require mocking in tests.

### Module Encapsulation & Dependencies

All services and functions within a module should be used only internally. The only resources that other modules may access are DTOs and utils, ensuring a clean and modular architecture with well-defined boundaries.

### Achieving Decoupling and Maintainability with DI & IoC

To build scalable and maintainable software, we aim for loose coupling, where components interact with minimal dependencies on each other. One way to achieve this is through Dependency Injection (DI), a design pattern that shifts the responsibility of creating and managing dependencies from within a class to an external source. Instead of hardcoding dependencies, they are "injected" from the outside, making the system more modular, testable, and easier to modify.

Inversion of Control (IoC) takes this idea further by reversing the traditional flow of control. Instead of a class managing its dependencies, an external framework or container takes over, handling object creation and lifecycle. DI is a key implementation of IoC, ensuring that dependencies are managed efficiently and reducing tight coupling between components.

Tools like _Wire_, a compile-time dependency injection framework for Go, help automate this process by generating code to wire dependencies together, simplifying configuration and improving efficiency.

For more about how to use wire in this project refer the the [wire page](./docs/wire.md).

## Running Tests

**Unit Tests**

```bash
make test
```

**Functional (Black Box) Tests**

```bash
go test ./tests/...
```

**Test Specific Component**

```bash
make stest C=internal/component-name
```

**Coverage Report**

```bash
make test/coverage
```
