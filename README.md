# Go Project Template

This repository provides example code for setting up an empty Go project following best practices.
For more information on recommended project structure, please look at this info [Golang Standards Project Layout](https://github.com/golang-standards/project-layout).

# External Packages
The template project uses the following external packages:
- A CLI based on the [Cobra](https://github.com/spf13/cobra) framework. This framework is used by many other Golang project, e.g Kubernetes, Docker etc.
- Logging via [Logrus](https://github.com/sirupsen/logrus)

## Project structure 
- **`cmd/`**  
  Contains the application's main executable logic. This is where `main.go` lives.

- **`internal/`**  
  Contains private application code. Anything inside `internal/` cannot be imported from outside the project (enforced by the Go compiler).

- **`pkg/`**  
  Contains public libraries or utilities that can be imported by other projects if needed.

- **`bin/`**  
  Stores built binaries, generated via the `Makefile`.

- **`Makefile`**  
  Automates build tasks such as compiling, building Docker images, running tests, etc.

- **`go.mod` and `go.sum`**  
  Define module requirements and manage dependencies.

##  Quick Start
### Build the project
```bash
go mod tidy
make build
```

### Run the binary
```bash
./bin/helloworld talk
```

or type:
```bash
go run cmd/main.go talk
```

```console
ERRO[0000] Error detected                                Error="This is an error"
INFO[0000] Talking...                                    Msg="Hello, World!" OtherMsg="Logging is cool!"
Hello, World!
```


### Build Docker container
```bash
make container
```

### Running tests
To run all tests;
```bash
make test 
```
Remember to update Makefile if adding more source directories with tests.

To run individual tests:
```bash
cd pkg/helloworld
go test -v --race
```

Always use the `--race` flag when running tests to detect race conditions during execution.  
The `--race` flag enables the Go race detector, helping you catch concurrency issues early during development.

```bash
 cd pkg/helloworld/
make test 
```

## Change the Project Name
To customize the project name, follow these steps:

1. Create a new Git repo.

2. **Edit `go.mod`** 
   Change the module path from: `module github.com/eislab-cps/go-template` to your new project path.

3. **Update Import Paths**  
Modify the import paths in the following files:

- `internal/cli/version.go`  
  Line 6:
  ```go
  "github.com/eislab-cps/go-template/pkg/build"
  ```

- `internal/cli/talk.go`  
  Line 4:
  ```go
  "github.com/eislab-cps/go-template/pkg/helloworld"
  ```

- `cmd/main.go`  
  Lines 4–5:
  ```go
  "github.com/eislab-cps/go-template/internal/cli"
  "github.com/eislab-cps/go-template/pkg/build"
  ```

Replace each instance of `github.com/eislab-cps/go-template` with your new module name.

## Continus 

## Other tips
- Run `go mod tidy` to clean up and verify dependencies.
- To store all dependencies in the `./vendor` directory, run:

  ```sh
  go mod vendor
  ```
- Github Co-pilot is very good at generating logging statement.

