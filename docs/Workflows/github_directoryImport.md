# GitHub Directory Import Workflow

## Overview

The `github_directoryImport.yml` workflow imports a specific directory from one GitHub repository into another repository. It uses the
`encoredigitalgroup/directory-import-action` to copy files and directories from a source repository to a target location with automatic commit functionality.

## Language/Tool Support

- **GitHub Repositories**: Any public or private GitHub repository
- **Directory Operations**: File and folder copying across repositories

## Features

- **Cross-Repository Import**: Import directories from external repositories
- **Flexible Target Naming**: Option to rename the imported directory
- **Automatic Commits**: Commits imported files with customizable messages
- **Path Customization**: Specify exact source and target paths

## Usage

```yaml
uses: ./.github/workflows/github_directoryImport.yml
with:
  sourceRepository: "owner/source-repo"
  sourceDirectory: "src/components"
  targetDirectory: "frontend/components"
  targetDirectoryName: "ui-components"
  commitMessage: "[Automated] Import UI components from design system"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Input                 | Type   | Required | Default                     | Description                                                 |
|-----------------------|--------|----------|-----------------------------|-------------------------------------------------------------|
| `sourceRepository`    | string | ✅        | -                           | The source repository to import from (format: "owner/repo") |
| `sourceDirectory`     | string | ✅        | -                           | Path to the directory in the source repository              |
| `targetDirectory`     | string | ✅        | -                           | Target directory path in the current repository             |
| `targetDirectoryName` | string | ❌        | `empty_value`               | New name for the imported directory                         |
| `commitMessage`       | string | ❌        | `[Automated] DirectorySync` | Commit message for the import                               |

## Secrets

| Secret  | Required | Description                                     |
|---------|----------|-------------------------------------------------|
| `token` | ✅        | GitHub token with repository access permissions |

## Workflow Steps

1. **Checkout Target Repository**: Checks out the current repository where files will be imported
2. **Checkout Source Repository**: Clones the source repository to a temporary directory
3. **Directory Import**: Uses the directory-import-action to copy files from source to target

## Example Configurations

### Basic Directory Import

```yaml
uses: ./.github/workflows/github_directoryImport.yml
with:
  sourceRepository: "company/shared-components"
  sourceDirectory: "src/ui"
  targetDirectory: "components"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

### Import with Custom Naming

```yaml
uses: ./.github/workflows/github_directoryImport.yml
with:
  sourceRepository: "design-team/component-library"
  sourceDirectory: "packages/buttons"
  targetDirectory: "src/components"
  targetDirectoryName: "button-components"
  commitMessage: "Import button components from design system"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

### Documentation Import

```yaml
uses: ./.github/workflows/github_directoryImport.yml
with:
  sourceRepository: "docs-team/api-documentation"
  sourceDirectory: "openapi"
  targetDirectory: "docs/api"
  commitMessage: "Update API documentation from docs repository"
secrets:
  token: ${{ secrets.GITHUB_TOKEN }}
```

## Use Cases

### Component Library Sync

Keep shared UI components up-to-date across multiple projects:

- Import latest components from a central design system
- Maintain consistency across application frontends
- Automatically update shared utilities and styles

### Documentation Sync

Synchronize documentation across repositories:

- Import API specs from backend repositories
- Keep README files and guides updated
- Share configuration templates

### Configuration Management

Distribute configuration files and templates:

- Import CI/CD templates
- Share development environment configurations
- Distribute coding standards and linting rules

## Requirements

- **Repository Access**: The `token` must have read access to the source repository
- **Write Permissions**: The workflow requires `contents: write` permission on the target repository
- **Valid Paths**: Source and target directories must be valid relative paths

## Troubleshooting

### Common Issues

**Permission Denied**

- Verify the GitHub token has access to the source repository
- Check if the source repository is private and token has appropriate permissions
- Ensure the token has write access to the target repository

**Source Directory Not Found**

- Verify the source directory path exists in the source repository
- Check for typos in the repository name or directory path
- Ensure the source repository is accessible

**Import Fails**

- Check if target directory already exists and has conflicting files
- Verify the directory-import-action version is compatible
- Review the action logs for specific error messages

**Commit Issues**

- Ensure the workflow has `contents: write` permissions
- Check if there are branch protection rules preventing commits
- Verify the commit message format is valid

## Related Workflows

- **github_syncBranch.yml**: For branch synchronization across repositories
- **php_gitStatusCheck.yml**: For checking if imported files need processing