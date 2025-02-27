# Onefootball Catalog (OFC)

FC is a Go-based command-line interface (CLI) designed to manage the grading system, featuring an integrated scraper that processes data points either in batches or individually for each component, populating scorecards. It adheres to the Unix principle of idempotency and modularity.

For a detailed documenation about features and definitions refers the [docs](./docs/index.md).

## Installation

**Prerequisites**

Ensure you have:

 - Go 1.20+ installed ([Download Go](https://go.dev/dl/))

**Clone the Repository**

```bash
git clone https://github.com/example/fact-collector.git
cd fact-collector
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

**Generate Wire Dependencies**

[Wire](https://github.com/google/wire) is a dependency injection tool for Go that eliminates the need for manually wiring dependencies. It generates code that initializes dependencies automatically based on provider functions.

**How Wire Works:**

1. Define an Interface for Each Service

    Every service that should be injected needs an interface.

    ```go
    type RepositoryInterface interface {
        FetchData() string
    }
    ```

2. Provide a Concrete Implementation

    A struct implements the interface.

    ```go
    type Repository struct {
        config *ConfigService
    }

    func (r *Repository) FetchData() string {
        return "data"
    }
    ```

3. Bind the Interface to the Implementation

    Wire needs an explicit binding to know which struct fulfills the interface.

    ```go
    var ProviderSet = wire.NewSet(
        NewRepository,
        wire.Bind(new(RepositoryInterface), new(*Repository)),
    )
    ```

4. Use Dependency in a Constructor

    Wire ensures dependencies are injected when calling constructors.

    ```go
    type Handler struct {
        repo RepositoryInterface
    }

    func NewHandler(repo RepositoryInterface) *Handler {
        return &Handler{repo: repo}
    }
    ```

5. Generate Dependency Code

    Running the command below will generate the required dependency injection code for the project:

    ```bash
    wire gen ./internal/app
    ```

    Wire will build the dependency chain, resolving the correct constructors automatically.

To wire all the subcommands at once you can also run:

```bash
  make wire-all
```


## Configuration

Usage

`To be defined`

Example:

```bash
go run main.go ...
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

## Project Structure

```bash
fact-collector/
├── cmd                                 # This is the root of the command
│   └── root.go                         # All modules are called here
├── go.mod
├── go.sum
├── internal
│   ├── app                             # Wire related folder
│   │   ├── wire.go                     # Register here the dependencies
│   │   └── wire_gen.go                 # Do not edit manually !!!
│   ├── modules                         # All modules are defined here
│   │   └── metric                      # metric is a module
│   │       ├── handler                 # Handlers expose functionalities
│   │       │   ├── handler.go
│   │       │   └── handler_test.go
│   │       └── repository              # Repositories handle data
│   │           ├── repository.go
│   │           └── repository_test.go
│   └── services                        # Services are common resources
│       ├── configservice
│       │   ├── config.go
│       │   └── config_test.go
│       ├── githubservice
│       │   └── github.go
│       └── keyringservice
│           └── keyring.go
├── main.go                             # TBR | point to cmd/root instead
└── tests                               # Functional tests
    └── cli_test.go
```

## Contribution

This repository follows a trunk-based development workflow. To contribute:

1. Clone the repository.

    ```
    git clone https://github.com/motain/fact-collector.git
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
