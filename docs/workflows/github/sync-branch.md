# GitHub Sync Branch Workflow

## Overview

The `github_syncBranch.yml` workflow synchronizes a branch with the latest changes, creates a version file based on the latest Git tag, and commits the version
information. It's designed for maintaining branch synchronization and version tracking across repository branches with optional prefix support for version labeling.

## Language/Tool Support

- **Git Operations**: Branch synchronization and version management
- **GitHub Repositories**: Any repository with Git tags and branch structure
- **Version Management**: Automatic version file generation from Git tags

## Features

- **Branch Synchronization**: Syncs specified branch with latest changes
- **Version File Generation**: Creates `.version` file with latest Git tag
- **Force Push Support**: Optional force push for branch updates
- **Directory Support**: Can work in subdirectories within repositories
- **Version Prefixing**: Add custom prefixes to version values
- **Automated Commits**: Commits version file with tag-based commit messages

## Usage

```yaml
uses: ./.github/workflows/github_syncBranch.yml
with:
  branch: "develop"
  force: false
  directory: "backend"
  prefix: "v"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Input       | Type    | Required | Default                   | Description                                  |
|-------------|---------|----------|---------------------------|----------------------------------------------|
| `branch`    | string  | ✅        | -                         | The branch to synchronize                    |
| `force`     | boolean | ❌        | `false`                   | Enable force push for branch synchronization |
| `directory` | string  | ❌        | `${{ github.workspace }}` | Directory path relative to workspace         |
| `prefix`    | string  | ❌        | `''`                      | Prefix to add to version value               |

## Secrets

| Secret  | Required | Description                                    |
|---------|----------|------------------------------------------------|
| `token` | ✅        | GitHub token with repository write permissions |

## Workflow Steps

1. **Repository Checkout**: Checks out the repository
2. **Directory Determination**: Determines working directory (root or subdirectory)
3. **Prefix Configuration**: Sets up version prefix if provided
4. **Branch Synchronization**: Syncs the specified branch using EncoreDigitalGroup/action-sync-branch
5. **Tag Retrieval**: Gets the latest Git tag from the repository
6. **Version File Creation**: Creates/updates `.version` file with prefixed tag
7. **Auto Commit**: Commits the version file with the tag as commit message

## Example Configurations

### Basic Branch Sync

```yaml
uses: ./.github/workflows/github_syncBranch.yml
with:
  branch: "main"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

### Development Branch Sync with Prefix

```yaml
uses: ./.github/workflows/github_syncBranch.yml
with:
  branch: "develop"
  prefix: "dev-"
  force: false
secrets:
  token: ${{ secrets.SYNC_TOKEN }}
```

### Subdirectory Version Management

```yaml
uses: ./.github/workflows/github_syncBranch.yml
with:
  branch: "release"
  directory: "api"
  prefix: "API v"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

### Force Sync for Hotfix

```yaml
uses: ./.github/workflows/github_syncBranch.yml
with:
  branch: "hotfix"
  force: true
  prefix: "hotfix-"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

## Version File Format

The workflow creates a `.version` file in the specified directory containing:

**Without Prefix:**

```
1.2.3
```

**With Prefix:**

```
v1.2.3
```

**With Custom Prefix:**

```
API v @ 1.2.3
```

## Use Cases

### Release Management

Maintain version consistency across branches:

- Sync release branches with latest tags
- Update version files for deployment processes
- Track version history across environments

### Multi-Environment Deployment

Keep environment branches synchronized:

- Sync staging branch with production tags
- Maintain development branch version tracking
- Coordinate feature branch version alignment

### Microservices Versioning

Manage versions in microservice repositories:

- Track service versions in subdirectories
- Sync API versions across service boundaries
- Maintain version consistency in monorepos

### Automated Release Preparation

Prepare branches for automated releases:

- Update version files before release workflows
- Sync branches before deployment pipelines
- Maintain version tracking for rollback scenarios

## Requirements

- **Git Tags**: Repository must have at least one Git tag
- **Branch Permissions**: Token must have write access to the target branch
- **Repository Structure**: Target directory must exist if specified
- **Sync Action Access**: Must have access to EncoreDigitalGroup/action-sync-branch

## Troubleshooting

### Common Issues

**No Git Tags Found**

- Ensure the repository has at least one Git tag
- Verify tags are properly formatted (e.g., v1.0.0, 1.0.0)
- Check if tags are accessible from the workflow context

**Branch Sync Fails**

- Verify the target branch exists
- Check if branch protection rules prevent force pushes
- Ensure the token has appropriate repository permissions

**Permission Denied**

- Verify the GitHub token has `contents: write` permissions
- Check if the token has access to the target repository
- Ensure branch protection rules allow the bot user to push

**Directory Not Found**

- Verify the specified directory exists in the repository
- Check the directory path format (should be relative to workspace)
- Ensure the directory is accessible from the workflow

**Force Push Issues**

- Check if branch protection rules prevent force pushes
- Verify the `force` parameter is set to `true` when needed
- Ensure the sync action supports force push operations

### Version File Issues

**Version File Not Created**

- Check if the working directory is writable
- Verify the previous tag step completed successfully
- Ensure the directory path is correctly formatted

**Incorrect Version Format**

- Check the prefix configuration
- Verify the Git tag format matches expectations
- Review the prefix output in workflow logs

## Safety Considerations

1. **Force Push Warning**: Use `force: true` carefully as it can overwrite branch history
2. **Branch Protection**: Ensure critical branches have appropriate protection rules
3. **Token Permissions**: Use minimal required permissions for the GitHub token
4. **Backup Important Branches**: Consider backing up before force synchronization

## Related Workflows

- **github_directoryImport.yml**: For importing files that may need version tracking
- **php_test.yml**: For testing synchronized code before deployment