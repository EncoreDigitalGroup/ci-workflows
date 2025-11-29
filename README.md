# CI Workflows

Encore Digital Group has written a number of reusable GitHub Actions and Workflows for common CI/CD tasks. These actions and workflows are designed to be used across
multiple projects to standardize and streamline development processes. We decided to open-source a large number of our GitHub Actions and Workflows for other developers
and organizations to use.

## What Workflows and Actions Are Included

### Workflows

| Workflow Name                | Language/Tool | Description                                                            | Link                                                            |
|------------------------------|---------------|------------------------------------------------------------------------|-----------------------------------------------------------------|
| GitHub Dependabot Auto Merge | GitHub        | Automatically merges Dependabot pull requests                          | [Documentation](docs/workflows/github/dependabot-auto-merge.md) |
| GitHub Directory Import      | GitHub        | Imports directory contents from external repositories                  | [Documentation](docs/workflows/github/directory-import.md)      |
| GitHub Directory Sync        | GitHub        | Synchronizes directory contents between repositories                   | [Documentation](docs/workflows/github/directory-sync.md)        |
| GitHub Sync Branch           | GitHub        | Synchronizes branches between repositories                             | [Documentation](docs/workflows/github/sync-branch.md)           |
| GitHub Update Changelog      | GitHub        | Updates changelog files automatically                                  | [Documentation](docs/workflows/github/update-changelog.md)      |
| Go Static Analysis           | Go            | Runs static analysis tools for Go projects                             | [Documentation](docs/workflows/go/static-analysis.md)           |
| PHP Static Analysis          | PHP           | Runs static analysis tools (PHPStan, etc.) for PHP projects            | [Documentation](docs/workflows/php/static-analysis.md)          |
| PHP Unit Tests               | PHP           | Runs PHP unit tests with Pest/PHPUnit and optional CodeCov integration | [Documentation](docs/workflows/php/test.md)                     |

### Actions

| Action Name                 | Language/Tool | Description                                                     | Link                                                                |
|-----------------------------|---------------|-----------------------------------------------------------------|---------------------------------------------------------------------|
| GitHub Create Release       | GitHub        | Creates GitHub releases with automated changelog generation     | [Documentation](./docs/actions/github/create-release.md)            |
| GitHub Format PR Title      | GitHub        | Formats pull request titles according to conventions            | [Documentation](./docs/actions/github/format-pull-request-title.md) |
| GitHub Git Status Check     | GitHub        | Validates git repository status for GitHub workflows            | [Documentation](./docs/actions/github/git-status-check.md)          |
| GitHub JSON Diff Alert      | GitHub        | Compares JSON files and posts differences as PR comments        | [Documentation](./docs/actions/github/json-diff-alert.md)           |
| Go Git Status Check         | Go            | Validates git repository status for Go projects                 | [Documentation](./docs/actions/go/git-status-check.md)              |
| PHP Code Style              | PHP           | Applies code style fixes using Rector and Duster (Laravel Pint) | [Documentation](./docs/actions/php/code-style.md)                   |
| PHP Duster                  | PHP           | Applies Laravel Pint code style fixes to PHP projects           | [Documentation](./docs/actions/php/duster.md)                       |
| PHP Git Status Check        | PHP           | Validates git repository status for PHP projects                | [Documentation](./docs/actions/php/git-status-check.md)             |
| PHP Rector                  | PHP           | Applies automated refactoring using Rector for PHP projects     | [Documentation](./docs/actions/php/rector.md)                       |
| TypeScript Git Status Check | TypeScript    | Validates git repository status for TypeScript projects         | [Documentation](./docs/actions/typescript/git-status-check.md)      |

## Planned Changes

You can review the planned changes for major versions of CI Workflows:

- [v3.0 Planned Changes](docs/planned-changes/v3.md)
- [v4.0 Planned Changes](docs/planned-changes/v4.md)