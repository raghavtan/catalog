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
