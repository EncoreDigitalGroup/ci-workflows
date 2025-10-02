# GitHub Git Status Check Action

## Overview

The `actions/github/gitStatusCheck` action determines whether specific files have been modified by checking git status for YAML files (*.yml and *.yaml). This action is
useful for conditional workflow execution, allowing subsequent jobs to run only when relevant files have changed.

## Language/Tool Support

- **Git**: Git repository status checking
- **YAML**: Focuses on YAML file changes (*.yml, *.yaml)
- **Shell Scripts**: Bash-based implementation
- **GitHub Actions**: Integration with GitHub Actions workflow logic

## Features

- **Selective File Monitoring**: Monitors only YAML files for changes
- **Branch-based Checking**: Configurable branch for comparison
- **Output Generation**: Provides boolean output for workflow decisions
- **Lightweight Execution**: Fast execution with minimal resource usage
- **Conditional Workflows**: Enables conditional job execution

## Usage

```yaml
- name: Check YAML File Changes
  id: yaml-check
  uses: ./actions/github/gitStatusCheck
  with:
    branch: main

- name: Run YAML Processing
  if: steps.yaml-check.outputs.shouldRun == 'true'
  run: echo "YAML files have changed, processing..."
```

## Inputs

| Input    | Type   | Required | Default | Description                       |
|----------|--------|----------|---------|-----------------------------------|
| `branch` | string | ❌        | `main`  | Branch to compare changes against |

## Outputs

| Output      | Type    | Description                                   |
|-------------|---------|-----------------------------------------------|
| `shouldRun` | boolean | Whether YAML files were modified (true/false) |

## Action Implementation

- **Script**: Bash entrypoint script (`entrypoint.sh`)
- **File Types**: Monitors `*.yml` and `*.yaml` files
- **Git Commands**: Uses git status and diff commands
- **Environment**: Uses `BRANCH_NAME` environment variable

## Usage Examples

### Conditional Workflow Execution

```yaml
name: YAML File Processing
on:
  pull_request:
    branches: [main]

jobs:
  check-yaml-changes:
    runs-on: ubuntu-latest
    outputs:
      yaml-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check for YAML Changes
        id: check
        uses: ./actions/github/gitStatusCheck
        with:
          branch: main

  process-yaml:
    needs: check-yaml-changes
    if: needs.check-yaml-changes.outputs.yaml-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Process YAML Files
        run: echo "Processing YAML files..."
```

### Workflow Configuration Validation

```yaml
name: Workflow Validation
on:
  pull_request:
    paths:
      - '.github/workflows/**'

jobs:
  check-workflow-changes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Workflow File Changes
        id: workflow-check
        uses: ./actions/github/gitStatusCheck
        with:
          branch: ${{ github.base_ref }}

      - name: Validate Workflows
        if: steps.workflow-check.outputs.shouldRun == 'true'
        run: |
          echo "Workflow files changed, running validation..."
          # Add workflow validation logic here
```

### Docker Compose Updates

```yaml
name: Docker Compose Validation
on:
  push:
    branches: [develop, staging]

jobs:
  check-compose-changes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Docker Compose Changes
        id: compose-check
        uses: ./actions/github/gitStatusCheck
        with:
          branch: main

      - name: Validate Docker Compose
        if: steps.compose-check.outputs.shouldRun == 'true'
        run: |
          echo "Docker Compose files changed"
          docker-compose config --quiet
```

### Kubernetes Manifest Updates

```yaml
name: Kubernetes Validation
on:
  pull_request:
    paths:
      - 'k8s/**'
      - 'kubernetes/**'

jobs:
  check-k8s-changes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Kubernetes Manifest Changes
        id: k8s-check
        uses: ./actions/github/gitStatusCheck
        with:
          branch: ${{ github.base_ref }}

      - name: Validate Kubernetes Manifests
        if: steps.k8s-check.outputs.shouldRun == 'true'
        run: |
          echo "Kubernetes manifests changed"
          kubectl apply --dry-run=client -f k8s/
```

## Integration Patterns

### Multi-Stage Conditional Pipeline

```yaml
name: Conditional Pipeline
on:
  pull_request:

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      yaml-changed: ${{ steps.yaml.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check YAML Changes
        id: yaml
        uses: ./actions/github/gitStatusCheck

  validate-yaml:
    needs: detect-changes
    if: needs.detect-changes.outputs.yaml-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Validate YAML Syntax
        run: yamllint .

  deploy-config:
    needs: [detect-changes, validate-yaml]
    if: needs.detect-changes.outputs.yaml-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy Configuration
        run: echo "Deploying updated configuration..."
```

### Branch-Specific Validation

```yaml
name: Branch-Specific Checks
on:
  push:
    branches: [main, develop, release/*]

jobs:
  check-changes:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Changes Against Main
        id: main-check
        uses: ./actions/github/gitStatusCheck
        with:
          branch: main

      - name: Production Deployment
        if: github.ref == 'refs/heads/main' && steps.main-check.outputs.shouldRun == 'true'
        run: echo "Deploying to production..."

      - name: Staging Deployment
        if: github.ref == 'refs/heads/develop' && steps.main-check.outputs.shouldRun == 'true'
        run: echo "Deploying to staging..."
```

### Performance Optimization

```yaml
name: Optimized Pipeline
on:
  pull_request:

jobs:
  quick-checks:
    runs-on: ubuntu-latest
    outputs:
      should-run-expensive: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check for Configuration Changes
        id: check
        uses: ./actions/github/gitStatusCheck

  expensive-validation:
    needs: quick-checks
    if: needs.quick-checks.outputs.should-run-expensive == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Run Expensive Validation
        run: |
          echo "Running comprehensive validation..."
          # Resource-intensive validation only when needed
```

## File Detection Logic

The action checks for changes in:

### Included File Types

- `*.yml` - YAML files
- `*.yaml` - YAML files (alternative extension)

### Common Use Cases

- **GitHub Workflows**: `.github/workflows/*.yml`
- **Docker Compose**: `docker-compose.yml`, `docker-compose.yaml`
- **Kubernetes**: `*.yaml`, `*.yml` in k8s directories
- **Configuration**: `config.yml`, `application.yaml`
- **CI/CD**: Various pipeline configuration files

## Branch Comparison

The action compares the current state against the specified branch:

### Default Behavior

```yaml
# Compares against main branch
uses: ./actions/github/gitStatusCheck
```

### Custom Branch

```yaml
# Compares against develop branch
uses: ./actions/github/gitStatusCheck
with:
  branch: develop
```

### Dynamic Branch

```yaml
# Compares against PR base branch
uses: ./actions/github/gitStatusCheck
with:
  branch: ${{ github.base_ref }}
```

## Environment Configuration

The action uses environment variables:

- `BRANCH_NAME`: Branch to compare against (from input)

## Use Cases

### Configuration Management

Monitor configuration file changes:

```yaml
name: Config Change Detection
on:
  pull_request:
    paths:
      - 'config/**'

jobs:
  config-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Config Changes
        id: config
        uses: ./actions/github/gitStatusCheck

      - name: Validate Configuration
        if: steps.config.outputs.shouldRun == 'true'
        run: echo "Configuration changed, validating..."
```

### Infrastructure as Code

Monitor infrastructure changes:

```yaml
name: Infrastructure Validation
on:
  pull_request:
    paths:
      - 'infrastructure/**'
      - 'terraform/**'

jobs:
  infra-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Infrastructure Changes
        id: infra
        uses: ./actions/github/gitStatusCheck

      - name: Plan Infrastructure
        if: steps.infra.outputs.shouldRun == 'true'
        run: terraform plan
```

### Deployment Automation

Conditional deployment based on changes:

```yaml
name: Smart Deployment
on:
  push:
    branches: [main]

jobs:
  check-deployment-needed:
    runs-on: ubuntu-latest
    outputs:
      deploy: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Deployment Files
        id: check
        uses: ./actions/github/gitStatusCheck

  deploy:
    needs: check-deployment-needed
    if: needs.check-deployment-needed.outputs.deploy == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy Application
        run: echo "Deploying application..."
```

## Troubleshooting

### Common Issues

**Action Always Returns False**

- Check that YAML files actually exist in the repository
- Ensure the branch comparison is correct

**Action Always Returns True**

- Check if git history is available
- Verify branch exists and is accessible
- Ensure proper git configuration

**Permission Issues**

- Ensure repository checkout is performed
- Verify git commands can access repository
- Check working directory is correct

### Debugging

**Check Git Status Manually**

```bash
# Debug git status locally
git status --porcelain | grep -E '\.(yml|yaml)$'
git diff --name-only origin/main | grep -E '\.(yml|yaml)$'
```

**Verify Branch Existence**

```bash
# Check if branch exists
git branch -r | grep origin/main
```

**Test File Detection**

```bash
# Test file pattern matching
find . -name "*.yml" -o -name "*.yaml"
```

## Performance Considerations

### Execution Time

- **Fast**: Minimal git operations
- **Lightweight**: Simple shell script execution
- **Efficient**: Early termination when files found

### Resource Usage

- **Low Memory**: Minimal memory footprint
- **Quick I/O**: Limited file system operations
- **Network**: No external network calls

### Optimization Tips

1. **Place Early**: Run early in workflow for maximum benefit
2. **Combine Checks**: Use outputs in multiple subsequent jobs

## Best Practices

### Workflow Design

1. **Early Execution**: Run checks early in the workflow
2. **Output Reuse**: Share outputs between multiple jobs
3. **Clear Naming**: Use descriptive step IDs and names
4. **Documentation**: Document conditional logic clearly

### Branch Strategy

1. **Consistent Branches**: Use consistent branch names for comparison
2. **Dynamic References**: Use dynamic branch references in PRs
3. **Default Branches**: Align with repository default branch

## Migration Guide

### From Path Filters

Replace workflow path filters with dynamic checks:

```yaml
# Before: Static path filtering
on:
  pull_request:
    paths:
      - '*.yml'
      - '*.yaml'

# After: Dynamic checking
on:
  pull_request:

jobs:
  check:
    steps:
      - uses: ./actions/github/gitStatusCheck
        id: yaml-check
      - if: steps.yaml-check.outputs.shouldRun == 'true'
        run: echo "Process YAML files"
```

### From Manual Git Commands

Replace manual git status checks:

```yaml
# Before: Manual git commands
- name: Check Changes
  run: |
    if git diff --name-only origin/main | grep -E '\.(yml|yaml)$'; then
      echo "yaml_changed=true" >> $GITHUB_OUTPUT
    fi

# After: Use action
- name: Check Changes
  id: check
  uses: ./actions/github/gitStatusCheck
```