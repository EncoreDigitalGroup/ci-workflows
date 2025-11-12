# CI Workflows

Encore Digital Group has written a number of reusable GitHub Actions and Workflows for common CI/CD tasks. These actions and workflows are designed to be used across
multiple projects to standardize and streamline development processes. We decided to open-source a large number of our GitHub Actions and Workflows for other developers
and organizations to use.

## What Workflows and Actions Are Included

### Workflows

| Workflow Name                | Language/Tool | Description                                                            | Link                                                       |
|------------------------------|---------------|------------------------------------------------------------------------|------------------------------------------------------------|
| GitHub Dependabot Auto Merge | GitHub        | Automatically merges Dependabot pull requests                          | [Documentation](./Workflows/github_dependabotAutoMerge.md) |
| GitHub Directory Import      | GitHub        | Imports directory contents from external repositories                  | [Documentation](./Workflows/github_directoryImport.md)     |
| GitHub Directory Sync        | GitHub        | Synchronizes directory contents between repositories                   | [Documentation](./Workflows/github_directorySync.md)       |
| GitHub Sync Branch           | GitHub        | Synchronizes branches between repositories                             | [Documentation](./Workflows/github_syncBranch.md)          |
| GitHub Update Changelog      | GitHub        | Updates changelog files automatically                                  | [Documentation](./Workflows/github_updateChangelog.md)     |
| Go Static Analysis           | Go            | Runs static analysis tools for Go projects                             | [Documentation](./Workflows/go_staticAnalysis.md)          |
| PHP Static Analysis          | PHP           | Runs static analysis tools (PHPStan, etc.) for PHP projects            | [Documentation](./Workflows/php_staticAnalysis.md)         |
| PHP Unit Tests               | PHP           | Runs PHP unit tests with Pest/PHPUnit and optional CodeCov integration | [Documentation](./Workflows/php_test.md)                   |

### Actions

| Action Name                 | Language/Tool | Description                                                     | Link                                                        |
|-----------------------------|---------------|-----------------------------------------------------------------|-------------------------------------------------------------|
| GitHub Create Release       | GitHub        | Creates GitHub releases with automated changelog generation     | [Documentation](./Actions/GitHub/createRelease.md)          |
| GitHub Format PR Title      | GitHub        | Formats pull request titles according to conventions            | [Documentation](./Actions/GitHub/formatPullRequestTitle.md) |
| GitHub Git Status Check     | GitHub        | Validates git repository status for GitHub workflows            | [Documentation](./Actions/GitHub/gitStatusCheck.md)         |
| GitHub JSON Diff Alert      | GitHub        | Compares JSON files and posts differences as PR comments        | [Documentation](./Actions/GitHub/jsonDiffAlert.md)          |
| Go Git Status Check         | Go            | Validates git repository status for Go projects                 | [Documentation](./Actions/Go/gitStatusCheck.md)             |
| PHP Code Style              | PHP           | Applies code style fixes using Rector and Duster (Laravel Pint) | [Documentation](./Actions/PHP/codeStyle.md)                 |
| PHP Duster                  | PHP           | Applies Laravel Pint code style fixes to PHP projects           | [Documentation](./Actions/PHP/duster.md)                    |
| PHP Git Status Check        | PHP           | Validates git repository status for PHP projects                | [Documentation](./Actions/PHP/gitStatusCheck.md)            |
| PHP Rector                  | PHP           | Applies automated refactoring using Rector for PHP projects     | [Documentation](./Actions/PHP/rector.md)                    |
| TypeScript Git Status Check | TypeScript    | Validates git repository status for TypeScript projects         | [Documentation](./Actions/TypeScript/gitStatusCheck.md)     |

## Planned Changes

You can review the planned changes for major versions of CI Workflows:

- [v3.0 Planned Changes](./PlannedChanges/v3.md)
- [v4.0 Planned Changes](./PlannedChanges/v4.md)