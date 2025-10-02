# TypeScript Git Status Check Action

## Overview

The `actions/ts/gitStatusCheck` action determines whether TypeScript-specific files have been modified by checking git status for TypeScript source files (*.ts), package
configuration (package.json), and dependency lock files (package-lock.json). This action enables conditional workflow execution for TypeScript projects, optimizing build
and test processes.

## Language/Tool Support

- **TypeScript**: TypeScript source files and type definitions
- **Node.js**: package.json and package-lock.json dependency management
- **JavaScript**: Mixed TypeScript/JavaScript projects
- **Git**: Git repository status checking
- **Shell Scripts**: Bash-based implementation

## Features

- **TypeScript-Specific Monitoring**: Monitors .ts files and Node.js dependencies
- **Branch-based Checking**: Configurable branch for comparison
- **Output Generation**: Provides boolean output for workflow decisions
- **Modern Frontend Support**: Handles modern TypeScript/Node.js projects
- **Dependency-Aware**: Monitors package.json and lock file changes

## Usage

```yaml
- name: Check TypeScript File Changes
  id: ts-check
  uses: ./actions/ts/gitStatusCheck
  with:
    branch: main

- name: Run TypeScript Build
  if: steps.ts-check.outputs.shouldRun == 'true'
  run: npm run build
```

## Inputs

| Input    | Type   | Required | Default | Description                       |
|----------|--------|----------|---------|-----------------------------------|
| `branch` | string | ❌        | `main`  | Branch to compare changes against |

## Outputs

| Output      | Type    | Description                                                 |
|-------------|---------|-------------------------------------------------------------|
| `shouldRun` | boolean | Whether TypeScript-related files were modified (true/false) |

## Monitored File Types

The action checks for changes in:

### TypeScript Files

- `*.ts` - TypeScript source files
- Type definitions, interfaces, and implementations
- Test files (`*.test.ts`, `*.spec.ts`)

### Node.js Configuration

- `package.json` - Project dependencies and scripts
- `package-lock.json` - Locked dependency versions

### File Examples

```
src/index.ts
src/components/Button.ts
src/types/User.ts
src/utils/helpers.ts
tests/unit/Button.test.ts
package.json
package-lock.json
```

## Usage Examples

### React TypeScript Project

```yaml
name: React TypeScript CI
on:
  pull_request:
    branches: [main]

jobs:
  check-ts-changes:
    runs-on: ubuntu-latest
    outputs:
      ts-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check for TypeScript Changes
        id: check
        uses: ./actions/ts/gitStatusCheck
        with:
          branch: main

  build-and-test:
    needs: check-ts-changes
    if: needs.check-ts-changes.outputs.ts-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'

      - name: Install Dependencies
        run: npm ci

      - name: TypeScript Type Check
        run: npm run type-check

      - name: Build
        run: npm run build

      - name: Run Tests
        run: npm test
```

### Node.js TypeScript API

```yaml
name: Node.js TypeScript API
on:
  push:
    branches: [main, develop]

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      api-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check API Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  api-tests:
    needs: detect-changes
    if: needs.detect-changes.outputs.api-changed == 'true'
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
    steps:
      - name: Setup TypeScript Environment
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Dependencies
        run: npm ci

      - name: Run API Tests
        run: npm run test:api

      - name: Integration Tests
        run: npm run test:integration
```

### Monorepo TypeScript Project

```yaml
name: TypeScript Monorepo
on:
  pull_request:

jobs:
  check-ts-changes:
    runs-on: ubuntu-latest
    outputs:
      ts-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check TypeScript Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  build-packages:
    needs: check-ts-changes
    if: needs.check-ts-changes.outputs.ts-changed == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        package: [frontend, backend, shared]
    steps:
      - name: Build ${{ matrix.package }}
        working-directory: packages/${{ matrix.package }}
        run: |
          npm ci
          npm run build
          npm test
```

### TypeScript Library Development

```yaml
name: TypeScript Library
on:
  pull_request:

jobs:
  check-library-changes:
    runs-on: ubuntu-latest
    outputs:
      library-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Library Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  validate-library:
    needs: check-library-changes
    if: needs.check-library-changes.outputs.library-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Dependencies
        run: npm ci

      - name: Type Check
        run: npx tsc --noEmit

      - name: Build Library
        run: npm run build

      - name: Test Library
        run: npm test

      - name: Generate Documentation
        run: npm run docs
```

## Integration Patterns

### Performance-Optimized Pipeline

```yaml
name: Optimized TypeScript Pipeline
on:
  pull_request:

jobs:
  quick-check:
    runs-on: ubuntu-latest
    outputs:
      ts-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check TypeScript Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  fast-validation:
    needs: quick-check
    if: needs.quick-check.outputs.ts-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Quick Type Check
        run: npx tsc --noEmit --incremental

  comprehensive-testing:
    needs: [quick-check, fast-validation]
    if: needs.quick-check.outputs.ts-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Full Test Suite
        run: |
          npm ci
          npm run test:unit
          npm run test:e2e
```

### Multi-Environment Deployment

```yaml
name: TypeScript Deployment Pipeline
on:
  push:
    branches: [main, staging, develop]

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      deploy-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Deployment Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  build-and-deploy:
    needs: detect-changes
    if: needs.detect-changes.outputs.deploy-needed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Build Application
        run: |
          npm ci
          npm run build

      - name: Deploy to Environment
        run: |
          if [ "${{ github.ref }}" == "refs/heads/main" ]; then
            echo "Deploying to production"
          elif [ "${{ github.ref }}" == "refs/heads/staging" ]; then
            echo "Deploying to staging"
          fi
```

### Frontend Framework Integration

```yaml
name: Frontend Framework Pipeline
on:
  pull_request:

jobs:
  detect-frontend-changes:
    runs-on: ubuntu-latest
    outputs:
      frontend-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Frontend Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  react-build:
    needs: detect-frontend-changes
    if: needs.detect-frontend-changes.outputs.frontend-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Build React App
        run: |
          npm ci
          npm run build
          npm run test

  angular-build:
    needs: detect-frontend-changes
    if: needs.detect-frontend-changes.outputs.frontend-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Build Angular App
        run: |
          npm ci
          npm run build
          npm run test:ci

  vue-build:
    needs: detect-frontend-changes
    if: needs.detect-frontend-changes.outputs.frontend-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Build Vue App
        run: |
          npm ci
          npm run build
          npm run test:unit
```

## TypeScript-Specific Use Cases

### Type Definition Validation

```yaml
name: Type Definition Validation
on:
  pull_request:
    paths:
      - 'src/types/**'
      - '**/*.d.ts'

jobs:
  validate-types:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Type Changes
        id: types
        uses: ./actions/ts/gitStatusCheck

      - name: Validate Type Definitions
        if: steps.types.outputs.shouldRun == 'true'
        run: |
          npm ci
          npx tsc --noEmit
          npm run type-coverage
```

### Dependency Security Scanning

```yaml
name: TypeScript Security Audit
on:
  pull_request:
    paths:
      - 'package.json'
      - 'package-lock.json'

jobs:
  security-audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Dependency Changes
        id: deps
        uses: ./actions/ts/gitStatusCheck

      - name: Security Audit
        if: steps.deps.outputs.shouldRun == 'true'
        run: |
          npm audit --audit-level moderate
          npx audit-ci
```

### Multi-Version Node.js Testing

```yaml
name: Node.js Version Compatibility
on:
  pull_request:

jobs:
  check-ts-changes:
    runs-on: ubuntu-latest
    outputs:
      test-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check TypeScript Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  compatibility-test:
    needs: check-ts-changes
    if: needs.check-ts-changes.outputs.test-needed == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: ['18', '20', '22']
    steps:
      - name: Setup Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v4
        with:
          node-version: ${{ matrix.node-version }}

      - name: Test Compatibility
        run: |
          npm ci
          npm test
```

### Code Generation Validation

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

      - name: Check Source Changes
        id: check
        uses: ./actions/ts/gitStatusCheck

  validate-generated:
    needs: check-source-changes
    if: needs.check-source-changes.outputs.generate-needed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Generate TypeScript Code
        run: |
          npm run generate:types
          npm run generate:api-client

      - name: Check for Changes
        run: |
          if ! git diff --exit-code; then
            echo "Generated code is out of date"
            exit 1
          fi
```

## File Detection Patterns

### Standard TypeScript Patterns

- Source files: `src/**/*.ts`
- Type definitions: `src/types/**/*.ts`, `**/*.d.ts`
- Tests: `src/**/*.test.ts`, `src/**/*.spec.ts`
- Configuration: `tsconfig.json`, `tsconfig.*.json`

### Framework-Specific Patterns

#### React

- Components: `src/components/**/*.tsx`
- Hooks: `src/hooks/**/*.ts`
- Utils: `src/utils/**/*.ts`

#### Angular

- Services: `src/app/**/*.service.ts`
- Components: `src/app/**/*.component.ts`
- Modules: `src/app/**/*.module.ts`

#### Vue

- Composables: `src/composables/**/*.ts`
- Stores: `src/stores/**/*.ts`
- Types: `src/types/**/*.ts`

#### Node.js

- Routes: `src/routes/**/*.ts`
- Controllers: `src/controllers/**/*.ts`
- Models: `src/models/**/*.ts`

## Performance Benefits

### Development Efficiency

- **Faster Builds**: Skip builds when no TypeScript files changed
- **Resource Optimization**: Avoid unnecessary dependency installation

### Workflow Optimization

- **Smart Caching**: Better dependency and build cache utilization
- **Parallel Processing**: Enable conditional parallel job execution
- **Resource Allocation**: Allocate build resources based on actual needs
- **Developer Experience**: Faster PR feedback and iteration

## Troubleshooting

### Common Issues

**Always Returns False**

- Verify TypeScript files exist in repository
- Check that package.json exists
- Ensure checkout action runs before this action

**Missing Node.js Dependencies**

- Verify package.json and package-lock.json are committed
- Check Node.js project structure
- Ensure proper TypeScript configuration

**Type Checking Issues**

- Verify tsconfig.json exists and is valid
- Check TypeScript dependency in package.json
- Ensure proper type definitions are installed

### Debugging Commands

```bash
# Check TypeScript files
find . -name "*.ts" -o -name "package.json" -o -name "package-lock.json"

# Check TypeScript configuration
cat tsconfig.json

# Check git differences
git diff --name-only origin/main | grep -E '\.(ts|json)$'

# Validate TypeScript setup
npx tsc --version
npx tsc --noEmit
```

## Best Practices

### Project Structure

1. **Standard Layout**: Follow TypeScript/Node.js conventions
2. **Type Organization**: Organize types in dedicated directories
3. **Configuration**: Use proper tsconfig.json configuration
4. **Dependencies**: Keep package-lock.json committed

### Workflow Design

1. **Early Checking**: Run checks early in the pipeline
2. **Conditional Execution**: Use proper conditional logic

### Performance Optimization

1. **Incremental Builds**: Use TypeScript incremental compilation
2. **Cache Strategy**: Implement effective caching
3. **Parallel Jobs**: Run checks in parallel with other operations

## Migration Guide

### From Path-Based Triggers

Replace path-based workflow triggers:

```yaml
# Before: Path-based filtering
on:
  pull_request:
    paths:
      - '**/*.ts'
      - 'package.json'
      - 'package-lock.json'

# After: Dynamic checking
on:
  pull_request:

jobs:
  check:
    steps:
      - uses: ./actions/ts/gitStatusCheck
        id: ts-check
```

### From Manual File Detection

Replace manual file checking:

```yaml
# Before: Manual commands
- name: Check TypeScript Changes
  run: |
    if git diff --name-only origin/main | grep -E '\.(ts|json)$'; then
      echo "ts_changed=true" >> $GITHUB_OUTPUT
    fi

# After: Use action
- name: Check TypeScript Changes
  id: check
  uses: ./actions/ts/gitStatusCheck
```

### Framework-Specific Integration

#### React Projects

```yaml
- uses: ./actions/ts/gitStatusCheck
  id: react-check

- name: Build React App
  if: steps.react-check.outputs.shouldRun == 'true'
  run: npm run build
```

#### Node.js APIs

```yaml
- uses: ./actions/ts/gitStatusCheck
  id: api-check

- name: API Tests
  if: steps.api-check.outputs.shouldRun == 'true'
  run: npm run test:api
```

#### Library Projects

```yaml
- uses: ./actions/ts/gitStatusCheck
  id: lib-check

- name: Build Library
  if: steps.lib-check.outputs.shouldRun == 'true'
  run: |
    npm run build
    npm run test
    npm run docs
```