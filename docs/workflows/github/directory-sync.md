# GitHub Directory Sync Workflow

## Overview

The `github_directorySync.yml` workflow synchronizes directories between repositories by copying files from a source directory in the current repository to a target
directory in another repository. This workflow is useful for maintaining shared code, documentation, or configuration files across multiple repositories.

## Language/Tool Support

- **GitHub**: Repository file management
- **File Types**: All file types and directory structures
- **Cross-Repository**: Supports copying between different repositories

## Features

- **Directory Synchronization**: Copies entire directories between repositories
- **Flexible Target Naming**: Option to rename target directory
- **Custom Commit Messages**: Configurable commit messages for tracking changes
- **Force Push**: Overwrites existing files in target repository
- **Automated Commits**: Commits changes using a bot account

## Usage

```yaml
uses: ./.github/workflows/github_directorySync.yml
with:
  source: "src/shared"
  targetRepo: "organization/target-repo"
  targetDirectory: "shared-components"
  targetDirectoryName: "components"
  commitMessage: "[Automated] Sync shared components"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Input                 | Type   | Required | Default                     | Description                                     |
|-----------------------|--------|----------|-----------------------------|-------------------------------------------------|
| `source`              | string | ✅        | -                           | Source directory path to copy from              |
| `targetRepo`          | string | ✅        | -                           | Target repository in format "owner/repo"        |
| `targetDirectory`     | string | ✅        | -                           | Target directory path in destination repository |
| `targetDirectoryName` | string | ❌        | `empty_value`               | Custom name for the target directory            |
| `commitMessage`       | string | ❌        | `[Automated] DirectorySync` | Commit message for the synchronization          |

## Secrets

| Secret  | Required | Description                                             |
|---------|----------|---------------------------------------------------------|
| `token` | ✅        | GitHub token with write access to the target repository |

## Workflow Steps

1. **Checkout**: Checks out the source repository
2. **Directory Copy**: Uses a specialized action to copy files to target repository
3. **Commit**: Commits changes to target repository with specified message
4. **Force Push**: Pushes changes to target repository, overwriting conflicts

## Use Cases

### Shared Library Distribution

Distribute common libraries across multiple projects:

```yaml
uses: ./.github/workflows/github_directorySync.yml
with:
  source: "libs/common"
  targetRepo: "company/frontend-app"
  targetDirectory: "src/lib"
  commitMessage: "[Sync] Update common libraries"
secrets:
  token: ${{ secrets.SYNC_TOKEN }}
```

### Documentation Synchronization

Keep documentation synchronized across repositories:

```yaml
uses: ./.github/workflows/github_directorySync.yml
with:
  source: "docs/api"
  targetRepo: "company/documentation-site"
  targetDirectory: "content/api"
  targetDirectoryName: "api-docs"
  commitMessage: "[Docs] Sync API documentation"
secrets:
  token: ${{ secrets.DOCS_SYNC_TOKEN }}
```

### Configuration Management

Synchronize configuration files across environments:

```yaml
uses: ./.github/workflows/github_directorySync.yml
with:
  source: "config/production"
  targetRepo: "company/production-env"
  targetDirectory: "config"
  commitMessage: "[Config] Update production configuration"
secrets:
  token: ${{ secrets.CONFIG_SYNC_TOKEN }}
```

## Example Workflows

### Scheduled Synchronization

```yaml
name: Daily Sync
on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM

jobs:
  sync-shared-components:
    uses: ./.github/workflows/github_directorySync.yml
    with:
      source: "packages/shared"
      targetRepo: "company/mobile-app"
      targetDirectory: "shared"
      commitMessage: "[Automated] Daily sync of shared components"
    secrets:
      token: ${{ secrets.SYNC_TOKEN }}
```

### Multi-Target Synchronization

```yaml
name: Sync to Multiple Repos
on:
  push:
    paths:
      - 'shared/**'

jobs:
  sync-frontend:
    uses: ./.github/workflows/github_directorySync.yml
    with:
      source: "shared/components"
      targetRepo: "company/frontend"
      targetDirectory: "src/shared"
      commitMessage: "[Sync] Update shared components from ${{ github.sha }}"
    secrets:
      token: ${{ secrets.SYNC_TOKEN }}

  sync-mobile:
    uses: ./.github/workflows/github_directorySync.yml
    with:
      source: "shared/utils"
      targetRepo: "company/mobile"
      targetDirectory: "lib/shared"
      commitMessage: "[Sync] Update shared utilities from ${{ github.sha }}"
    secrets:
      token: ${{ secrets.SYNC_TOKEN }}
```

### Conditional Synchronization

```yaml
name: Conditional Sync
on:
  pull_request:
    types: [closed]

jobs:
  sync-on-merge:
    if: github.event.pull_request.merged == true
    uses: ./.github/workflows/github_directorySync.yml
    with:
      source: "dist/production"
      targetRepo: "company/deployment-repo"
      targetDirectory: "artifacts"
      commitMessage: "[Release] Sync production build from PR #${{ github.event.number }}"
    secrets:
      token: ${{ secrets.DEPLOY_TOKEN }}
```

## Security Considerations

### Token Permissions

- **Repository Access**: Token must have write access to target repository
- **Scope Limitation**: Use tokens with minimal required permissions
- **Secret Management**: Store tokens as repository secrets

### File Validation

- **Content Review**: Review synchronized content before deployment
- **Access Control**: Ensure appropriate access controls on target repositories
- **Audit Trail**: Monitor synchronization activities through commit history

## Required Permissions

The GitHub token must have the following permissions for the target repository:

```yaml
permissions:
  contents: write    # Required for creating commits
  metadata: read     # Required for repository access
```

## Bot Configuration

The workflow uses a bot account for commits:

- **Email**: `ghbot@encoredigitalgroup.com`
- **Username**: `EncoreBot`
- **Force Push**: Enabled to overwrite conflicts

## Troubleshooting

### Common Issues

**Synchronization Failures**

- Verify token has write access to target repository
- Check that source directory exists and contains files
- Ensure target repository is accessible

**Permission Errors**

- Confirm token scope includes target repository
- Verify repository settings allow token access
- Check branch protection rules don't block automated commits

**File Conflicts**

- The workflow uses force push to resolve conflicts
- Manual resolution may be needed for complex conflicts
- Consider coordination between teams making changes

### Monitoring

**Failed Synchronizations**

- Set up workflow failure notifications
- Monitor target repositories for unexpected changes
- Review commit history for synchronization activities

**Performance Optimization**

- Limit synchronization to specific file types if needed
- Consider using path filters to trigger only relevant syncs
- Monitor repository size growth from frequent syncs

## Best Practices

1. **Selective Synchronization**: Only sync necessary files and directories
2. **Clear Commit Messages**: Use descriptive commit messages for tracking
3. **Regular Monitoring**: Monitor target repositories for sync status
4. **Documentation**: Document sync relationships between repositories
5. **Testing**: Test synchronization in development environments first
6. **Coordination**: Coordinate with teams using target repositories

## Related Workflows

- **github_directoryImport.yml**: For importing external content
- **Release workflows**: For distributing synchronized content
- **Testing workflows**: For validating synchronized content

## Migration Scenarios

### From Manual File Copying

1. Identify frequently copied directories
2. Set up automated sync workflows
3. Document new synchronization processes
4. Train teams on new automated workflows

### From Monorepo to Multi-Repo

1. Identify shared components to synchronize
2. Create source-of-truth repositories
3. Set up sync workflows to distribute components
4. Gradually transition to multi-repo architecture