# Dependabot Auto-Merge Workflow

This GitHub Actions workflow automates the process of approving and merging Dependabot pull requests for semver-minor and semver-patch updates.

## Workflow File

The workflow file is located at `.github/workflows/dependabotAutoMerge.yml`.

## Usage

To use this workflow, you can call it from another workflow or trigger it manually. Below is an example of how to call this workflow:

```yaml
name: Pull Request

on:
    pull_request_target:

jobs:
    # The Rest of Your PR Workflow

    AutoMerge:
        needs: [ StaticAnalysis, Test ] # Encore Digital Group runs this step only after Unit Tests and Static Analysis have passed.
        name: AutoMerge
        uses: EncoreDigitalGroup/.github/.github/workflows/dependabotAutoMerge.yml@v1
```

This workflow will approve and merge Dependabot pull requests based on the provided conditions and secrets. Adjust the conditions and steps as needed to fit your process.