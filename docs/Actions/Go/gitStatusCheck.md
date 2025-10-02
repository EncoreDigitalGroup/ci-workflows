# Go Git Status Check Action

## Overview

The `actions/go/gitStatusCheck` action determines whether Go-specific files have been modified by checking git status for Go source files (*.go), Go modules (go.mod), and
Go dependency files (go.sum). This action enables conditional workflow execution for Go projects, allowing build, test, and analysis jobs to run only when relevant Go
files have changed.

## Language/Tool Support

- **Go**: Go source files and module system
- **Git**: Git repository status checking
- **Go Modules**: go.mod and go.sum dependency management
- **Shell Scripts**: Bash-based implementation
- **GitHub Actions**: Integration with GitHub Actions workflow logic

## Features

- **Go-Specific File Monitoring**: Monitors Go source and module files
- **Branch-based Checking**: Configurable branch for comparison
- **Output Generation**: Provides boolean output for workflow decisions
- **Module-Aware**: Understands Go module structure
- **Performance Optimization**: Enables conditional execution to save resources

## Usage

```yaml
- name: Check Go File Changes
  id: go-check
  uses: ./actions/go/gitStatusCheck
  with:
    branch: main

- name: Run Go Build
  if: steps.go-check.outputs.shouldRun == 'true'
  run: go build ./...
```

## Inputs

| Input    | Type   | Required | Default | Description                       |
|----------|--------|----------|---------|-----------------------------------|
| `branch` | string | ❌        | `main`  | Branch to compare changes against |

## Outputs

| Output      | Type    | Description                                 |
|-------------|---------|---------------------------------------------|
| `shouldRun` | boolean | Whether Go files were modified (true/false) |

## Monitored File Types

The action checks for changes in:

### Go Source Files

- `*.go` - All Go source files
- Package files across all directories
- Test files (`*_test.go`)

### Go Module Files

- `go.mod` - Module definition and requirements
- `go.sum` - Cryptographic checksums for dependencies

### File Examples

```
main.go
cmd/server/main.go
pkg/api/handler.go
internal/config/config.go
main_test.go
go.mod
go.sum
```

## Usage Examples

### Go Build Pipeline

```yaml
name: Go Build Pipeline
on:
  pull_request:
    branches: [main]

jobs:
  check-go-changes:
    runs-on: ubuntu-latest
    outputs:
      go-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check for Go Changes
        id: check
        uses: ./actions/go/gitStatusCheck
        with:
          branch: main

  build:
    needs: check-go-changes
    if: needs.check-go-changes.outputs.go-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Build
        run: go build ./...

  test:
    needs: check-go-changes
    if: needs.check-go-changes.outputs.go-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Test
        run: go test ./...
```

### Dependency Change Detection

```yaml
name: Dependency Management
on:
  pull_request:
    paths:
      - 'go.mod'
      - 'go.sum'

jobs:
  check-deps:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Go Dependencies
        id: deps-check
        uses: ./actions/go/gitStatusCheck
        with:
          branch: ${{ github.base_ref }}

      - name: Security Audit
        if: steps.deps-check.outputs.shouldRun == 'true'
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

      - name: License Check
        if: steps.deps-check.outputs.shouldRun == 'true'
        run: |
          go mod download
          # Run license compliance check
```

### Multi-Module Project

```yaml
name: Multi-Module Go Project
on:
  push:
    branches: [main, develop]

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      go-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Go Changes
        id: check
        uses: ./actions/go/gitStatusCheck

  build-modules:
    needs: detect-changes
    if: needs.detect-changes.outputs.go-changed == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        module: [./cmd/api, ./cmd/worker, ./pkg/shared]
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Build Module
        working-directory: ${{ matrix.module }}
        run: go build
```

### Static Analysis Pipeline

```yaml
name: Go Static Analysis
on:
  pull_request:

jobs:
  check-go-changes:
    runs-on: ubuntu-latest
    outputs:
      should-analyze: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Go File Changes
        id: check
        uses: ./actions/go/gitStatusCheck
        with:
          branch: ${{ github.base_ref }}

  static-analysis:
    needs: check-go-changes
    if: needs.check-go-changes.outputs.should-analyze == 'true'
    uses: ./.github/workflows/go_staticAnalysis.yml
    with:
      branch: ${{ github.head_ref }}
      onlyNew: true
```

## Integration Patterns

### Performance-Optimized Pipeline

```yaml
name: Optimized Go Pipeline
on:
  pull_request:

jobs:
  quick-check:
    runs-on: ubuntu-latest
    outputs:
      go-changed: ${{ steps.go.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Go Changes
        id: go
        uses: ./actions/go/gitStatusCheck

  fast-feedback:
    needs: quick-check
    if: needs.quick-check.outputs.go-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Fast Compilation Check
        run: go build -o /dev/null ./...

  comprehensive-testing:
    needs: [quick-check, fast-feedback]
    if: needs.quick-check.outputs.go-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Comprehensive Test Suite
        run: |
          go test -race ./...
          go test -bench=. ./...
```

### Cross-Platform Conditional Builds

```yaml
name: Cross-Platform Go Build
on:
  push:
    tags:
      - 'v*'

jobs:
  check-go-changes:
    runs-on: ubuntu-latest
    outputs:
      build-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Go Changes Since Last Tag
        id: check
        uses: ./actions/go/gitStatusCheck
        with:
          branch: ${{ github.event.before }}

  cross-compile:
    needs: check-go-changes
    if: needs.check-go-changes.outputs.build-needed == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux, windows, darwin]
        goarch: [amd64, arm64]
    steps:
      - name: Build for ${{ matrix.goos }}-${{ matrix.goarch }}
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: go build -o dist/app-${{ matrix.goos }}-${{ matrix.goarch }}
```

### Workspace-Aware Checking

```yaml
name: Go Workspace Pipeline
on:
  pull_request:

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      workspace-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Go Workspace Changes
        id: check
        uses: ./actions/go/gitStatusCheck

  workspace-analysis:
    needs: detect-changes
    if: needs.detect-changes.outputs.workspace-changed == 'true'
    uses: ./.github/workflows/go_staticAnalysis.yml
    with:
      workspace: true
      timeout: "5m"
```

## Go-Specific Use Cases

### Module Dependency Updates

```yaml
name: Dependency Update Validation
on:
  pull_request:
    paths:
      - 'go.mod'
      - 'go.sum'

jobs:
  validate-deps:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Dependency Changes
        id: deps
        uses: ./actions/go/gitStatusCheck

      - name: Validate Dependencies
        if: steps.deps.outputs.shouldRun == 'true'
        run: |
          go mod verify
          go mod tidy
          git diff --exit-code go.mod go.sum
```

### Go Version Compatibility

```yaml
name: Go Version Compatibility
on:
  pull_request:

jobs:
  check-go-changes:
    runs-on: ubuntu-latest
    outputs:
      test-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Go Changes
        id: check
        uses: ./actions/go/gitStatusCheck

  compatibility-test:
    needs: check-go-changes
    if: needs.check-go-changes.outputs.test-needed == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.22', '1.23', '1.24']
    steps:
      - name: Test with Go ${{ matrix.go-version }}
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}

      - name: Run Tests
        run: go test ./...
```

### Go Generate Validation

```yaml
name: Generated Code Validation
on:
  pull_request:

jobs:
  check-source-changes:
    runs-on: ubuntu-latest
    outputs:
      generate-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Go Source Changes
        id: check
        uses: ./actions/go/gitStatusCheck

  validate-generated:
    needs: check-source-changes
    if: needs.check-source-changes.outputs.generate-needed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Run Go Generate
        run: go generate ./...

      - name: Check for Changes
        run: |
          if ! git diff --exit-code; then
            echo "Generated code is out of date"
            exit 1
          fi
```

## Branch Comparison Strategies

### Against Main Branch

```yaml
uses: ./actions/go/gitStatusCheck
with:
  branch: main
```

### Against PR Base

```yaml
uses: ./actions/go/gitStatusCheck
with:
  branch: ${{ github.base_ref }}
```

### Against Previous Tag

```yaml
uses: ./actions/go/gitStatusCheck
with:
  branch: ${{ github.event.before }}
```

## Troubleshooting

### Common Issues

**Always Returns False**

- Verify Go files exist in repository
- Check git history is available
- Ensure checkout action runs first

**Always Returns True**

- Check branch comparison logic
- Verify git diff commands work correctly
- Ensure proper working directory

**Missing Go Files**

- Verify file extensions are correct (*.go)
- Check for hidden or ignored files
- Ensure go.mod/go.sum exist if using modules

### Debugging Commands

```bash
# Check Go files manually
find . -name "*.go" -o -name "go.mod" -o -name "go.sum"

# Check git differences
git diff --name-only origin/main | grep -E '\.(go|mod|sum)$'

# Verify Go module structure
go list -m all
```

## Best Practices

### Workflow Design

1. **Early Execution**: Run checks at the beginning of workflows
2. **Output Sharing**: Share check results across multiple jobs
3. **Conditional Dependencies**: Use proper job dependencies
4. **Clear Documentation**: Document conditional logic

### Go Project Structure

1. **Module Organization**: Follow Go module best practices
2. **File Naming**: Use standard Go file naming conventions
3. **Directory Structure**: Organize code logically
4. **Dependency Management**: Keep go.mod and go.sum clean

## Migration Guide

### From Generic Path Filters

Replace generic path filters with Go-specific checks:

```yaml
# Before: Generic path filtering
on:
  pull_request:
    paths:
      - '**/*.go'
      - 'go.mod'
      - 'go.sum'

# After: Go-specific checking
on:
  pull_request:

jobs:
  check:
    steps:
      - uses: ./actions/go/gitStatusCheck
        id: go-check
```

### From Manual Git Commands

Replace manual git status checks:

```yaml
# Before: Manual commands
- name: Check Go Changes
  run: |
    if git diff --name-only origin/main | grep -E '\.(go|mod|sum)$'; then
      echo "go_changed=true" >> $GITHUB_OUTPUT
    fi

# After: Use action
- name: Check Go Changes
  id: check
  uses: ./actions/go/gitStatusCheck
```