# GitHub Create Release Action

## Overview

The `actions/github/createRelease` action creates GitHub releases with automated changelog generation and configurable release options. This Docker-based action provides
comprehensive release management capabilities including pre-release flags, draft releases, and Dependabot PR inclusion control.

## Language/Tool Support

- **GitHub**: Release management and API integration
- **Docker**: Containerized execution for consistent behavior
- **Changelog**: Automated release notes generation
- **Semantic Versioning**: Tag-based release naming

## Features

- **Automated Release Creation**: Creates GitHub releases from tags
- **Release Notes Generation**: Automatically generates changelog from commits
- **Pre-release Support**: Mark releases as pre-release versions
- **Draft Releases**: Create draft releases for review before publishing
- **Dependabot Integration**: Option to include/exclude Dependabot PRs from notes
- **Docker-based**: Consistent execution environment

## Usage

```yaml
- name: Create github Release
  uses: ./actions/github/createRelease
  with:
    token: ${{ secrets.GITHUB_TOKEN }}
    repository: ${{ github.repository }}
    tagName: ${{ github.ref_name }}
    preRelease: false
    generateReleaseNotes: true
    isDraft: false
    includeDependabot: false
```

## Inputs

| Input                  | Type    | Required | Default | Description                                    |
|------------------------|---------|----------|---------|------------------------------------------------|
| `token`                | string  | ✅        | -       | GitHub token with repository write permissions |
| `repository`           | string  | ✅        | -       | GitHub repository in format "owner/repo"       |
| `tagName`              | string  | ✅        | -       | Git tag name for the release                   |
| `preRelease`           | boolean | ❌        | `false` | Mark release as pre-release                    |
| `generateReleaseNotes` | boolean | ❌        | `true`  | Generate automatic release notes               |
| `isDraft`              | boolean | ❌        | `false` | Create release as draft                        |
| `includeDependabot`    | boolean | ❌        | `false` | Include Dependabot PRs in release notes        |

## Action Implementation

This action runs in a Docker container:

- **Image**: `ghcr.io/encoredigitalgroup/gh-action-create-github-release:latest`
- **Environment**: All inputs are passed as environment variables
- **Execution**: Containerized Go application

## Usage Examples

### Basic Release Creation

```yaml
name: Create Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Create Release
        uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: ${{ github.ref_name }}
```

### Pre-release Creation

```yaml
name: Create Pre-release
on:
  push:
    tags:
      - 'v*-rc*'
      - 'v*-beta*'
      - 'v*-alpha*'

jobs:
  prerelease:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Create Pre-release
        uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: ${{ github.ref_name }}
          preRelease: true
          generateReleaseNotes: true
```

### Draft Release for Review

```yaml
name: Create Draft Release
on:
  workflow_dispatch:
    inputs:
      tag:
        description: 'Release tag'
        required: true
        type: string

jobs:
  draft-release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Create Draft Release
        uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: ${{ inputs.tag }}
          isDraft: true
          generateReleaseNotes: true
          includeDependabot: false
```

### Release with Dependabot Changes

```yaml
name: Release with Dependencies
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Create Release with Dependabot
        uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: ${{ github.ref_name }}
          generateReleaseNotes: true
          includeDependabot: true
```

## Integration Patterns

### Automated Release Pipeline

```yaml
name: Release Pipeline
on:
  push:
    tags:
      - 'v[0-9]+.[0-9]+.[0-9]+'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Build Application
        run: echo "Building..."

  test:
    runs-on: ubuntu-latest
    steps:
      - name: Run Tests
        run: echo "Testing..."

  release:
    needs: [build, test]
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Create Release
        uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: ${{ github.ref_name }}

  deploy:
    needs: release
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Production
        run: echo "Deploying..."
```

### Multi-Environment Release

```yaml
name: Multi-Environment Release
on:
  push:
    tags:
      - 'v*'

jobs:
  determine-release-type:
    runs-on: ubuntu-latest
    outputs:
      is-prerelease: ${{ steps.check.outputs.prerelease }}
    steps:
      - id: check
        run: |
          if [[ "${{ github.ref_name }}" =~ (alpha|beta|rc) ]]; then
            echo "prerelease=true" >> $GITHUB_OUTPUT
          else
            echo "prerelease=false" >> $GITHUB_OUTPUT
          fi

  create-release:
    needs: determine-release-type
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Create Release
        uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: ${{ github.ref_name }}
          preRelease: ${{ needs.determine-release-type.outputs.is-prerelease }}
```

## Release Notes Generation

The action automatically generates release notes when `generateReleaseNotes: true`:

### Generated Content Includes

- **All Pull Requests**: All pull requests are automatically included.
- **Dependencies**: Dependabot updates (if `includeDependabot: true`), otherwise dependabot PR's are excluded.
- **Contributors**: List of contributors to the release

## Required Permissions

The GitHub token must have the following permissions:

```yaml
permissions:
  contents: write      # Required for creating releases
  pull-requests: read  # Required for generating release notes
```

## Use Cases

### Version Release Management

```yaml
# Semantic version releases
on:
  push:
    tags:
      - 'v[0-9]+.[0-9]+.[0-9]+'
      - '!v[0-9]+.[0-9]+.[0-9]+-*'  # Exclude pre-releases
```

### Hotfix Releases

```yaml
# Hotfix releases
on:
  push:
    tags:
      - 'hotfix-*'

jobs:
  hotfix-release:
    steps:
      - uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: ${{ github.ref_name }}
          generateReleaseNotes: false  # Manual notes for hotfixes
```

### Feature Branch Releases

```yaml
# Feature preview releases
on:
  workflow_dispatch:
    inputs:
      branch:
        description: 'Feature branch'
        required: true

jobs:
  feature-release:
    steps:
      - name: Create Feature Tag
        run: git tag "feature-${{ inputs.branch }}-$(date +%s)"

      - uses: ./actions/github/createRelease
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          repository: ${{ github.repository }}
          tagName: "feature-${{ inputs.branch }}-$(date +%s)"
          preRelease: true
          isDraft: true
```

## Troubleshooting

### Common Issues

**Release Creation Fails**

- Verify token has `contents: write` permission
- Ensure tag exists and is accessible
- Check repository exists and is accessible

**Release Notes Generation Issues**

- Verify token has `pull-requests: read` permission
- Check if repository has sufficient commit history
- Ensure pull requests are properly labeled

### Error Resolution

**Permission Denied**

```yaml
# Ensure proper permissions
permissions:
  contents: write
  pull-requests: read
```

**Tag Not Found**

```bash
# Verify tag exists
git tag -l
git push origin --tags
```

**Release Already Exists**

The action will fail if a release for the tag already exists. Consider:

- Using different tag names
- Deleting existing releases if recreating
- Implementing idempotent release logic

## Best Practices

### Tag Management

1. **Semantic Versioning**: Use consistent version tagging (v1.2.3)
2. **Tag Protection**: Protect release tags from deletion
3. **Automated Tagging**: Use workflow automation for tag creation
4. **Tag Validation**: Validate tag format before release creation

### Security Considerations

1. **Token Security**: Use repository secrets for tokens
2. **Permission Minimization**: Use minimal required permissions
3. **Audit Trail**: Monitor release creation activities
4. **Access Control**: Limit who can create releases