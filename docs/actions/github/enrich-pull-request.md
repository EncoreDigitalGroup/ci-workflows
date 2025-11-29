# GitHub Enrich Pull Request Action

## Overview

The `actions/github/enrichPullRequest` action automatically enriches pull request titles and descriptions with information from project management systems. It supports
multiple strategies including branch name parsing and Jira integration, providing enhanced PR context and consistency across development workflows.

## Language/Tool Support

- **GitHub**: Pull request management and API integration
- **Jira**: Issue tracking and project management integration
- **Branch Naming**: Automatic parsing of branch name conventions

## Features

- **Multiple Enrichment Strategies**: Support for branch-name and Jira strategies
- **Automatic PR Title Formatting**: Converts branch names and issue keys to readable titles
- **Jira Integration**: Syncs Jira issue titles and descriptions to pull requests
- **Custom Formatting Rules**: User-defined formatting preferences
- **Label Management**: Automatic label creation and assignment for Jira sync tracking
- **Parent Issue Support**: Includes parent issue prefixes for hierarchical issues

## Usage

```yaml
-   name: Enrich Pull Request
    uses: EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest@v3
    with:
        repository: ${{ github.repository }}
        pullRequestNumber: ${{ github.event.number }}
        branch: ${{ github.head_ref }}
        token: ${{ secrets.GITHUB_TOKEN }}
        strategy: "branch-name"
        customFormatting: "api:API,ui:User Interface"
```

## Inputs

| Input                       | Type    | Required | Default                | Description                                            |
|-----------------------------|---------|----------|------------------------|--------------------------------------------------------|
| `repository`                | string  | ✅        | -                      | GitHub repository in format "owner/repo"               |
| `pullRequestNumber`         | string  | ✅        | -                      | Pull request number to update                          |
| `branch`                    | string  | ✅        | -                      | Branch name to parse for enrichment                    |
| `token`                     | string  | ✅        | -                      | GitHub token with pull request write permissions       |
| `strategy`                  | string  | ❌        | `"branch-name"`        | Enrichment strategy: "branch-name" or "jira"           |
| `customFormatting`          | string  | ❌        | `""`                   | Custom word formatting rules (comma-separated pairs)   |
| `jiraURL`                   | string  | ❌        | `""`                   | URL to your Jira instance (required for jira strategy) |
| `jiraEmail`                 | string  | ❌        | `""`                   | Jira authentication email (required for jira strategy) |
| `jiraToken`                 | string  | ❌        | `""`                   | Jira authentication token (required for jira strategy) |
| `jiraEnableSyncLabel`       | boolean | ❌        | `true`                 | Create and assign sync completion label                |
| `jiraEnableSyncDescription` | boolean | ❌        | `true`                 | Sync Jira description to PR description                |
| `jiraSyncLabelName`         | string  | ❌        | `"jira-sync-complete"` | Name of the sync completion label                      |

## Action Implementation

This action runs in a Docker container:

- **Image**: `ghcr.io/encoredigitalgroup/gh-action-enrich-pull-request:latest`
- **Environment**: All inputs are passed as environment variables
- **Execution**: Containerized Go application with multiple strategy drivers

## Enrichment Strategies

### Branch Name Strategy

Extracts issue keys and descriptions from branch names using regex patterns:

**Supported Branch Patterns:**

- `(epic|feature|bugfix|hotfix)/[A-Z]+-[0-9]+-summary`
- `[A-Z]+-[0-9]+-summary`

**Examples:**

- `feature/PROJ-123-user-authentication` → "[PROJ-123] User Authentication"
- `bugfix/ISSUE-456-login-fix` → "[ISSUE-456] Login Fix"
- `TASK-789-api-improvements` → "[TASK-789] API Improvements"

### Jira Strategy

Integrates with Jira to fetch issue information and enrich PRs:

**Features:**

- Fetches issue title from Jira
- Optionally syncs issue description to PR (enabled by default)
- Adds parent issue prefixes for subtasks
- Creates sync completion labels (enabled by default)
- Prevents duplicate syncing

## Usage Examples

### Basic Branch Name Enrichment

```yaml
name: Enrich PR with Branch Info
on:
    pull_request:
        types: [ opened, synchronize ]

jobs:
    enrich-pr:
        runs-on: ubuntu-latest
        permissions:
            pull-requests: write
        steps:
            -   name: Enrich Pull Request
                uses: EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest@v3
                with:
                    repository: ${{ github.repository }}
                    pullRequestNumber: ${{ github.event.number }}
                    branch: ${{ github.head_ref }}
                    token: ${{ secrets.GITHUB_TOKEN }}
                    strategy: "branch-name"
```

### Jira Integration

```yaml
name: Enrich PR with Jira
on:
    pull_request:
        types: [ opened, synchronize ]

jobs:
    enrich-pr:
        runs-on: ubuntu-latest
        permissions:
            pull-requests: write
        steps:
            -   name: Enrich with Jira Info
                uses: EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest@v3
                with:
                    repository: ${{ github.repository }}
                    pullRequestNumber: ${{ github.event.number }}
                    branch: ${{ github.head_ref }}
                    token: ${{ secrets.GITHUB_TOKEN }}
                    strategy: "jira"
                    jiraURL: ${{ vars.JIRA_URL }}
                    jiraEmail: ${{ vars.JIRA_EMAIL }}
                    jiraToken: ${{ secrets.JIRA_TOKEN }}
                    jiraEnableSyncLabel: true
                    jiraEnableSyncDescription: true
```

### Custom Formatting Rules

```yaml
name: Enrich with Custom Formatting
on:
    pull_request:
        types: [ opened, synchronize ]

jobs:
    enrich-pr:
        runs-on: ubuntu-latest
        permissions:
            pull-requests: write
        steps:
            -   name: Enrich with Custom Rules
                uses: EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest@v3
                with:
                    repository: ${{ github.repository }}
                    pullRequestNumber: ${{ github.event.number }}
                    branch: ${{ github.head_ref }}
                    token: ${{ secrets.GITHUB_TOKEN }}
                    strategy: "branch-name"
                    customFormatting: "api:API Integration,ui:User Interface,db:Database Operations,auth:Authentication"
```

## Integration Patterns

### Complete PR Workflow

```yaml
name: PR Quality Gate
on:
    pull_request:
        types: [ opened, synchronize ]

jobs:
    enrich:
        runs-on: ubuntu-latest
        permissions:
            pull-requests: write
        steps:
            -   name: Enrich PR Information
                uses: EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest@v3
                with:
                    repository: ${{ github.repository }}
                    pullRequestNumber: ${{ github.event.number }}
                    branch: ${{ github.head_ref }}
                    token: ${{ secrets.GITHUB_TOKEN }}
                    strategy: "jira"
                    jiraURL: ${{ vars.JIRA_URL }}
                    jiraEmail: ${{ vars.JIRA_EMAIL }}
                    jiraToken: ${{ secrets.JIRA_TOKEN }}

    validate:
        needs: enrich
        runs-on: ubuntu-latest
        steps:
            -   name: Validate PR
                run: echo "Running validation..."

    test:
        needs: enrich
        runs-on: ubuntu-latest
        steps:
            -   name: Run Tests
                run: echo "Running tests..."
```

## Branch Name Patterns

### Supported Formats

The branch-name strategy recognizes these patterns:

**With Issue Type Prefix:**

```
epic/PROJ-123-epic-title
feature/PROJ-456-new-feature
bugfix/PROJ-789-bug-description
hotfix/PROJ-101-critical-fix
```

**Without Issue Type Prefix:**

```
PROJ-123-task-description
ISSUE-456-improvement-title
TICKET-789-maintenance-task
```

### Pattern Matching

- **Issue Key**: `[A-Z]+-[0-9]+` (e.g., PROJ-123, ISSUE-456)
- **Description**: Hyphen-separated words after issue key
- **Type Prefixes**: epic, feature, bugfix, hotfix (optional)

## Jira Integration Details

### Authentication

```yaml
# Using organization variables and secrets
jiraURL: ${{ vars.JIRA_URL }}           # https://yourorg.atlassian.net
jiraEmail: ${{ vars.JIRA_EMAIL }}       # someone@yourcompany.com
jiraToken: ${{ secrets.JIRA_TOKEN }}    # Jira API token or PAT
```

### Issue Information

The Jira strategy fetches and uses:

- **Issue Summary**: Used as PR title
- **Issue Description**: Optionally synced to PR description
- **Parent Issue**: Added as prefix for subtasks (excluding epics)
- **Issue Status**: Used for validation

### Label Management

When `jiraEnableSyncLabel: true`:

- Creates label if it doesn't exist
- Label description: "Indicates that Jira synchronization has been completed for this PR"
- Prevents duplicate syncing on subsequent runs

## Required Permissions

The GitHub token must have the following permissions:

```yaml
permissions:
    pull-requests: write  # Required for updating PR title/description and labels
    issues: write         # Required for creating the sync label if it doesn't exist
    contents: read        # Required for accessing repository and branch information
```

## Custom Formatting Rules

Define custom word formatting with key-value pairs:

### Format

```
"key1:value1,key2:value2,key3:value3"
```

### Examples

```yaml
# Technical abbreviations
customFormatting: "api:API,ui:UI,db:Database,auth:Authentication,ci:Continuous Integration"

# Business domains
customFormatting: "crm:Customer Relations,hr:Human Resources,fin:Finance,ops:Operations"

# Component names
customFormatting: "comp:Component,svc:Service,lib:Library,util:Utility,cfg:Configuration"
```

## Troubleshooting

### Common Issues

**PR Not Updated**

- Verify token has `pull-requests: write` permission
- Check pull request number is correct
- Ensure branch name matches expected patterns

**Jira Connection Issues**

- Verify Jira URL format: `https://yourorg.atlassian.net`
- Check Jira token permissions
- Ensure issue key exists in Jira

**Branch Pattern Mismatch**

- Verify branch name follows supported patterns
- Check regex matching in action logs
- Test with simpler branch names

### Error Resolution

**Authentication Failed**

```yaml
# Ensure proper Jira authentication
jiraToken: ${{ secrets.JIRA_API_TOKEN }}  # Use API token, not password
```

**Issue Not Found**

```bash
# Verify issue exists and is accessible
curl -H "Authorization: Bearer $JIRA_TOKEN" \
     "https://yourorg.atlassian.net/rest/api/3/issue/PROJ-123"
```

**Permission Denied**

```yaml
# Ensure proper github permissions
permissions:
    pull-requests: write
    issues: write
    contents: read
```

## Best Practices

### Branch Naming Conventions

1. **Consistent Prefixes**: Use standard issue type prefixes
2. **Issue Keys**: Always include project issue keys
3. **Descriptive Names**: Use clear, hyphen-separated descriptions
4. **Lowercase**: Use lowercase for consistency

### Jira Integration

1. **Service Account**: Create a service account in your Jira system dedicated to integrations.
2. **API Tokens**: Use Jira API tokens instead of passwords
3. **Issue Templates**: Maintain consistent issue title formats
4. **Parent Relationships**: Properly structure issue hierarchies