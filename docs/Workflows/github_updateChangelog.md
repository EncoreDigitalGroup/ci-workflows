# GitHub Update Changelog Workflow

## Overview

The `github_updateChangelog.yml` workflow automatically updates the project's CHANGELOG.md file when a new release is created. It extracts release information and appends
it to the changelog, maintaining a comprehensive history of project changes.

## Language/Tool Support

- **GitHub**: Release management integration
- **File Format**: Markdown changelog files
- **Automation**: Automated changelog maintenance

## Features

- **Automatic Updates**: Updates changelog on release creation
- **Release Notes Integration**: Incorporates GitHub release notes
- **Version Tracking**: Maintains chronological version history
- **Auto-Commit**: Automatically commits changelog updates
- **Main Branch Updates**: Updates changelog on the main branch

## Triggers

- **workflow_call**: Can be called from other workflows (typically release workflows)

## Usage

```yaml
uses: ./.github/workflows/github_updateChangelog.yml
```

## Required Permissions

```yaml
permissions:
  contents: write  # Required for updating files and committing changes
```

## Workflow Steps

1. **Checkout**: Checks out the repository on the main branch
2. **Extract Release Data**: Gets release name and body from GitHub release event
3. **Update Changelog**: Uses changelog-updater-action to format and add entry
4. **Commit Changes**: Commits the updated CHANGELOG.md file

## Integration Examples

### With Release Workflow

```yaml
name: Release Pipeline
on:
  release:
    types: [published]

jobs:
  update-changelog:
    uses: ./.github/workflows/github_updateChangelog.yml

  deploy:
    needs: update-changelog
    runs-on: ubuntu-latest
    steps:
      - name: Deploy Release
        run: echo "Deploying ${{ github.event.release.tag_name }}"
```

### Automated Release Process

```yaml
name: Automated Release
on:
  push:
    tags:
      - 'v*'

jobs:
  create-release:
    runs-on: ubuntu-latest
    steps:
      - name: Create Release
        uses: actions/create-release@v1
        with:
          tag_name: ${{ github.ref }}
          release_name: Release ${{ github.ref }}
          generate_release_notes: true

  update-changelog:
    needs: create-release
    uses: ./.github/workflows/github_updateChangelog.yml
```

## Changelog Format

The workflow maintains a standard changelog format compatible with [Keep a Changelog](https://keepachangelog.com/):

```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [v2.1.0] - 2024-01-15

### Added
- New feature for user authentication
- Support for custom themes

### Fixed
- Fixed bug in data validation
- Resolved memory leak issue

### Changed
- Updated API endpoints
- Improved error handling

## [v2.0.0] - 2024-01-01

### Added
- Complete rewrite of core functionality
- New dashboard interface

### Breaking Changes
- Removed deprecated API endpoints
- Changed configuration file format
```

## Release Notes Integration

The workflow automatically integrates GitHub release notes:

### Example Release Creation

```yaml
- name: Create Release with Notes
  uses: actions/create-release@v1
  with:
    tag_name: v1.2.0
    release_name: "Version 1.2.0 - Feature Update"
    body: |
      ## What's New
      - Added user profile management
      - Improved search functionality

      ## Bug Fixes
      - Fixed login redirect issue
      - Resolved mobile layout problems

      ## Breaking Changes
      - Updated API authentication method
```

## Configuration Options

### Changelog Location

The workflow expects a `CHANGELOG.md` file in the repository root. For custom locations:

```yaml
# Custom configuration would require workflow modification
# Default: CHANGELOG.md in repository root
```

### Commit Message Customization

The workflow uses a standard commit message: "Update Changelog"

For custom commit messages, you can modify the workflow:

```yaml
- name: Commit Changelog
  uses: stefanzweifel/git-auto-commit-action@v5
  with:
    commit_message: "docs: update changelog for ${{ github.event.release.name }}"
    file_pattern: CHANGELOG.md
```

## Use Cases

### Version Release Management

```yaml
name: Release Management
on:
  workflow_dispatch:
    inputs:
      version:
        description: 'Release version'
        required: true
        type: string

jobs:
  create-release:
    runs-on: ubuntu-latest
    steps:
      - name: Create GitHub Release
        # Create release steps...

  update-docs:
    needs: create-release
    uses: ./.github/workflows/github_updateChangelog.yml
```

### Hotfix Release Process

```yaml
name: Hotfix Release
on:
  push:
    branches:
      - 'hotfix/*'

jobs:
  hotfix-release:
    if: github.event_name == 'push'
    runs-on: ubuntu-latest
    steps:
      - name: Create Hotfix Release
        # Release creation steps...

  update-changelog:
    needs: hotfix-release
    uses: ./.github/workflows/github_updateChangelog.yml
```

## Best Practices

### Release Notes Quality

1. **Structured Format**: Use consistent formatting for release notes
2. **User-Focused**: Write notes from user perspective
3. **Categorization**: Group changes by type (Added, Fixed, Changed, etc.)
4. **Breaking Changes**: Clearly highlight breaking changes

### Changelog Maintenance

1. **Regular Updates**: Keep changelog current with each release
2. **Clear Versioning**: Use semantic versioning for releases
3. **Link to Issues**: Reference relevant issues and pull requests
4. **Migration Guides**: Include migration instructions for breaking changes

### Automation Integration

1. **Pre-Release Checks**: Validate changelog format before release
2. **Notification**: Set up notifications for changelog updates
3. **Review Process**: Review auto-generated changelog entries
4. **Backup**: Maintain changelog backups

## Troubleshooting

### Common Issues

**Changelog Not Updating**

- Verify workflow is triggered by release events
- Check repository permissions for workflow
- Ensure CHANGELOG.md file exists in repository root

**Commit Failures**

- Confirm workflow has `contents: write` permission
- Check if main branch is protected and allows automated commits
- Verify git configuration in workflow

**Format Issues**

- Review release notes format
- Check changelog structure compatibility
- Ensure markdown formatting is correct

### Error Resolution

**Permission Denied**

```yaml
# Ensure proper permissions
permissions:
  contents: write
```

**File Not Found**

```bash
# Create initial CHANGELOG.md if missing
echo "# Changelog\n\nAll notable changes will be documented here." > CHANGELOG.md
```

**Merge Conflicts**

- The workflow targets main branch directly
- Resolve any existing changelog conflicts before release
- Consider branch protection rules

## Related Workflows

- **Release creation workflows**: For generating releases
- **github_createRelease.yml**: For automated release creation
- **Version tagging workflows**: For semantic versioning

## Customization Options

### Custom Changelog Format

Modify the workflow to use different changelog formats:

```yaml
- name: Custom Changelog Update
  run: |
    echo "## ${{ github.event.release.name }}" >> CHANGELOG.md
    echo "${{ github.event.release.body }}" >> CHANGELOG.md
```

### Multiple Changelog Files

For projects with multiple changelog files:

```yaml
- name: Update Multiple Changelogs
  run: |
    # Update main changelog
    # Update API changelog
    # Update UI changelog
```

### Notification Integration

Add notifications for changelog updates:

```yaml
- name: Notify Team
  uses: slack-action
  with:
    message: "Changelog updated for ${{ github.event.release.name }}"
```