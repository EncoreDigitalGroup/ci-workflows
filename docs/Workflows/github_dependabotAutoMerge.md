# GitHub Dependabot Auto-Merge Workflow

## Overview

The `github_dependabotAutoMerge.yml` workflow automatically approves and merges Dependabot pull requests for semver-minor and semver-patch dependency updates. This
workflow helps maintain dependency freshness while reducing manual intervention for low-risk updates.

## Language/Tool Support

- **GitHub**: Dependabot integration
- **Dependencies**: All package managers supported by Dependabot
- **Merge Strategy**: Squash merging for clean commit history

## Features

- **Automatic Approval**: Automatically approves Dependabot PRs
- **Selective Auto-Merge**: Only merges minor and patch updates automatically
- **Metadata Analysis**: Uses Dependabot metadata to determine update type
- **Safe Merging**: Excludes major version updates from auto-merge
- **Merge Group Support**: Works with GitHub merge groups

## Triggers

- **workflow_call**: Can be called from other workflows
- **merge_group**: Triggered by GitHub merge group events

## Usage

```yaml
uses: ./.github/workflows/github_dependabotAutoMerge.yml
```

## Required Permissions

```yaml
permissions:
  pull-requests: write  # Required for approving PRs
  contents: write       # Required for merging PRs
```

## Workflow Steps

1. **Dependabot Detection**: Only runs if the actor is `dependabot[bot]`
2. **Metadata Extraction**: Fetches Dependabot metadata including update type
3. **Approval**: Automatically approves the pull request
4. **Conditional Auto-Merge**: Merges PRs based on semantic version update type:
    - **semver-minor**: Minor version updates (1.2.0 → 1.3.0)
    - **semver-patch**: Patch updates (1.2.0 → 1.2.1)
    - **semver-major**: Major updates (excluded from auto-merge)

## Update Types

### Auto-Merged Updates

- **Minor Updates**: New features, backwards-compatible
- **Patch Updates**: Bug fixes, security patches

### Manual Review Required

- **Major Updates**: Breaking changes requiring manual review
- **Unknown Update Types**: Updates that don't match semver patterns

## Configuration Example

### Basic Dependabot Integration

```yaml
name: Dependabot Auto-Merge
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  auto-merge:
    if: ${{ github.actor == 'dependabot[bot]' }}
    uses: ./.github/workflows/github_dependabotAutoMerge.yml
```

### With Branch Protection

```yaml
name: CI Pipeline
on:
  pull_request:
    branches: [main]

jobs:
  tests:
    runs-on: ubuntu-latest
    steps:
      - name: Run Tests
        run: npm test

  dependabot-auto-merge:
    needs: tests
    if: ${{ github.actor == 'dependabot[bot]' && success() }}
    uses: ./.github/workflows/github_dependabotAutoMerge.yml
```

## Best Practices

### Repository Setup

1. **Configure Dependabot**: Set up `.github/dependabot.yml`
2. **Branch Protection**: Require status checks before merging
3. **Auto-Merge Settings**: Enable auto-merge in repository settings

### Security Considerations

- **Review Major Updates**: Always manually review breaking changes
- **Monitor Dependencies**: Regularly audit dependency updates
- **Test Coverage**: Ensure comprehensive tests before enabling auto-merge

### Example Dependabot Configuration

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    reviewers:
      - "security-team"
    assignees:
      - "maintainer"
```

## Integration Patterns

### With Testing Workflows

```yaml
name: Dependabot Workflow
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Tests
        run: npm test

  auto-merge:
    needs: test
    if: ${{ github.actor == 'dependabot[bot]' && success() }}
    uses: ./.github/workflows/github_dependabotAutoMerge.yml
```

### With Security Scanning

```yaml
name: Security and Auto-Merge
on:
  pull_request:

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Security Audit
        run: npm audit

  dependabot-merge:
    needs: security-scan
    if: ${{ github.actor == 'dependabot[bot]' && success() }}
    uses: ./.github/workflows/github_dependabotAutoMerge.yml
```

## Troubleshooting

### Common Issues

**Auto-Merge Not Working**

- Verify repository auto-merge is enabled
- Check branch protection rules allow auto-merge
- Ensure required status checks are passing

**PRs Not Being Approved**

- Confirm workflow has `pull-requests: write` permission
- Verify Dependabot actor detection is working
- Check if workflow is triggered correctly

**Merge Conflicts**

- Dependabot will automatically rebase PRs
- Manual intervention may be needed for complex conflicts
- Consider configuring Dependabot rebase strategy

### Monitoring and Alerts

**Failed Auto-Merges**

- Set up notifications for workflow failures
- Monitor merge queue for blocked PRs
- Review failed auto-merge attempts regularly

**Security Considerations**

- Audit auto-merged dependencies regularly
- Set up security alerts for vulnerable dependencies
- Consider additional approval for security-related updates

## Related Workflows

- **github_createRelease.yml**: For managing releases after dependency updates
- **Security scanning workflows**: For vulnerability assessment
- **Testing workflows**: For validating dependency updates

## Migration Guide

### From Manual Dependabot Management

1. Enable repository auto-merge feature
2. Configure branch protection rules
3. Add this workflow to your repository
4. Test with a sample Dependabot PR

### Customization Options

**Custom Update Types**

Modify the workflow to handle additional update types:

```yaml
- name: Auto-merge custom updates
  if: ${{ steps.metadata.outputs.package-ecosystem == 'npm' }}
  run: gh pr merge --auto --squash "$PR_URL"
```

**Custom Merge Strategy**

Change from squash to merge or rebase:

```yaml
- name: Merge with rebase
  run: gh pr merge --auto --rebase "$PR_URL"
```