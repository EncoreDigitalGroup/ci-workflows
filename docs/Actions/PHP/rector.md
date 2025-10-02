# PHP Rector Action

## Overview

The `actions/php/rector` action applies automated refactoring and modernization to PHP projects using Rector. It focuses specifically on code modernization, type
improvements, and automated refactoring while automatically committing changes and managing git blame history.

## Language/Tool Support

- **PHP**: All PHP projects
- **Rector**: Automated refactoring and code modernization
- **Framework Agnostic**: Works with any PHP framework (Laravel, Symfony, etc.)

## Features

- **Automated Refactoring**: Uses Rector for comprehensive code modernization
- **Type Improvements**: Adds type declarations and improves type safety
- **Framework Upgrades**: Helps with framework version migrations
- **Automatic Commits**: Commits changes with descriptive "Rectifying" message
- **Git Blame Management**: Automatically adds commit hashes to `.git-blame-ignore-revs`
- **Composer Integration**: Handles dependency management
- **Force Push**: Ensures changes are pushed to the target branch

## Usage

```yaml
- name: Apply Rector Refactoring
  uses: ./actions/php/rector
  with:
    repository: ${{ github.repository }}
    branch: ${{ github.head_ref }}
    phpVersion: '8.3'
```

## Inputs

| Input              | Type   | Required | Default | Description                                       |
|--------------------|--------|----------|---------|---------------------------------------------------|
| `repository`       | string | ✅        | -       | GitHub repository in format "owner/repo"          |
| `branch`           | string | ❌        | `main`  | Target branch for applying refactoring changes    |
| `phpVersion`       | string | ❌        | `8.3`   | PHP version to use for refactoring operations     |

## Required Dependencies

Your project must include Rector in `composer.json`:

```json
{
  "require-dev": {
    "rector/rector": "^1.0"
  }
}
```

## Action Steps

1. **Setup PHP**: Configures the specified PHP version
2. **Checkout Repository**: Retrieves code from the specified repository and branch
3. **Set Permissions**: Makes the entrypoint script executable
4. **Run Rector**: Executes the refactoring process
5. **Auto Commit**: Commits changes if any modifications were made

## Refactoring Process

### Rector Execution

1. Installs Composer dependencies with optimized flags
2. Runs `rector process` to apply refactoring rules
3. Commits changes with message "Rectifying"
4. Adds commit hash to `.git-blame-ignore-revs`
5. Creates additional commit to ignore Rector changes in git blame
6. Force pushes changes to origin

## Configuration

### Basic Rector Configuration (rector.php)

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

### Advanced Configuration with Framework Support

```php
<?php

declare(strict_types=1);

use Rector\Config\RectorConfig;
use Rector\Set\ValueObject\LevelSetList;
use Rector\Laravel\Set\LaravelSetList;
use Rector\TypeDeclaration\Rector\ClassMethod\AddVoidReturnTypeWhereNoReturnRector;

return static function (RectorConfig $rectorConfig): void {
    $rectorConfig->paths([
        __DIR__ . '/app',
        __DIR__ . '/database',
        __DIR__ . '/tests',
    ]);

    $rectorConfig->sets([
        LevelSetList::UP_TO_PHP_83,
        LaravelSetList::LARAVEL_100,
    ]);

    $rectorConfig->rules([
        AddVoidReturnTypeWhereNoReturnRector::class,
    ]);

    $rectorConfig->skip([
        __DIR__ . '/bootstrap',
        __DIR__ . '/storage',
        __DIR__ . '/vendor',
    ]);
};
```

## Usage Examples

### Basic Usage in Workflow

```yaml
name: Code Refactoring

on:
  pull_request:
    branches: [ main ]

jobs:
  refactor:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Apply Rector Refactoring
        uses: ./actions/php/rector
        with:
          repository: ${{ github.repository }}
          branch: ${{ github.head_ref }}
```

### Scheduled Refactoring

```yaml
name: Weekly Refactoring

on:
  schedule:
    - cron: '0 2 * * 1'  # Every Monday at 2 AM

jobs:
  refactor:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Apply Rector Refactoring
        uses: ./actions/php/rector
        with:
          repository: ${{ github.repository }}
          branch: 'main'
          phpVersion: '8.3'
```

### Multiple PHP Versions Testing

```yaml
strategy:
  matrix:
    php-version: ['8.2', '8.3']

steps:
  - name: Apply Rector for PHP ${{ matrix.php-version }}
    uses: ./actions/php/rector
    with:
      repository: ${{ github.repository }}
      branch: ${{ github.head_ref }}
      phpVersion: ${{ matrix.php-version }}
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
2. **Adding Commit Hash**: The Rector commit is added to the ignore file
3. **Separate Ignore Commit**: Creates dedicated commit for updating the ignore file
4. **Git Configuration**: Teams can configure git to use this file:

```bash
git config blame.ignoreRevsFile .git-blame-ignore-revs
```

## Common Rector Rules

### PHP Modernization

- **Type Declarations**: Adds parameter and return types
- **Null Coalescing**: Converts isset() checks to null coalescing operator
- **Arrow Functions**: Converts simple closures to arrow functions
- **Constructor Property Promotion**: Modernizes constructor property assignments

### Framework-Specific

- **Laravel**: Migration to newer Laravel versions
- **Symfony**: Symfony framework upgrades
- **PHPUnit**: Test method improvements
- **Doctrine**: ORM modernization

### Code Quality

- **Dead Code Removal**: Removes unused code
- **Simplification**: Simplifies complex expressions
- **Performance**: Applies performance improvements
- **Standards**: Enforces coding standards

## Troubleshooting

### Common Issues

**Rector Installation Fails**

- Ensure `rector/rector` is in `composer.json`
- Verify Composer dependencies are properly configured

**Refactoring Fails**

- Check for syntax errors in PHP files
- Verify Rector configuration is valid
- Ensure sufficient memory limits for large projects

**Configuration Errors**

- Validate `rector.php` syntax
- Check path configurations match project structure
- Verify rule compatibility with PHP version

**Git Push Failures**

- Confirm workflow has `contents: write` permission
- Check branch protection rules allow automated commits
- Verify branch exists and is accessible

### Performance Optimization

**Memory Management**

- Increase PHP memory limit for large projects
- Use selective path configuration
- Exclude vendor and build directories

**Processing Optimization**

- Configure parallel processing when available
- Use caching for repeated runs
- Implement incremental processing for large codebases

## Security Considerations

- **Force Push**: Action performs force pushes; configure branch protection appropriately
- **Automated Commits**: Review automated refactoring changes
- **Code Review**: Always review Rector changes before merging

## Related Actions

- **actions/php/codeStyle**: Combined Rector and Duster action
- **actions/php/duster**: Standalone Duster code style action
- **php_staticAnalysis.yml**: Code quality analysis workflow
- **php_test.yml**: Unit testing workflow

## Best Practices

1. **Incremental Adoption**: Start with conservative rule sets
2. **Testing Integration**: Always run tests after refactoring
3. **Code Review**: Review all automated changes
4. **Configuration Management**: Version control Rector configuration
5. **Documentation**: Document custom rules and exclusions
6. **Regular Updates**: Keep Rector updated for latest rules

## Migration Scenarios

### Legacy PHP Versions

```php
// rector.php for PHP 7.4 to 8.3 migration
return static function (RectorConfig $rectorConfig): void {
    $rectorConfig->sets([
        LevelSetList::UP_TO_PHP_83,
    ]);
};
```

### Framework Upgrades

```php
// Laravel upgrade example
return static function (RectorConfig $rectorConfig): void {
    $rectorConfig->sets([
        LaravelSetList::LARAVEL_90,
        LaravelSetList::LARAVEL_100,
    ]);
};
```

## Rule Categories

### Level Sets

- **UP_TO_PHP_81**: Modernize to PHP 8.1
- **UP_TO_PHP_82**: Modernize to PHP 8.2
- **UP_TO_PHP_83**: Modernize to PHP 8.3

### Framework Sets

- **LARAVEL_XX**: Laravel version-specific upgrades
- **SYMFONY_XX**: Symfony version-specific upgrades
- **DOCTRINE_XX**: Doctrine ORM upgrades

### Quality Sets

- **CODE_QUALITY**: General code quality improvements
- **DEAD_CODE**: Dead code removal
- **TYPE_DECLARATION**: Type declaration additions

## Custom Rules

### Creating Custom Rules

```php
use Rector\Config\RectorConfig;
use App\Rector\CustomRule;

return static function (RectorConfig $rectorConfig): void {
    $rectorConfig->rule(CustomRule::class);
};
```

### Rule Configuration

```php
use Rector\Arguments\Rector\ClassMethod\ArgumentAdderRector;
use Rector\Arguments\ValueObject\ArgumentAdder;

return static function (RectorConfig $rectorConfig): void {
    $rectorConfig->ruleWithConfiguration(ArgumentAdderRector::class, [
        new ArgumentAdder('SomeClass', 'someMethod', 0, 'newArgument', 'string'),
    ]);
};
```