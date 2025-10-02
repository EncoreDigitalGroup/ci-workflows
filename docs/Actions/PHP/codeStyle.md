# PHP Code Style Action

## Overview

The `actions/php/codeStyle` action applies comprehensive code style fixes and refactoring to PHP projects by combining both Rector (automated refactoring) and Duster
(Laravel Pint code style fixes). It automatically commits changes and manages git blame history.

## Language/Tool Support

- **PHP**: All PHP projects
- **Laravel**: Full support with Laravel Pint integration
- **Rector**: Automated refactoring and modernization
- **Duster**: Laravel Pint-based code style enforcement

## Features

- **Dual Tool Integration**: Combines Rector and Duster for comprehensive code improvement
- **Automatic Commits**: Commits changes separately for each tool with descriptive messages
- **Git Blame Management**: Automatically adds commit hashes to `.git-blame-ignore-revs`
- **Directory Support**: Can target subdirectories within repositories
- **Composer Integration**: Handles dependency caching
- **Force Push**: Ensures changes are pushed to the target branch

## Usage

```yaml
- name: Apply Code Style
  uses: ./actions/php/codeStyle
  with:
    repository: ${{ github.repository }}
    branch: ${{ github.head_ref }}
    phpVersion: '8.3'
    directory: '.'
```

## Inputs

| Input              | Type   | Required | Default                   | Description                                       |
|--------------------|--------|----------|---------------------------|---------------------------------------------------|
| `repository`       | string | ✅        | -                         | GitHub repository in format "owner/repo"          |
| `branch`           | string | ❌        | `main`                    | Target branch for applying code style changes     |
| `phpVersion`       | string | ❌        | `8.3`                     | PHP version to use for code style operations      |
| `directory`        | string | ❌        | `${{ github.workspace }}` | Directory path relative to workspace              |

## Required Dependencies

Your project must include both tools in `composer.json`:

```json
{
  "require-dev": {
    "tightenco/duster": "^2.0",
    "rector/rector": "^1.0"
  }
}
```

## Action Steps

1. **Setup PHP**: Configures the specified PHP version
2. **Checkout Repository**: Retrieves code from the specified repository and branch
3. **Set Permissions**: Makes the entrypoint script executable
4. **Determine Working Directory**: Handles workspace vs subdirectory configuration
5. **Cache Dependencies**: Caches Composer dependencies for performance
6. **Apply Code Style**: Runs the entrypoint script that executes both tools
7. **Force Cache**: Saves cache if it wasn't hit

## Code Style Process

### Phase 1: Rector (Automated Refactoring)

1. Runs `rector process` to apply automated refactoring rules
2. Commits changes with message "Rectifying"
3. Adds commit hash to `.git-blame-ignore-revs`
4. Creates additional commit to ignore Rector changes in git blame
5. Force pushes changes to origin

### Phase 2: Duster (Code Style Fixes)

1. Runs `duster fix` to apply Laravel Pint code style rules
2. Commits changes with message "Dusting"
3. Adds commit hash to `.git-blame-ignore-revs`
4. Creates additional commit to ignore Duster changes in git blame
5. Force pushes changes to origin

## Configuration Files

### Rector Configuration (rector.php)

```php
<?php

declare(strict_types=1);

use Rector\Config\RectorConfig;
use Rector\Set\ValueObject\LevelSetList;

return static function (RectorConfig $rectorConfig): void {
    $rectorConfig->paths([
        __DIR__ . '/app',
        __DIR__ . '/src',
    ]);

    $rectorConfig->sets([
        LevelSetList::UP_TO_PHP_83,
    ]);
};
```

### Duster Configuration (duster.json or pint.json)

```json
{
    "preset": "laravel",
    "rules": {
        "simplified_null_return": true,
        "not_operator_with_successor_space": true
    }
}
```

## Usage Examples

### Basic Usage in Workflow

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
```

### Subdirectory Usage

```yaml
- name: Apply Code Style to Backend
  uses: ./actions/php/codeStyle
  with:
    repository: ${{ github.repository }}
    branch: ${{ github.head_ref }}
    directory: 'backend'
    phpVersion: '8.2'
```

### Multiple PHP Versions

```yaml
strategy:
  matrix:
    php-version: ['8.2', '8.3']
    directory: ['api', 'admin', 'frontend']

steps:
  - name: Apply Code Style
    uses: ./actions/php/codeStyle
    with:
      repository: ${{ github.repository }}
      branch: ${{ github.head_ref }}
      phpVersion: ${{ matrix.php-version }}
      directory: ${{ matrix.directory }}
```

## Required Permissions

The workflow using this action must have:

```yaml
permissions:
  contents: write  # Required for committing and pushing changes
```

## Git Blame Integration

The action automatically manages git blame history by:

1. **Creating `.git-blame-ignore-revs`**: File containing commit hashes to ignore
2. **Adding Commit Hashes**: Each code style commit is added to the ignore file
3. **Separate Ignore Commits**: Creates dedicated commits for updating the ignore file
4. **Git Configuration**: Teams can configure git to use this file:

```bash
git config blame.ignoreRevsFile .git-blame-ignore-revs
```

## Troubleshooting

### Common Issues

**Rector Failures**

- Ensure `rector.php` configuration is valid
- Check that Rector rules are compatible with your PHP version
- Verify project structure matches configured paths

**Duster Failures**

- Confirm Duster/Pint configuration is correct
- Check for syntax errors that prevent code style fixes
- Ensure proper Laravel Pint preset is specified

**Git Push Failures**

- Confirm workflow has `contents: write` permission
- Check that the target branch allows force pushes
- Verify branch protection rules don't prevent automated commits

**Directory Issues**

- Ensure the specified directory exists in the repository
- Check that `composer.json` exists in the target directory
- Verify relative paths are correct

### Performance Optimization

**Caching**

- The action automatically caches Composer dependencies
- Cache keys are based on `composer.lock` file hash
- Force cache save ensures cache is updated when needed

**Dependency Installation**

- Uses `--no-ansi --no-interaction --no-progress` for faster installation
- Includes `--prefer-dist --ignore-platform-reqs` for reliability

## Security Considerations

- **Force Push**: Action performs force pushes; ensure branch protection is configured appropriately
- **Automated Commits**: Review automated commits in sensitive repositories
- **Token Permissions**: Use minimal required permissions for Composer authentication

## Related Actions

- **actions/php/rector**: Standalone Rector refactoring
- **actions/php/duster**: Standalone Duster code style fixes
- **php_staticAnalysis.yml**: Code quality analysis workflow
- **php_test.yml**: Unit testing workflow

## Best Practices

1. **Run on Pull Requests**: Apply code style fixes before merging
2. **Review Changes**: Always review automated code style changes
3. **Configure Rules**: Customize Rector and Duster rules for your project
4. **Test After Changes**: Run tests after code style changes
5. **Branch Protection**: Configure appropriate branch protection rules
6. **Regular Updates**: Keep Rector and Duster dependencies updated

## Migration from php_dusterFix.yml

This action replaces the deprecated `php_dusterFix.yml` workflow with enhanced features:

- **Added Rector**: Now includes automated refactoring
- **Better Git Blame**: Improved git blame ignore functionality
- **Enhanced Caching**: More efficient dependency caching
- **Directory Support**: Better handling of subdirectories
- **Separate Commits**: Individual commits for each tool