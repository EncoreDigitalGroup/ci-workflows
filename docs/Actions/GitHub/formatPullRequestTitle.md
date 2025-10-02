# GitHub Format Pull Request Title Action

## Overview

The `actions/github/formatPullRequestTitle` action automatically formats pull request titles based on branch naming conventions. It can enhance titles with proper
capitalization, spacing, and custom formatting rules, making PR titles more readable and consistent across the project.

## Language/Tool Support

- **GitHub**: Pull request management and API integration
- **Branch Naming**: Automatic parsing of branch name conventions
- **Custom Formatting**: User-defined formatting rules support

## Features

- **Automatic Title Formatting**: Converts branch names to readable PR titles
- **Custom Formatting Rules**: Support for project-specific formatting preferences
- **Intelligent Parsing**: Smart handling of common branch naming patterns
- **Caching**: Optimized caching for build artifacts and dependencies

## Usage

```yaml
- name: Format Pull Request Title
  uses: ./actions/github/formatPullRequestTitle
  with:
    repository: ${{ github.repository }}
    pullRequestNumber: ${{ github.event.number }}
    branch: ${{ github.head_ref }}
    token: ${{ secrets.GITHUB_TOKEN }}
    customFormatting: "api:API,ui:User Interface,db:Database"
```

## Inputs

| Input               | Type    | Required | Default | Description                                          |
|---------------------|---------|----------|---------|------------------------------------------------------|
| `repository`        | string  | ✅        | -       | GitHub repository in format "owner/repo"             |
| `pullRequestNumber` | string  | ✅        | -       | Pull request number to update                        |
| `branch`            | string  | ✅        | -       | Branch name to parse for title formatting            |
| `token`             | string  | ✅        | -       | GitHub token with pull request write permissions     |
| `customFormatting`  | string  | ❌        | `""`    | Custom word formatting rules (comma-separated pairs) |

## Branch Name Patterns

The action recognizes common branch naming conventions:

### Supported Patterns

- **Feature branches**: `feature/user-authentication` → "Feature: User Authentication"
- **Bug fixes**: `fix/login-issue` → "Fix: Login Issue"
- **Hotfixes**: `hotfix/security-patch` → "Hotfix: Security Patch"
- **Improvements**: `improvement/api-performance` → "Improvement: API Performance"
- **Chores**: `chore/update-dependencies` → "Chore: Update Dependencies"

### Custom Patterns

Define custom formatting with the `customFormatting` input:

```yaml
customFormatting: "feat:Feature,fix:Bug Fix,docs:Documentation,test:Testing"
```

## Usage Examples

### Basic PR Title Formatting

```yaml
name: Format PR Title
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  format-title:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - name: Format Pull Request Title
        uses: ./actions/github/formatPullRequestTitle
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          branch: ${{ github.head_ref }}
          token: ${{ secrets.GITHUB_TOKEN }}
```

### Custom Formatting Rules

```yaml
name: Custom PR Formatting
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  format-title:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - name: Format with Custom Rules
        uses: ./actions/github/formatPullRequestTitle
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          branch: ${{ github.head_ref }}
          token: ${{ secrets.GITHUB_TOKEN }}
          customFormatting: "api:API Integration,ui:User Interface,auth:Authentication,db:Database Operations"
```

### Multi-Team Configuration

```yaml
name: Team-Specific Formatting
on:
  pull_request:
    types: [opened]

jobs:
  format-frontend:
    if: startsWith(github.head_ref, 'frontend/')
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - name: Format Frontend PR
        uses: ./actions/github/formatPullRequestTitle
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          branch: ${{ github.head_ref }}
          token: ${{ secrets.GITHUB_TOKEN }}
          customFormatting: "ui:User Interface,ux:User Experience,css:Styling"

  format-backend:
    if: startsWith(github.head_ref, 'backend/')
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - name: Format Backend PR
        uses: ./actions/github/formatPullRequestTitle
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          branch: ${{ github.head_ref }}
          token: ${{ secrets.GITHUB_TOKEN }}
          customFormatting: "api:API Development,db:Database,auth:Authentication"
          useGo: true
```

## Formatting Examples

### Branch to Title Conversion

| Branch Name                     | Formatted Title                |
|---------------------------------|--------------------------------|
| `feature/user-login`            | Feature: User Login            |
| `fix/broken-api`                | Fix: Broken API                |
| `hotfix/security-vulnerability` | Hotfix: Security Vulnerability |
| `chore/update-deps`             | Chore: Update Deps             |
| `improvement/load-time`         | Improvement: Load Time         |
| `docs/api-documentation`        | Docs: API Documentation        |

### With Custom Formatting

```yaml
customFormatting: "api:API,ui:User Interface,db:Database"
```

| Branch Name                  | Formatted Title                   |
|------------------------------|-----------------------------------|
| `feature/api-integration`    | Feature: API Integration          |
| `fix/ui-layout`              | Fix: User Interface Layout        |
| `improvement/db-performance` | Improvement: Database Performance |

## Integration Patterns

### Automated PR Workflow

```yaml
name: PR Automation
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  format-and-validate:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - name: Format PR Title
        uses: ./actions/github/formatPullRequestTitle
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          branch: ${{ github.head_ref }}
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: Validate PR
        run: echo "Running additional validation..."
```

### Branch Protection Integration

```yaml
name: PR Quality Gate
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  format-title:
    runs-on: ubuntu-latest
    permissions:
      pull-requests: write
    steps:
      - name: Format Title
        uses: ./actions/github/formatPullRequestTitle
        with:
          repository: ${{ github.repository }}
          pullRequestNumber: ${{ github.event.number }}
          branch: ${{ github.head_ref }}
          token: ${{ secrets.GITHUB_TOKEN }}

  required-checks:
    needs: format-title
    runs-on: ubuntu-latest
    steps:
      - name: Run Tests
        run: echo "Running required tests..."
```

## Required Permissions

The GitHub token must have the following permissions:

```yaml
permissions:
  pull-requests: write  # Required for updating PR titles
  contents: read        # Required for accessing repository
```

## Custom Formatting Rules

Define custom formatting with key-value pairs:

### Format

```
"key1:value1,key2:value2,key3:value3"
```

### Examples

```yaml
# Technical abbreviations
customFormatting: "api:API,ui:UI,db:Database,auth:Authentication"

# Business domains
customFormatting: "crm:Customer Relations,hr:Human Resources,fin:Finance"

# Team prefixes
customFormatting: "fe:Frontend,be:Backend,qa:Quality Assurance,ops:Operations"
```

## Troubleshooting

### Common Issues

**Title Not Updated**

- Verify token has `pull-requests: write` permission
- Check pull request number is correct
- Ensure repository is accessible

**Formatting Not Applied**

- Verify branch name follows expected patterns
- Check custom formatting syntax
- Test with simplified formatting rules

### Error Resolution

**Permission Denied**

```yaml
# Ensure proper permissions
permissions:
  pull-requests: write
  contents: read
```

### Branch Naming Conventions

Establish consistent branch naming:

```yaml
# Recommended patterns
feature/description      # New features
fix/description         # Bug fixes
hotfix/description      # Critical fixes
chore/description       # Maintenance tasks
docs/description        # Documentation
test/description        # Testing improvements
```

## Best Practices

### Branch Naming

1. **Consistent Prefixes**: Use standard prefixes (feature, fix, hotfix)
2. **Descriptive Names**: Use clear, descriptive branch names
3. **Hyphen Separation**: Use hyphens instead of underscores
4. **Lowercase**: Use lowercase for consistency

### Custom Formatting

1. **Team Standards**: Align formatting rules with team conventions
2. **Abbreviation Consistency**: Use consistent abbreviations
3. **Documentation**: Document custom formatting rules
4. **Regular Review**: Review and update formatting rules periodically

### Workflow Integration

1. **Early Execution**: Run formatting early in PR workflow
2. **Conditional Logic**: Use conditions for team-specific formatting