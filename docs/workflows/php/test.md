# PHP Unit Tests Workflow

## Overview

The `php_test.yml` workflow runs PHP unit tests using Pest or PHPUnit with optional CodeCov integration. It supports Laravel projects with environment decryption and
provides flexible configuration options.

## Language/Tool Support

- **PHP**: All PHP projects with Pest or PHPUnit testing frameworks
- **Laravel**: Full support including environment decryption

## Features

- **Flexible Testing**: Supports both Pest and PHPUnit via composer scripts
- **CodeCov Integration**: Optional code coverage reporting
- **Laravel Support**: Environment decryption for Laravel projects
- **Directory Support**: Can target subdirectories within repositories

## Usage

```yaml
uses: ./.github/workflows/php_test.yml
with:
  branch: "main"
  phpVersion: "8.3"
  useCodeCov: true
  codeCovSlug: "your-org/your-repo"
secrets:
  laravelEnvDecryptionKey: ${{ secrets.LARAVEL_ENV_DECRYPTION_KEY }}
  codeCovToken: ${{ secrets.CODECOV_TOKEN }}
```

## Inputs

| Input                        | Type    | Required | Default                   | Description                          |
|------------------------------|---------|----------|---------------------------|--------------------------------------|
| `branch`                     | string  | ❌        | `main`                    | The branch to test                   |
| `phpVersion`                 | string  | ❌        | `8.3`                     | PHP version to use                   |
| `useLaravelEnvDecryptionKey` | boolean | ❌        | `false`                   | Use Laravel environment decryption   |
| `directory`                  | string  | ❌        | `${{ github.workspace }}` | Directory path relative to workspace |
| `useCodeCov`                 | boolean | ❌        | `false`                   | Enable CodeCov integration           |
| `codeCovSlug`                | string  | ❌        | `''`                      | CodeCov repository slug              |

## Secrets

| Secret                    | Required | Description                                       |
|---------------------------|----------|---------------------------------------------------|
| `laravelEnvDecryptionKey` | ❌        | Laravel environment decryption key                |
| `codeCovToken`            | ❌        | CodeCov authentication token                      |

## Required Setup

### Composer Test Script

**This is required in v3.0+**. Your `composer.json` must include a test script:

```json
{
  "scripts": {
    "test": "pest"
  }
}
```

Or for PHPUnit:

```json
{
  "scripts": {
    "test": "phpunit"
  }
}
```

For CodeCov integration:

```json
{
  "scripts": {
    "test": "pest --parallel --coverage-clover coverage.xml --log-junit junit.xml"
  }
}
```

## Workflow Behavior

The workflow has two jobs that run conditionally:

1. **Test Job**: Runs when `useCodeCov: false` - Standard testing without coverage
2. **CodeCov Job**: Runs when `useCodeCov: true` - Testing with coverage reporting

## Example Configurations

### Basic PHP Testing

```yaml
uses: ./.github/workflows/php_test.yml
with:
  phpVersion: "8.3"
```

### Laravel with Environment Decryption

```yaml
uses: ./.github/workflows/php_test.yml
with:
  phpVersion: "8.3"
  useLaravelEnvDecryptionKey: true
secrets:
  laravelEnvDecryptionKey: ${{ secrets.LARAVEL_ENV_DECRYPTION_KEY }}
```

### With CodeCov Integration

```yaml
uses: ./.github/workflows/php_test.yml
with:
  phpVersion: "8.3"
  useCodeCov: true
  codeCovSlug: "your-org/your-repo"
secrets:
  codeCovToken: ${{ secrets.CODECOV_TOKEN }}
```

### Subdirectory Testing

```yaml
uses: ./.github/workflows/php_test.yml
with:
  directory: "backend"
  phpVersion: "8.3"
```

## Laravel Environment Decryption

When `useLaravelEnvDecryptionKey: true`, the workflow will:

1. Decrypt the testing environment file using the provided key
2. Move the decrypted file to `.env` for testing

Requires:

- Laravel project with encrypted environment files
- `LARAVEL_ENV_DECRYPTION_KEY` secret configured

## CodeCov Integration

When `useCodeCov: true`, the workflow will:

1. Run tests with coverage reporting
2. Upload coverage data to CodeCov
3. Upload test results to CodeCov

Requires:

- `codeCovToken` secret
- `codeCovSlug` input with repository identifier
- Test script configured to generate `coverage.xml` and `junit.xml`

## Migration from v2 to v3

### Removed Parameters

- `enforceCoverage` - No longer configurable
- `minCodeCoverage` - No longer configurable
- `runParallel` - No longer configurable

### Required Changes

1. **Add composer test script** - Required for all projects
2. **Update workflow calls** - Remove deprecated parameters
3. **Configure coverage in test script** - Instead of workflow parameters

## Troubleshooting

### Common Issues

**Test Script Not Found**

- Ensure `composer.json` has a `test` script defined
- Verify the test command is correct for your testing framework

**Laravel Environment Issues**

- Ensure the decryption key is correct
- Check that encrypted environment files exist

**CodeCov Upload Fails**

- Verify `codeCovToken` is valid
- Check `codeCovSlug` format (should be "owner/repo")
- Ensure test script generates required coverage files

## Related Workflows

- **php_staticAnalysis.yml**: For code quality analysis
- **actions/php/codeStyle**: For code style enforcement