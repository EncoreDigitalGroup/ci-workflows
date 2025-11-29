---
title: PHP Unit
---

# PHP Test PHPUnit Workflow v3 Upgrade Guide

## Overview

The `php_test_phpunit.yml` workflow has been **removed** in v3. Projects using PHPUnit should migrate to the consolidated `php_test.yml` workflow.

## Breaking Changes

### Workflow Removal

The dedicated PHPUnit workflow (`php_test_phpunit.yml`) has been completely removed in v3. All PHP testing is now handled by the unified `php_test.yml` workflow.

## Migration Steps

### 1. Update Workflow References

**Before (v2):**
```yaml
uses: ./.github/workflows/php_test_phpunit.yml
with:
  phpVersion: '8.3'
```

**After (v3):**
```yaml
uses: ./.github/workflows/php_test.yml
with:
  phpVersion: '8.3'
```

### 2. Create Composer Test Script

The new unified workflow uses `composer run test` instead of directly calling PHPUnit. You **must** add a test script to your `composer.json`:

```json
{
  "scripts": {
    "test": "phpunit"
  }
}
```

For projects that need coverage reporting with CodeCov, use:

```json
{
  "scripts": {
    "test": "phpunit --coverage-clover coverage.xml --log-junit junit.xml"
  }
}
```

### 3. Configure PHPUnit Options

Since the workflow no longer directly calls PHPUnit with specific flags, configure your testing options via composer script flags:

#### Composer Script with Flags

```json
{
  "scripts": {
    "test": "phpunit --coverage-clover coverage.xml --log-junit junit.xml --testdox"
  }
}
```

### 4. CodeCov Integration Added

You now have the option to send your test results to Sentry CodeCov. To do so, ensure the `useCodeCov` parameter is set to `true`.

```yaml
uses: ./.github/workflows/php_test.yml
with:
  useCodeCov: true
  codeCovSlug: your-org/your-repo
secrets:
  codeCovToken: ${{ secrets.CODECOV_TOKEN }}
```

## Key Differences

| Feature       | php_test_phpunit.yml (v2) | php_test.yml (v3)                    |
|---------------|---------------------------|--------------------------------------|
| Test Command  | `./vendor/bin/phpunit`    | `composer run test`                  |
| Configuration | Workflow parameters       | Composer scripts + PHPUnit config    |
| CodeCov       | Not supported             | Conditional support via `useCodeCov` |
| Flexibility   | Limited                   | High (via composer scripts)          |

## Migration Checklist

- [ ] Replace workflow references from `php_test_phpunit.yml` to `php_test.yml`
- [ ] Add `test` script to `composer.json` with appropriate PHPUnit command
- [ ] Configure PHPUnit options via `phpunit.xml` or composer script flags
- [ ] If using CodeCov, set `useCodeCov: true` and provide required secrets
- [ ] Test the updated workflow to ensure PHPUnit runs correctly
- [ ] Verify coverage and test result uploads work (if using CodeCov)

## Benefits of Migration

1. **Unified Testing**: Single workflow for all PHP testing frameworks
2. **Project-specific Configuration**: Use composer scripts for custom test commands
3. **Better Flexibility**: Configure testing behavior at the project level
4. **CodeCov Integration**: Optional but integrated coverage reporting
5. **Reduced Maintenance**: Fewer workflow files to maintain