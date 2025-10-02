# PHP Duster Action

## Overview

The `actions/php/duster` action applies Laravel Pint code style fixes to PHP projects. It focuses specifically on code style enforcement using Duster (Laravel Pint) and
automatically commits changes while managing git blame history.

## Language/Tool Support

- **PHP**: All PHP projects
- **Laravel**: Optimized for Laravel projects with Laravel Pint
- **Duster**: Laravel Pint-based code style enforcement

## Features

- **Laravel Pint Integration**: Uses Duster for Laravel Pint code style fixes
- **Automatic Commits**: Commits changes with descriptive "Dusting" message
- **Git Blame Management**: Automatically adds commit hashes to `.git-blame-ignore-revs`
- **Directory Support**: Can target subdirectories within repositories
- **Composer Integration**: Handles dependency caching
- **Force Push**: Ensures changes are pushed to the target branch

## Usage

```yaml
- name: Apply Duster Code Style
  uses: ./actions/php/duster
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

Your project must include Duster in `composer.json`:

```json
{
  "require-dev": {
    "tightenco/duster": "^2.0"
  }
}
```

## Action Steps

1. **Setup PHP**: Configures the specified PHP version
2. **Checkout Repository**: Retrieves code from the specified repository and branch
3. **Set Permissions**: Makes the entrypoint script executable
4. **Determine Working Directory**: Handles workspace vs subdirectory configuration
5. **Cache Dependencies**: Caches Composer dependencies for performance
6. **Apply Code Style**: Runs the entrypoint script that executes Duster
7. **Force Cache**: Saves cache if it wasn't hit

## Code Style Process

### Duster Execution

1. Runs `duster fix` to apply Laravel Pint code style rules
2. Commits changes with message "Dusting"
3. Adds commit hash to `.git-blame-ignore-revs`
4. Creates additional commit to ignore Duster changes in git blame
5. Force pushes changes to origin

## Configuration

### Duster Configuration (duster.json)

```json
{
    "preset": "laravel",
    "rules": {
        "simplified_null_return": true,
        "not_operator_with_successor_space": true,
        "binary_operator_spaces": {
            "default": "single_space"
        }
    }
}
```

### Alternative: Pint Configuration (pint.json)

```json
{
    "preset": "laravel",
    "rules": {
        "no_unused_imports": true,
        "ordered_imports": {
            "sort_algorithm": "alpha"
        }
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
      - name: Apply Duster Code Style
        uses: ./actions/php/duster
        with:
          repository: ${{ github.repository }}
          branch: ${{ github.head_ref }}
```

### Subdirectory Usage

```yaml
- name: Apply Duster to API
  uses: ./actions/php/duster
  with:
    repository: ${{ github.repository }}
    branch: ${{ github.head_ref }}
    directory: 'api'
    phpVersion: '8.2'
```

### Multiple Directories

```yaml
strategy:
  matrix:
    directory: ['backend', 'admin-panel', 'api']

steps:
  - name: Apply Duster Code Style
    uses: ./actions/php/duster
    with:
      repository: ${{ github.repository }}
      branch: ${{ github.head_ref }}
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
2. **Adding Commit Hash**: The Duster commit is added to the ignore file
3. **Separate Ignore Commit**: Creates dedicated commit for updating the ignore file
4. **Git Configuration**: Teams can configure git to use this file:

```bash
git config blame.ignoreRevsFile .git-blame-ignore-revs
```

## Duster vs Laravel Pint

Duster is a wrapper around Laravel Pint that provides:

- **Enhanced Configuration**: Additional configuration options
- **Better Integration**: Improved CI/CD integration
- **Extended Rules**: Additional code style rules beyond standard Pint
- **Performance**: Optimized for large codebases

## Code Style Rules

Common Duster/Pint rules include:

### PHP Standards

- PSR-12 compliance
- Array syntax normalization
- Method chaining alignment
- Import ordering

### Laravel-Specific

- Blade directive formatting
- Eloquent method chaining
- Route definition formatting
- Configuration array formatting

## Troubleshooting

### Common Issues

**Duster Installation Fails**

- Ensure `tightenco/duster` is in `composer.json`
- Verify Composer dependencies are properly configured

**Code Style Fixes Fail**

- Check for syntax errors in PHP files
- Verify Duster configuration is valid
- Ensure proper file permissions

**Git Push Failures**

- Confirm workflow has `contents: write` permission
- Check branch protection rules allow automated commits
- Verify branch exists and is accessible

**Configuration Issues**

- Validate `duster.json` or `pint.json` syntax
- Check preset compatibility with codebase
- Verify rule configuration matches project needs

### Performance Optimization

**Caching**

- Action automatically caches Composer dependencies
- Cache keys based on `composer.lock` file hash
- Force cache save ensures updated cache

**Selective Processing**

- Configure Duster to process specific directories
- Exclude vendor and build directories
- Use `.dusterignore` file for exclusions

## Security Considerations

- **Force Push**: Action performs force pushes; configure branch protection appropriately
- **Automated Commits**: Review automated commits in sensitive repositories
- **Token Permissions**: Use minimal required permissions

## Related Actions

- **actions/php/codeStyle**: Combined Rector and Duster action
- **actions/php/rector**: Standalone Rector refactoring action
- **php_staticAnalysis.yml**: Code quality analysis workflow
- **php_test.yml**: Unit testing workflow

## Best Practices

1. **Regular Application**: Run on every pull request
2. **Review Changes**: Always review automated code style changes
3. **Custom Rules**: Configure rules to match team preferences
4. **Pre-commit Hooks**: Consider local pre-commit hooks for immediate feedback
5. **Documentation**: Document custom code style rules for the team
6. **Integration**: Combine with static analysis and testing workflows

## Migration Scenarios

### From Manual Pint Usage

Replace manual Laravel Pint commands:

```bash
# Before
./vendor/bin/pint

# After
# Use this action in GitHub workflows
```

### From Other Code Style Tools

Migrating from PHP-CS-Fixer or similar tools:

1. Remove old tool configurations
2. Add Duster dependency
3. Configure Duster rules to match existing style
4. Update CI/CD workflows to use this action

## Configuration Examples

### Basic Laravel Project

```json
{
    "preset": "laravel"
}
```

### Custom Rules

```json
{
    "preset": "laravel",
    "rules": {
        "array_syntax": {
            "syntax": "short"
        },
        "binary_operator_spaces": {
            "default": "single_space"
        },
        "blank_line_before_statement": {
            "statements": ["return"]
        }
    }
}
```

### Exclusions

```json
{
    "preset": "laravel",
    "exclude": [
        "bootstrap/cache",
        "storage",
        "vendor"
    ]
}
```