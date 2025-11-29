# PHP Git Status Check Action

## Overview

The `actions/php/gitStatusCheck` action determines whether PHP-specific files have been modified by checking git status for PHP source files (*.php),
Composer files (composer.json, composer.lock), and Node.js files (package.json, package-lock.json). This action is ideal for PHP projects that may also include frontend
assets, enabling conditional workflow execution.

## Language/Tool Support

- **PHP**: PHP source files and Composer dependency management
- **Composer**: composer.json and composer.lock files
- **Node.js**: package.json and package-lock.json files (for mixed PHP/JS projects)
- **Git**: Git repository status checking
- **Shell Scripts**: Bash-based implementation

## Features

- **Multi-Language Support**: Monitors PHP, Composer, and Node.js files
- **Branch-based Checking**: Configurable branch for comparison
- **Output Generation**: Provides boolean output for workflow decisions
- **Laravel-Friendly**: Supports Laravel projects with mixed assets
- **Modern PHP Projects**: Handles composer and npm dependencies

## Usage

```yaml
- name: Check php File Changes
  id: php-check
  uses: ./actions/php/gitStatusCheck
  with:
    branch: main

- name: Run php Tests
  if: steps.php-check.outputs.shouldRun == 'true'
  run: composer test
```

## Inputs

| Input    | Type   | Required | Default | Description                       |
|----------|--------|----------|---------|-----------------------------------|
| `branch` | string | ❌        | `main`  | Branch to compare changes against |

## Outputs

| Output      | Type    | Description                                          |
|-------------|---------|------------------------------------------------------|
| `shouldRun` | boolean | Whether PHP-related files were modified (true/false) |

## Monitored File Types

The action checks for changes in:

### PHP Files

- `*.php` - All PHP source files
- Application code, configuration, migrations
- Test files and development tools

### Composer Files

- `composer.json` - Package dependencies and project metadata
- `composer.lock` - Locked dependency versions

### Node.js Files (for mixed projects)

- `package.json` - Frontend dependencies and scripts
- `package-lock.json` - Locked frontend dependency versions

### File Examples

```
index.php
app/Models/User.php
config/app.php
database/migrations/create_users_table.php
tests/Feature/UserTest.php
composer.json
composer.lock
package.json (for frontend assets)
package-lock.json (for frontend assets)
```

## Usage Examples

### PHP Laravel Project

```yaml
name: Laravel CI Pipeline
on:
  pull_request:
    branches: [main]

jobs:
  check-php-changes:
    runs-on: ubuntu-latest
    outputs:
      php-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check for php Changes
        id: check
        uses: ./actions/php/gitStatusCheck
        with:
          branch: main

  php-tests:
    needs: check-php-changes
    if: needs.check-php-changes.outputs.php-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup php
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.3'

      - name: Install Dependencies
        run: composer install

      - name: Run Tests
        run: composer test

  static-analysis:
    needs: check-php-changes
    if: needs.check-php-changes.outputs.php-changed == 'true'
    uses: ./.github/workflows/php_staticAnalysis.yml
    with:
      phpVersion: '8.3'
```

### Composer Dependency Validation

```yaml
name: Dependency Security Check
on:
  pull_request:
    paths:
      - 'composer.json'
      - 'composer.lock'

jobs:
  check-composer:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check Composer Changes
        id: composer-check
        uses: ./actions/php/gitStatusCheck
        with:
          branch: ${{ github.base_ref }}

      - name: Security Audit
        if: steps.composer-check.outputs.shouldRun == 'true'
        run: composer audit

      - name: Validate Dependencies
        if: steps.composer-check.outputs.shouldRun == 'true'
        run: composer validate --strict
```

### Mixed PHP/JavaScript Project

```yaml
name: Full Stack Validation
on:
  push:
    branches: [main, develop]

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      backend-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Backend Changes
        id: check
        uses: ./actions/php/gitStatusCheck

  backend-pipeline:
    needs: detect-changes
    if: needs.detect-changes.outputs.backend-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Setup php Environment
        uses: shivammathur/setup-php@v2
        with:
          php-version: '8.3'

      - name: Run php Tests
        run: composer test

      - name: Build Assets (if package.json changed)
        run: |
          if [ -f package.json ]; then
            npm ci
            npm run build
          fi
```

### WordPress Plugin Development

```yaml
name: WordPress Plugin CI
on:
  pull_request:

jobs:
  check-plugin-changes:
    runs-on: ubuntu-latest
    outputs:
      plugin-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Plugin Changes
        id: check
        uses: ./actions/php/gitStatusCheck

  wordpress-tests:
    needs: check-plugin-changes
    if: needs.check-plugin-changes.outputs.plugin-changed == 'true'
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:5.7
        env:
          MYSQL_ROOT_PASSWORD: password
    steps:
      - name: Setup WordPress Test Environment
        run: |
          # WordPress test setup
          composer install
          npm install
          npm run build

      - name: Run Plugin Tests
        run: composer test
```

## Integration Patterns

### Performance-Optimized Pipeline

```yaml
name: Optimized php Pipeline
on:
  pull_request:

jobs:
  quick-check:
    runs-on: ubuntu-latest
    outputs:
      php-changed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check php Changes
        id: check
        uses: ./actions/php/gitStatusCheck

  fast-validation:
    needs: quick-check
    if: needs.quick-check.outputs.php-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: php Syntax Check
        run: find . -name "*.php" -exec php -l {} \;

  comprehensive-testing:
    needs: [quick-check, fast-validation]
    if: needs.quick-check.outputs.php-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Full Test Suite
        run: |
          composer install
          composer test:unit
          composer test:integration
```

### Multi-Environment Deployment

```yaml
name: Multi-Environment php Deployment
on:
  push:
    branches: [main, staging, develop]

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      deploy-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Deployment Changes
        id: check
        uses: ./actions/php/gitStatusCheck

  deploy-staging:
    needs: detect-changes
    if: github.ref == 'refs/heads/staging' && needs.detect-changes.outputs.deploy-needed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Staging
        run: echo "Deploying to staging environment..."

  deploy-production:
    needs: detect-changes
    if: github.ref == 'refs/heads/main' && needs.detect-changes.outputs.deploy-needed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Production
        run: echo "Deploying to production environment..."
```

### Code Quality Pipeline

```yaml
name: php Code Quality
on:
  pull_request:

jobs:
  detect-php-changes:
    runs-on: ubuntu-latest
    outputs:
      quality-check-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check php Code Changes
        id: check
        uses: ./actions/php/gitStatusCheck

  code-style:
    needs: detect-php-changes
    if: needs.detect-php-changes.outputs.quality-check-needed == 'true'
    uses: ./.github/workflows/php_dusterFix.yml

  static-analysis:
    needs: detect-php-changes
    if: needs.detect-php-changes.outputs.quality-check-needed == 'true'
    uses: ./.github/workflows/php_staticAnalysis.yml
    with:
      phpVersion: '8.3'
```

## PHP-Specific Use Cases

### Laravel Application Testing

```yaml
name: Laravel Testing
on:
  pull_request:

jobs:
  check-laravel-changes:
    runs-on: ubuntu-latest
    outputs:
      test-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check Laravel Changes
        id: check
        uses: ./actions/php/gitStatusCheck

  laravel-tests:
    needs: check-laravel-changes
    if: needs.check-laravel-changes.outputs.test-needed == 'true'
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: password
          MYSQL_DATABASE: testing
    steps:
      - name: Setup Laravel Environment
        run: |
          cp .env.testing .env
          php artisan key:generate
          php artisan migrate

      - name: Run Feature Tests
        run: php artisan test
```

### Composer Security Monitoring

```yaml
name: Security Monitoring
on:
  schedule:
    - cron: '0 2 * * 1'  # Weekly on Monday

jobs:
  check-dependencies:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check for Dependency Changes
        id: deps
        uses: ./actions/php/gitStatusCheck

      - name: Security Audit
        if: steps.deps.outputs.shouldRun == 'true'
        run: |
          composer audit
          composer outdated --direct
```

### Multi-Version PHP Testing

```yaml
name: php Version Compatibility
on:
  pull_request:

jobs:
  check-php-changes:
    runs-on: ubuntu-latest
    outputs:
      test-needed: ${{ steps.check.outputs.shouldRun }}
    steps:
      - uses: actions/checkout@v4

      - name: Check php Changes
        id: check
        uses: ./actions/php/gitStatusCheck

  compatibility-test:
    needs: check-php-changes
    if: needs.check-php-changes.outputs.test-needed == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        php-version: ['8.1', '8.2', '8.3']
    steps:
      - name: Setup php ${{ matrix.php-version }}
        uses: shivammathur/setup-php@v2
        with:
          php-version: ${{ matrix.php-version }}

      - name: Test Compatibility
        run: composer test
```

## File Detection Patterns

### Standard PHP Patterns

- Application files: `app/**/*.php`
- Configuration: `config/**/*.php`
- Database: `database/**/*.php`
- Tests: `tests/**/*.php`

### Framework-Specific Patterns

#### Laravel

- Models: `app/Models/*.php`
- Controllers: `app/Http/Controllers/*.php`
- Migrations: `database/migrations/*.php`
- Seeders: `database/seeders/*.php`

#### Symfony

- Controllers: `src/Controller/*.php`
- Entities: `src/Entity/*.php`
- Services: `src/Service/*.php`

#### WordPress

- Plugin files: `*.php`
- Theme files: `*.php`
- Custom post types: `inc/*.php`

## Troubleshooting

### Common Issues

**Always Returns False**

- Verify PHP files exist in repository
- Check that composer.json exists for PHP projects
- Ensure checkout action runs before this action

**Missing Dependencies**

- Verify composer.json and composer.lock are committed
- Check package.json exists if using Node.js assets
- Ensure proper project structure

**File Detection Issues**

- Test file patterns manually: `find . -name "*.php"`
- Check git status: `git status --porcelain`
- Verify branch comparison works

### Debugging Commands

```bash
# Check php files
find . -name "*.php" -o -name "composer.json" -o -name "composer.lock"

# Check Node.js files
find . -name "package.json" -o -name "package-lock.json"

# Check git differences
git diff --name-only origin/main | grep -E '\.(php|json)$'
```

## Best Practices

### Project Structure

1. **Standard Layout**: Follow PHP framework conventions
2. **Dependency Management**: Keep composer.lock committed
3. **Mixed Projects**: Organize PHP and frontend code clearly

### Workflow Design

1. **Early Checking**: Run checks early in pipeline
2. **Conditional Logic**: Use proper conditional execution

## Migration Guide

### From Path-Based Filters

Replace path-based workflow triggers:

```yaml
# Before: Path-based filtering
on:
  pull_request:
    paths:
      - '**/*.php'
      - 'composer.json'
      - 'composer.lock'

# After: Dynamic checking
on:
  pull_request:

jobs:
  check:
    steps:
      - uses: ./actions/php/gitStatusCheck
        id: php-check
```

### From Manual Commands

Replace manual file checking:

```yaml
# Before: Manual commands
- name: Check php Changes
  run: |
    if git diff --name-only origin/main | grep -E '\.(php|json)$'; then
      echo "php_changed=true" >> $GITHUB_OUTPUT
    fi

# After: Use action
- name: Check php Changes
  id: check
  uses: ./actions/php/gitStatusCheck
```

### Framework-Specific Integration

#### Laravel Projects

```yaml
- uses: ./actions/php/gitStatusCheck
  id: laravel-check

- name: Laravel Tests
  if: steps.laravel-check.outputs.shouldRun == 'true'
  run: php artisan test
```

#### Symfony Projects

```yaml
- uses: ./actions/php/gitStatusCheck
  id: symfony-check

- name: Symfony Tests
  if: steps.symfony-check.outputs.shouldRun == 'true'
  run: php bin/phpunit
```