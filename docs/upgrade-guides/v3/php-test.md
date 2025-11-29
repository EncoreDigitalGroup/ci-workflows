---
title: PHP Test
---

# PHP Test Workflow v3 Upgrade Guide

## Overview

The PHP test workflow has been significantly simplified in v3, removing several configuration options and consolidating the testing approach.

## Breaking Changes

### Removed Input Parameters

The following input parameters have been **removed** in v3:

- `enforceCoverage` - Coverage enforcement is no longer configurable
- `minCodeCoverage` - Minimum coverage percentage is no longer configurable
- `runParallel` - Parallel execution is no longer configurable

### Changed Input Parameters

- `phpVersion`: Default changed from `'8.2'` to `'8.3'`

### Removed Job

The separate `TestWithCodeCov` job has been **removed**. CodeCov functionality is now integrated into the main `Test` job.

## Migration Steps

### 1. Update Workflow Calls

**Before (v2):**
```yaml
uses: ./.github/workflows/php_test.yml
with:
  enforceCoverage: true
  minCodeCoverage: 85
  runParallel: true
  phpVersion: '8.2'
```

**After (v3):**
```yaml
uses: ./.github/workflows/php_test.yml
with:
  phpVersion: '8.3'  # Optional, defaults to 8.3
```

### 2. Update Test Command

The workflow now uses `composer run test` instead of directly calling Pest with various flags. Ensure your `composer.json` has a test script defined:

```json
{
  "scripts": {
    "test": "pest --parallel --coverage-clover coverage.xml --log-junit junit.xml"
  }
}
```

### 3. CodeCov Integration

CodeCov uploads are now conditional based on the `useCodeCov` input parameter:

- Coverage upload only happens when `useCodeCov: true`
- Test results upload only happens when `useCodeCov: true`

## Key Improvements

1. **Simplified Configuration**: Fewer input parameters to manage
2. **Consolidated Testing**: Single job handles both regular testing and CodeCov
3. **Composer-based Testing**: Uses `composer run test` for better project-specific configuration
4. **Updated PHP Default**: Uses PHP 8.3 by default
5. **Conditional CodeCov**: CodeCov uploads only occur when explicitly enabled

## Migration Checklist

- [ ] Remove `enforceCoverage`, `minCodeCoverage`, and `runParallel` parameters from workflow calls
- [ ] Update `phpVersion` to `'8.3'` or remove to use default
- [ ] Add `test` script to `composer.json` with desired Pest or PHPUnit configuration
- [ ] Ensure `useCodeCov: true` is set if CodeCov integration is needed
- [ ] Test the updated workflow to ensure it works as expected