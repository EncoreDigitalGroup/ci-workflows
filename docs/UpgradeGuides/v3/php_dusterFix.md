# PHP Duster Fix Workflow v3 Migration Guide

## Overview

The `php_dusterFix.yml` workflow has been **removed** in v3 and replaced with the more comprehensive `actions/php/codeStyle` action that combines both Rector and Duster
code style enforcement.

## Breaking Changes

### Workflow Removal

The dedicated Duster Fix workflow (`php_dusterFix.yml`) has been completely removed in v3. Code style enforcement is now handled by the unified `actions/php/codeStyle`
action.

### Enhanced Functionality

The replacement action provides:

- **Rector** automated refactoring
- **Duster** code style fixes via Laravel Pint
- Git blame ignore functionality for both tools
- Better directory support
- Improved caching

## Migration Steps

### 1. Update Workflow References

**Before (v2):**

```yaml
uses: ./.github/workflows/php_dusterFix.yml
with:
  phpVersion: '8.3'
  useComposerAuthJson: true
```

**After (v3):**

```yaml
- name: Apply Code Style
  uses: ./actions/php/codeStyle
  with:
    repository: ${{ github.repository }}
    branch: ${{ github.head_ref }}
    phpVersion: '8.3'  # Optional, defaults to 8.3
    directory: 'src'     # Optional, defaults to workspace root
```

### 2. Required Parameters

The new action requires these parameters:

- `repository`: The GitHub repository (usually `${{ github.repository }}`)
- `branch`: The branch to work on (usually `${{ github.head_ref }}` for PRs)

### 3. Update Repository Permissions

Ensure your workflow has the necessary permissions:

```yaml
permissions:
  contents: write
```

### 4. Complete Workflow Example

```yaml
name: Code Style

on:
  pull_request:
    branches: [ main ]

jobs:
  code-style:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Apply Code Style
        uses: ./actions/php/codeStyle
        with:
          repository: ${{ github.repository }}
          branch: ${{ github.head_ref }}
          phpVersion: '8.3'
```

## Key Differences

| Feature           | php_dusterFix.yml (v2) | actions/php/codeStyle (v3)      |
|-------------------|------------------------|---------------------------------|
| Tools             | Duster only            | Rector + Duster                 |
| Directory Support | Root only              | Configurable via `directory`    |
| Git Config        | Workflow-level         | Action-level                    |
| Caching           | Basic                  | Enhanced with proper cache keys |
| Commits           | Single commit          | Separate commits per tool       |
| Git Blame         | Single ignore entry    | Separate ignore entries         |

## What the New Action Does

1. **Rector Processing**: Runs automated refactoring first
    - Commits changes with message "Rectifying"
    - Adds commit hash to `.git-blame-ignore-revs`
    - Creates commit to ignore Rector changes in git blame

2. **Duster Processing**: Runs Laravel Pint code style fixes
    - Commits changes with message "Dusting"
    - Adds commit hash to `.git-blame-ignore-revs`
    - Creates commit to ignore Duster changes in git blame

3. **Force Push**: Pushes all changes to the branch

## Dependencies Required

Ensure your project has the necessary dependencies:

```json
{
  "require-dev": {
    "tightenco/duster": "^3.0",
    "rector/rector": "^2.0"
  }
}
```

## Migration Checklist

- [ ] Replace workflow calls to use `actions/php/codeStyle` action
- [ ] Update required parameters (`repository`, `branch`)
- [ ] Add necessary permissions (`contents: write`) to workflow
- [ ] Install Rector and Duster dependencies in your project
- [ ] Test the new action to ensure both Duster and Rector run correctly
- [ ] Verify that commits are properly ignored in git blame

## Benefits of Migration

1. **Enhanced Code Quality**: Both refactoring (Rector) and style fixes (Duster)
2. **Better Git History**: Proper git blame ignore for code style commits
3. **Improved Performance**: Better caching and dependency management
4. **Flexible Directory Support**: Can target subdirectories within repositories
5. **Unified Approach**: Single action for all PHP code style enforcement

## Troubleshooting

- **Missing Dependencies**: Ensure both `tightenco/duster` and `rector/rector` are in your `composer.json`
- **Permission Errors**: Verify the workflow has `contents: write` permission
- **Directory Issues**: Check that the `directory` parameter points to a valid path