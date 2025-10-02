# CI Pipelines v3 Upgrade Guide

## Overview

Version 3 of the CI Pipelines introduces significant changes aimed at simplifying workflows, consolidating functionality, and improving maintainability. This major
version includes both breaking changes and complete removal of some workflows.

## 📋 Migration Summary

| Workflow                      | Change Type             | Status               | Migration Guide                                            |
|-------------------------------|-------------------------|----------------------|------------------------------------------------------------|
| `php_test.yml`                | 🔄 **Breaking Changes** | Modified             | [php_test.md](./php_test.md)                               |
| `php_test_phpunit.yml`        | ❌ **Removed**           | Deleted              | [php_test_phpunit.md](./php_test_phpunit.md)               |
| `php_dusterFix.yml`           | 🔄 **Replaced**         | Replaced with Action | [php_dusterFix.md](./php_dusterFix.md)                     |

## 🚨 Breaking Changes

### High Impact Changes

- **PHP Test Consolidation**: All PHP testing now uses a single unified workflow
- **Composer Script Requirement**: Tests now require `composer run test` instead of direct tool calls
- **PHP Version Update**: Default PHP version updated from 8.2 to 8.3

### Removed Workflows

- **PHPUnit Workflow**: Consolidated into main PHP test workflow
- **Duster Fix Workflow**: Replaced with enhanced `actions/php/codeStyle` action

## 📚 Detailed Migration Guides

### PHP Testing Changes

#### [php_test.yml - Breaking Changes](./php_test.md)

- ❌ Removed: `enforceCoverage`, `minCodeCoverage`, `runParallel` parameters
- 🔄 Changed: Default PHP version from 8.2 to 8.3
- 🔄 Changed: Now uses `composer run test` command
- ✅ Required: Must configure test script in `composer.json`

#### [php_test_phpunit.yml - Workflow Removed](./php_test_phpunit.md)

- ❌ **Complete Removal**: Workflow no longer exists
- 🔄 **Migration Path**: Use unified `php_test.yml` workflow
- ✅ **Required**: Create `composer.json` test script with PHPUnit command
- 📋 **Action**: Update all workflow references

### Code Style Changes

#### [php_dusterFix.yml - Replaced with Action](./php_dusterFix.md)

- ❌ **Workflow Removed**: No longer available as reusable workflow
- 🔄 **Replacement**: Use `actions/php/codeStyle` action
- ✨ **Enhanced**: Now includes both Duster and Rector
- ✅ **Required**: Update to action-based approach with new parameters

## 🚀 Quick Start Migration

### 1. Assess Impact

```bash
# Find affected workflows in your repository
grep -r "php_test_phpunit.yml\|php_dusterFix.yml\|newRelic_changeTracking.yml" .github/workflows/
```

### 2. Update PHP Testing

Add the following to your `composer.json` file:

```json
{
  "scripts": {
    "test": "pest --parallel --coverage-clover coverage.xml --log-junit junit.xml"
  }
}
```

### 3. Update Workflow References

```yaml
# Before
uses: ./.github/workflows/php_test_phpunit.yml

# After
uses: ./.github/workflows/php_test.yml
```

## 📊 Impact Assessment

### Low Impact Projects

- Projects using only basic `php_test.yml` functionality
- Projects not using removed workflows

### Medium Impact Projects

- Projects using `php_test.yml` with coverage parameters
- Projects using `php_test_phpunit.yml`

### High Impact Projects

- Projects using `php_dusterFix.yml`
- Projects with complex PHP testing configurations

## ✅ Migration Checklist

### Pre-Migration

- [ ] Inventory current workflow usage
- [ ] Review project-specific requirements
- [ ] Plan migration strategy for each affected workflow

### PHP Testing Migration

- [ ] Add `test` script to `composer.json`
- [ ] Update `php_test.yml` workflow calls
- [ ] Remove coverage-related parameters
- [ ] Test new workflow functionality

### Code Style Migration

- [ ] Replace `php_dusterFix.yml` with `actions/php/codeStyle`
- [ ] Update required parameters
- [ ] Ensure Rector dependency is installed
- [ ] Test combined Duster + Rector functionality

### Post-Migration

- [ ] Verify all workflows execute successfully
- [ ] Update team documentation
- [ ] Remove unused workflow references

## 🆘 Getting Help

If you encounter issues during migration:

1. **Review Specific Guides**: Each migration guide contains detailed troubleshooting sections
2. **Check Dependencies**: Ensure all required packages are installed
3. **Validate Syntax**: Use GitHub's workflow syntax validation
4. **Test Incrementally**: Migrate one workflow at a time

## 🎯 Benefits of v3

- **Simplified Configuration**: Fewer parameters to manage
- **Enhanced Functionality**: Better code style tools and testing flexibility
- **Improved Performance**: Better caching and dependency management
- **Reduced Maintenance**: Fewer workflow files to maintain
- **Project Control**: More flexibility through composer scripts and project-level configuration

---

**Important**: This is a major version upgrade with breaking changes. Plan your migration carefully and test thoroughly before deploying to production workflows.