# Go Static Analysis Workflow

## Overview

The `go_staticAnalysis.yml` workflow performs static code analysis on Go projects using golangci-lint. It supports both standard Go projects and Go workspaces with
multiple modules, providing comprehensive code quality checks, linting, and best practice enforcement.

## Language/Tool Support

- **Go**: All Go versions with golangci-lint compatibility
- **Go Workspaces**: Multi-module workspace support
- **golangci-lint**: Comprehensive Go linter with multiple analyzers
- **Module Detection**: Automatic Go module discovery

## Features

- **Workspace Support**: Automatic detection and analysis of Go workspaces
- **Multi-Module Analysis**: Parallel analysis of multiple Go modules
- **Configurable Timeout**: Adjustable analysis timeout for large projects
- **Selective Analysis**: Option to analyze only new issues
- **Flexible Configuration**: Supports custom golangci-lint configurations

## Usage

```yaml
uses: ./.github/workflows/go_staticAnalysis.yml
with:
  branch: "main"
  goVersion: "1.24"
  timeout: "5m"
  onlyNew: false
  workspace: true
```

## Inputs

| Input       | Type    | Required | Default | Description                                        |
|-------------|---------|----------|---------|----------------------------------------------------|
| `branch`    | string  | ❌        | `main`  | Branch to analyze                                  |
| `goVersion` | string  | ❌        | `1.24`  | Go version to use for analysis                     |
| `timeout`   | string  | ❌        | `1m`    | Timeout for golangci-lint execution                |
| `onlyNew`   | boolean | ❌        | `false` | Analyze only new issues since last run             |
| `workspace` | boolean | ❌        | `false` | Enable Go workspace mode for multi-module analysis |

## Workflow Jobs

### Standard Analysis (workspace: false)

For single-module Go projects:

```yaml
- name: Standard Go Analysis
  uses: golangci/golangci-lint-action@v8
  with:
    only-new-issues: false
    args: --timeout=1m
```

### Workspace Analysis (workspace: true)

For Go workspaces with multiple modules:

1. **Module Detection**: Automatically discovers all Go modules
2. **Parallel Analysis**: Runs golangci-lint on each module in parallel
3. **Matrix Strategy**: Uses GitHub Actions matrix for efficient execution

## Configuration Examples

### Basic Go Project Analysis

```yaml
uses: ./.github/workflows/go_staticAnalysis.yml
with:
  goVersion: "1.24"
  timeout: "2m"
```

### Go Workspace Analysis

```yaml
uses: ./.github/workflows/go_staticAnalysis.yml
with:
  workspace: true
  goVersion: "1.24"
  timeout: "5m"
  onlyNew: false
```

### Pull Request Analysis

```yaml
uses: ./.github/workflows/go_staticAnalysis.yml
with:
  branch: ${{ github.head_ref }}
  onlyNew: true
  timeout: "3m"
```

### Multi-Version Testing

```yaml
name: Multi-Version Analysis
on:
  pull_request:

jobs:
  static-analysis:
    strategy:
      matrix:
        go-version: ['1.22', '1.23', '1.24']
    uses: ./.github/workflows/go_staticAnalysis.yml
    with:
      goVersion: ${{ matrix.go-version }}
```

## Go Workspace Support

### Workspace Structure

The workflow supports Go workspaces with this structure:

```
project/
├── go.work
├── go.work.sum
├── module1/
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── module2/
│   ├── go.mod
│   ├── go.sum
│   └── main.go
└── shared/
    ├── go.mod
    ├── go.sum
    └── utils.go
```

### Workspace Configuration (go.work)

```go
go 1.24

use (
    ./module1
    ./module2
    ./shared
)
```

## golangci-lint Configuration

### Basic Configuration (.golangci.yml)

```yaml
run:
  timeout: 5m
  tests: true
  modules-download-mode: readonly

linters-settings:
  govet:
    check-shadowing: true
  gocyclo:
    min-complexity: 15
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gocyclo
    - dupl
    - goconst
    - gofmt
    - goimports
```

### Advanced Configuration

```yaml
run:
  timeout: 10m
  issues-exit-code: 1
  tests: true
  skip-dirs:
    - vendor
    - third_party
  skip-files:
    - ".*\\.pb\\.go$"
    - ".*_generated\\.go$"

linters-settings:
  revive:
    severity: warning
    rules:
      - name: exported
        severity: error
      - name: indent-error-flow
        severity: warning

  gosec:
    excludes:
      - G404  # Use of weak random number generator
      - G501  # Import blacklist: crypto/md5

  gocritic:
    enabled-tags:
      - performance
      - style
      - experimental
    disabled-checks:
      - wrapperFunc
      - hugeParam

linters:
  enable-all: true
  disable:
    - maligned
    - prealloc
    - gochecknoglobals
    - gochecknoinits

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
        - dupl
    - path: internal/
      text: "exported"
      linters:
        - revive
```

## Integration Patterns

### Pre-commit Analysis

```yaml
name: Code Quality
on:
  pull_request:
    branches: [main, develop]

jobs:
  static-analysis:
    uses: ./.github/workflows/go_staticAnalysis.yml
    with:
      onlyNew: true
      timeout: "3m"

  tests:
    needs: static-analysis
    runs-on: ubuntu-latest
    steps:
      - name: Run Tests
        run: go test ./...
```

### Release Quality Gate

```yaml
name: Release Pipeline
on:
  push:
    tags:
      - 'v*'

jobs:
  quality-gate:
    uses: ./.github/workflows/go_staticAnalysis.yml
    with:
      timeout: "10m"
      onlyNew: false

  build-release:
    needs: quality-gate
    runs-on: ubuntu-latest
    steps:
      - name: Build Release
        run: go build -o release/app
```

### Scheduled Analysis

```yaml
name: Weekly Code Analysis
on:
  schedule:
    - cron: '0 2 * * 1'  # Every Monday at 2 AM

jobs:
  comprehensive-analysis:
    uses: ./.github/workflows/go_staticAnalysis.yml
    with:
      timeout: "15m"
      onlyNew: false
      workspace: true
```

## Performance Optimization

### Timeout Configuration

```yaml
# For small projects
timeout: "1m"

# For medium projects
timeout: "5m"

# For large projects or workspaces
timeout: "15m"
```

### Selective Analysis

```yaml
# Analyze only new issues for faster feedback
onlyNew: true

# Full analysis for comprehensive checking
onlyNew: false
```

### Module-Specific Configuration

For workspaces, create module-specific configurations:

```
module1/.golangci.yml  # Configuration for module1
module2/.golangci.yml  # Configuration for module2
.golangci.yml          # Root configuration
```

## Troubleshooting

### Common Issues

**Analysis Timeout**

- Increase timeout value for large projects
- Consider using `onlyNew: true` for faster analysis
- Optimize golangci-lint configuration

**Module Detection Failures**

- Verify `go.mod` files exist in all modules
- Check Go workspace configuration (`go.work`)
- Ensure proper module structure

**Memory Issues**

- Reduce enabled linters for large codebases
- Use selective analysis with `onlyNew: true`
- Configure skip patterns for generated files

### Configuration Debugging

**Validate Configuration**

```bash
# Test golangci-lint configuration locally
golangci-lint run --config .golangci.yml
```

**Check Module Structure**

```bash
# List workspace modules
go work use -r .
```

**Performance Analysis**

```bash
# Run with verbose output
golangci-lint run -v --config .golangci.yml
```

## Best Practices

### Configuration Management

1. **Consistent Standards**: Use same linting rules across all modules
2. **Gradual Adoption**: Enable linters incrementally for existing projects
3. **Team Alignment**: Ensure team agrees on linting standards
4. **Documentation**: Document custom linting rules and exceptions

### Workflow Integration

1. **Early Feedback**: Run analysis on pull requests
2. **Quality Gates**: Block merges on analysis failures
3. **Regular Monitoring**: Schedule comprehensive analysis
4. **Performance Balance**: Balance thoroughness with execution time

### Error Handling

1. **Issue Tracking**: Track and resolve linting issues systematically
2. **Exclusion Rules**: Use exclusions judiciously for special cases
3. **Technical Debt**: Address linting debt in dedicated sprints
4. **Metrics**: Monitor code quality metrics over time

## Related Workflows

- **Go testing workflows**: For runtime testing after static analysis
- **go_build.yml**: For building Go applications
- **Security scanning workflows**: For vulnerability analysis

## Migration Guide

### From Manual Linting

1. Install golangci-lint locally for development
2. Create initial `.golangci.yml` configuration
3. Integrate workflow into CI/CD pipeline
4. Gradually increase linting strictness

### Workspace Migration

1. Convert multi-repository setup to Go workspace
2. Create `go.work` file with module references
3. Enable workspace mode in analysis workflow
4. Consolidate module-specific configurations