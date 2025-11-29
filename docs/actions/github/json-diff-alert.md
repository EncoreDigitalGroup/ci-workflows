# GitHub JSON Diff Alert Action

## Overview

The `actions/github/jsonDiffAlert` action compares JSON files between a source file and multiple destination files, then posts detailed difference reports as pull request
comments. This Docker-based action is useful for tracking configuration changes, API schema differences, and ensuring consistency across JSON files.

## Language/Tool Support

- **JSON**: All valid JSON file formats
- **GitHub**: Pull request integration and commenting
- **Docker**: Containerized execution for consistent behavior
- **Cross-File Comparison**: Multiple destination file support

## Features

- **Multi-File Comparison**: Compare one source against multiple destination files
- **Detailed Diff Reports**: Shows missing keys and new keys
- **PR Integration**: Automatically posts comments on pull requests
- **Flexible Path Resolution**: Supports relative and absolute paths
- **Validation Warnings**: Reports file validation issues
- **Root Directory Support**: Configurable base directory for path resolution

## Usage

```yaml
- name: Compare JSON Files
  uses: ./actions/github/jsonDiffAlert
  with:
    repository: ${{ github.repository }}
    pullRequestNumber: ${{ github.event.number }}
    token: ${{ secrets.GITHUB_TOKEN }}
    sourceFile: "config/production.json"
    destinationFiles: "config/staging.json,config/development.json"
    rootDirectory: ${{ github.workspace }}
```

## Inputs

| Input               | Type   | Required | Default                   | Description                                      |
|---------------------|--------|----------|---------------------------|--------------------------------------------------|
| `repository`        | string | ✅        | -                         | GitHub repository in format "owner/repo"         |
| `pullRequestNumber` | string | ✅        | -                         | Pull request number for commenting               |
| `token`             | string | ✅        | -                         | GitHub token with pull request write permissions |
| `sourceFile`        | string | ✅        | -                         | Path to source JSON file to compare from         |
| `destinationFiles`  | string | ✅        | -                         | Comma-separated list of destination JSON files   |
| `rootDirectory`     | string | ❌        | `${{ github.workspace }}` | Root directory for resolving relative file paths |

## Action Implementation

This action runs in a Docker container:

- **Image**: `ghcr.io/encoredigitalgroup/gh-action-json-diff-alert:latest`
- **Processing**: Go application for JSON comparison and GitHub API integration

## Comparison Logic

### Key Detection

The action extracts all JSON keys recursively:

- **Nested Objects**: Flattened using dot notation (`parent.child`)
- **Array Elements**: Indexed notation (`array[0].property`)
- **Deep Nesting**: Supports unlimited nesting depth

### Difference Types

1. **New in Source**: Keys present in source but missing from destinations
2. **Missing from Source**: Keys in destinations but not in source
3. **Validation Warnings**: Files that couldn't be processed

## Usage Examples

### API Schema Validation

```yaml
name: API Schema Check
on:
  pull_request:
    paths:
      - 'api/schema/**'

jobs:
  schema-diff:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - uses: actions/checkout@v4

      - name: Compare API Schemas
        uses: ./actions/github/jsonDiffAlert
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          token: ${{ secrets.GITHUB_TOKEN }}
          sourceFile: "api/schema/v2.json"
          destinationFiles: "api/schema/v1.json"
```

### Configuration Consistency

```yaml
name: Config Consistency Check
on:
  pull_request:
    paths:
      - 'config/**'

jobs:
  config-diff:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - uses: actions/checkout@v4

      - name: Check Configuration Consistency
        uses: ./actions/github/jsonDiffAlert
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          token: ${{ secrets.GITHUB_TOKEN }}
          sourceFile: "config/production.json"
          destinationFiles: "config/staging.json,config/development.json,config/test.json"
```

### Multi-Environment Validation

```yaml
name: Environment Config Validation
on:
  pull_request:
    paths:
      - 'environments/**'

jobs:
  env-validation:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - uses: actions/checkout@v4

      - name: Validate Environment Configurations
        uses: ./actions/github/jsonDiffAlert
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          token: ${{ secrets.GITHUB_TOKEN }}
          sourceFile: "environments/production/config.json"
          destinationFiles: "environments/staging/config.json,environments/dev/config.json"
          rootDirectory: ${{ github.workspace }}
```

### Package.json Comparison

```yaml
name: Package Consistency
on:
  pull_request:
    paths:
      - '**/package.json'

jobs:
  package-diff:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - uses: actions/checkout@v4

      - name: Compare Package Files
        uses: ./actions/github/jsonDiffAlert
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          token: ${{ secrets.GITHUB_TOKEN }}
          sourceFile: "frontend/package.json"
          destinationFiles: "backend/package.json,mobile/package.json"
```

## Comment Format

The action posts structured comments on pull requests:

### Example Comment Output

```markdown
## JSON Key Differences Report

### ⚠️ New keys in `config/production.json`:

**database.ssl.enabled:**
- Missing from: config/staging.json
- Missing from: config/development.json

**api.rateLimit.requests:**
- Missing from: config/development.json

### ⚠️ Keys present in destination files but missing from `config/production.json`:

**config/development.json:**
- `debug.verbose`
- `logging.level`

**config/staging.json:**
- `monitoring.enabled`

### ⚠️ File Processing Warnings:

- **config/invalid.json**: Invalid JSON file: unexpected token at line 15
```

## Integration Patterns

### Automated Configuration Auditing

```yaml
name: Config Audit
on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly on Monday

jobs:
  audit-configs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Create Audit PR
        id: create-pr
        uses: peter-evans/create-pull-request@v5
        with:
          title: "Weekly Configuration Audit"
          body: "Automated configuration consistency check"

      - name: Run Config Diff
        uses: ./actions/github/jsonDiffAlert
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ steps.create-pr.outputs.pull-request-number }}
          token: ${{ secrets.GITHUB_TOKEN }}
          sourceFile: "config/master.json"
          destinationFiles: "config/prod.json,config/stage.json,config/dev.json"
```

### Pre-Deployment Validation

```yaml
name: Pre-Deploy Validation
on:
  workflow_call:
    inputs:
      environment:
        required: true
        type: string

jobs:
  config-validation:
    runs-on: ubuntu-latest
    steps:
      - name: Validate Configuration
        uses: ./actions/github/jsonDiffAlert
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          token: ${{ secrets.GITHUB_TOKEN }}
          sourceFile: "config/template.json"
          destinationFiles: "config/${{ inputs.environment }}.json"
```

### Multi-Project Synchronization

```yaml
name: Cross-Project Config Sync
on:
  pull_request:
    paths:
      - 'shared-config/**'

jobs:
  config-sync-check:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - uses: actions/checkout@v4

      - name: Check Shared Configuration
        uses: ./actions/github/jsonDiffAlert
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          token: ${{ secrets.GITHUB_TOKEN }}
          sourceFile: "shared-config/base.json"
          destinationFiles: "projects/frontend/config.json,projects/backend/config.json,projects/mobile/config.json"
```

## Required Permissions

The GitHub token must have the following permissions:

```yaml
permissions:
  pull-requests: write  # Required for posting comments
  contents: read        # Required for accessing repository files
```

## File Path Resolution

### Relative Paths

Paths are resolved relative to `rootDirectory`:

```yaml
# With rootDirectory: "/workspace"
sourceFile: "config/app.json"        # Resolves to: /workspace/config/app.json
destinationFiles: "env/prod.json"    # Resolves to: /workspace/env/prod.json
```

## Use Cases

### API Versioning

Track changes between API versions:

```yaml
sourceFile: "api/v2/schema.json"
destinationFiles: "api/v1/schema.json,api/v3/schema.json"
```

### Environment Parity

Ensure consistency across environments:

```yaml
sourceFile: "config/production.json"
destinationFiles: "config/staging.json,config/development.json"
```

### Feature Flag Management

Monitor feature flag differences:

```yaml
sourceFile: "features/production.json"
destinationFiles: "features/staging.json,features/canary.json"
```

### Translation Files

Compare translation completeness:

```yaml
sourceFile: "i18n/en.json"
destinationFiles: "i18n/es.json,i18n/fr.json,i18n/de.json"
```

## Troubleshooting

### Common Issues

**No Comment Posted**

- Verify token has `pull-requests: write` permission
- Check pull request number is correct
- Ensure repository is accessible
- Is no comment on success enabled?

**File Not Found Errors**

- Verify file paths are correct relative to `rootDirectory`
- Check files exist in the repository
- Ensure proper checkout action runs first

**Invalid JSON Warnings**

- Validate JSON syntax in source/destination files
- Check for trailing commas or other JSON syntax errors
- Use JSON validators during development

### Error Resolution

**Permission Denied**

```yaml
# Ensure proper permissions
permissions:
  pull-requests: write
  contents: read
```

**Path Resolution Issues**

```yaml
# Use absolute paths or verify rootDirectory
rootDirectory: ${{ github.workspace }}
sourceFile: "./config/app.json"  # Relative to rootDirectory
```

## Best Practices

### File Organization

1. **Consistent Structure**: Use consistent JSON structure across environments
2. **Clear Naming**: Use descriptive file names for easy identification
3. **Logical Grouping**: Group related configurations together
4. **Documentation**: Comment on complex configuration differences

### Workflow Integration

1. **Targeted Triggers**: Use path filters to run only on relevant changes
2. **Early Validation**: Run comparisons early in the development process
3. **Review Process**: Include configuration reviews in PR process
4. **Notification**: Set up alerts for critical configuration differences

### Performance Optimization

1. **Selective Comparison**: Only compare relevant files
2. **Path Filtering**: Use workflow path filters
3. **Parallel Processing**: Run multiple comparisons in parallel
4. **Caching**: Cache results for repeated comparisons

## Migration Guide

### From Manual Configuration Review

1. Identify configuration files that need comparison
2. Set up automated comparison workflows
3. Train team on interpreting diff reports
4. Integrate into existing review process