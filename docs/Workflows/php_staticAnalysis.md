# PHP Static Analysis Workflow

## Overview

The `php_staticAnalysis.yml` workflow performs static code analysis on PHP projects using PHPStan. It analyzes code quality, type safety, and potential issues without
executing the code. The workflow supports Laravel projects with environment decryption, custom analysis paths, and subdirectory execution with comprehensive Composer
dependency management.

## Language/Tool Support

- **PHP**: All PHP versions with PHPStan compatibility
- **Laravel**: Full support including encrypted environment handling
- **PHPStan**: Static analysis tool for PHP code quality and type checking
- **Composer**: Dependency management with private repository support

## Features

- **PHPStan Integration**: Comprehensive static analysis using PHPStan
- **Laravel Environment Support**: Decryption of encrypted testing environments
- **Flexible Path Analysis**: Configurable paths for targeted analysis
- **Subdirectory Support**: Can analyze specific directories within repositories
- **Dependency Caching**: Optimized Composer dependency caching
- **Custom Analysis Paths**: Configurable analysis scope

## Usage

```yaml
uses: ./.github/workflows/php_staticAnalysis.yml
with:
  branch: "main"
  phpVersion: "8.3"
  path: "app/ src/"
  useLaravelEnvDecryptionKey: true
  directory: "backend"
secrets:
  laravelEnvDecryptionKey: ${{ secrets.LARAVEL_ENV_DECRYPTION_KEY }}
```

## Inputs

| Input                        | Type    | Required | Default                   | Description                                       |
|------------------------------|---------|----------|---------------------------|---------------------------------------------------|
| `branch`                     | string  | ❌        | `main`                    | The branch to analyze                             |
| `path`                       | string  | ❌        | `app/ app_modules/`       | Paths to analyze (space-separated)                |
| `phpVersion`                 | string  | ❌        | `8.2`                     | PHP version to use for analysis                   |
| `useLaravelEnvDecryptionKey` | boolean | ❌        | `false`                   | Enable Laravel environment decryption             |
| `directory`                  | string  | ❌        | `${{ github.workspace }}` | Directory path relative to workspace              |

## Secrets

| Secret                    | Required | Description                                                                         |
|---------------------------|----------|-------------------------------------------------------------------------------------|
| `laravelEnvDecryptionKey` | ❌        | Laravel environment decryption key (required if `useLaravelEnvDecryptionKey: true`) |

## Workflow Steps

1. **PHP Setup**: Configures PHP environment with specified version
2. **Repository Checkout**: Checks out the source code from specified branch
3. **Directory Configuration**: Determines working directory (root or subdirectory)
4. **Composer Caching**: Configures dependency caching for performance
5. **Dependency Installation**: Installs PHP dependencies with platform requirements ignored
6. **Laravel Environment Setup**: Decrypts testing environment (if enabled)
7. **PHPStan Analysis**: Executes static analysis on specified paths

## Example Configurations

### Basic PHP Analysis

```yaml
uses: ./.github/workflows/php_staticAnalysis.yml
with:
  phpVersion: "8.3"
  path: "src/ app/"
```

### Laravel Project Analysis

```yaml
uses: ./.github/workflows/php_staticAnalysis.yml
with:
  branch: "develop"
  phpVersion: "8.3"
  useLaravelEnvDecryptionKey: true
  path: "app/ routes/ database/"
secrets:
  laravelEnvDecryptionKey: ${{ secrets.LARAVEL_ENV_DECRYPTION_KEY }}
```

### Subdirectory Analysis

```yaml
uses: ./.github/workflows/php_staticAnalysis.yml
with:
  directory: "api"
  phpVersion: "8.2"
  path: "src/ tests/"
```

### Multi-Path Analysis

```yaml
uses: ./.github/workflows/php_staticAnalysis.yml
with:
  phpVersion: "8.3"
  path: "app/ src/ lib/ packages/"
```

## PHPStan Configuration

The workflow expects PHPStan to be configured in your project. Create a `phpstan.neon` file:

```neon
# phpstan.neon
parameters:
    level: 8
    paths:
        - app
        - src
    excludePaths:
        - app/Console/Commands/stubs
        - tests
    checkMissingIterableValueType: false
    checkGenericClassInNonGenericObjectType: false
```

For Laravel projects:

```neon
# phpstan.neon
includes:
    - ./vendor/nunomaduro/larastan/extension.neon

parameters:
    level: 5
    paths:
        - app
    excludePaths:
        - app/Console/Commands/stubs
    checkMissingIterableValueType: false
    checkGenericClassInNonGenericObjectType: false
```

## Analysis Paths

The `path` parameter accepts space-separated directory paths:

- **Single Path**: `"app/"`
- **Multiple Paths**: `"app/ src/ lib/"`
- **Nested Paths**: `"app/Http/ app/Models/ app/Services/"`
- **Pattern Support**: Use paths that PHPStan can understand

## Laravel Environment Setup

When `useLaravelEnvDecryptionKey: true`, the workflow:

1. Decrypts the testing environment using Laravel's encryption
2. Moves `.env.testing` to `.env` for analysis context
3. Provides proper environment for Laravel-specific analysis

**Environment Preparation:**

```bash
# Encrypt testing environment
php artisan env:encrypt --env=testing

# The workflow will decrypt it using:
php artisan env:decrypt --env=testing --key=${{ secrets.laravelEnvDecryptionKey }}
```

## Use Cases

### Code Quality Assurance

Maintain high code quality standards:

- Detect type inconsistencies and potential bugs
- Enforce coding standards and best practices
- Identify unused variables and dead code
- Validate method signatures and return types

### Laravel Application Analysis

Comprehensive Laravel project analysis:

- Analyze models, controllers, and services
- Validate Eloquent relationships and queries
- Check facade usage and dependency injection
- Verify route parameter types and middleware

### Library Development

Ensure library code quality:

- Validate public API consistency
- Check backward compatibility
- Analyze package dependencies
- Verify documentation alignment with code

### Legacy Code Improvement

Gradually improve existing codebases:

- Identify areas needing refactoring
- Find potential security vulnerabilities
- Detect performance bottlenecks
- Plan migration strategies

## Requirements

- **PHPStan Installation**: PHPStan must be installed via Composer
- **Configuration File**: `phpstan.neon` or `phpstan.neon.dist` must exist
- **PHP Compatibility**: Code must be compatible with specified PHP version
- **Composer Setup**: Valid `composer.json` with PHPStan dependency

## Troubleshooting

### Common Issues

**PHPStan Not Found**

- Ensure PHPStan is installed via Composer: `composer require --dev phpstan/phpstan`
- Verify the vendor directory exists and contains PHPStan
- Check if Composer dependencies were installed correctly

**Configuration File Missing**

- Create a `phpstan.neon` configuration file in the project root
- Ensure the configuration file syntax is valid
- Check if the configuration includes necessary extensions

**Memory Limit Exceeded**

- Increase PHP memory limit in the analysis configuration
- Consider reducing analysis scope or level
- Use PHPStan's `--memory-limit` option in configuration

**Laravel Analysis Issues**

- Install Larastan for Laravel-specific analysis: `composer require --dev nunomaduro/larastan`
- Ensure proper environment configuration
- Verify Laravel encryption key is correct

**Path Not Found**

- Verify specified paths exist in the repository
- Check path formatting (should be relative to working directory)
- Ensure directories contain PHP files to analyze

### Performance Optimization

**Analysis Level Adjustment**

```neon
# Start with lower levels for large codebases
parameters:
    level: 0  # Start here and gradually increase
```

**Memory Configuration**

```neon
parameters:
    tmpDir: var/cache/phpstan
    memory_limit: 1G
```

**Parallel Processing**

```neon
parameters:
    parallel:
        jobSize: 20
        processTimeout: 300.0
```

### Custom PHPStan Configuration

**Exclude Problematic Files**

```neon
parameters:
    excludePaths:
        - tests/
        - database/migrations/
        - storage/
```

**Custom Rules**

```neon
parameters:
    ignoreErrors:
        - '#Call to an undefined method#'
        - message: '#Access to an undefined property#'
          path: app/Models/*
```

## Integration Patterns

### Pre-commit Analysis

```yaml
name: Code Quality
on:
  pull_request:
    branches: [main]

jobs:
  static-analysis:
    uses: ./.github/workflows/php_staticAnalysis.yml
    with:
      phpVersion: "8.3"
```

### Multi-Version Testing

```yaml
jobs:
  static-analysis:
    strategy:
      matrix:
        php-version: ['8.1', '8.2', '8.3']
    uses: ./.github/workflows/php_staticAnalysis.yml
    with:
      phpVersion: ${{ matrix.php-version }}
```

## Related Workflows

- **php_test.yml**: For runtime testing after static analysis