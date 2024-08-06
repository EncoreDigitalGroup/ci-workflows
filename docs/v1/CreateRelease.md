# Create Release Workflow

This GitHub Actions workflow automates the process of creating a release in your repository. It allows you to configure whether the release is a pre-release, a draft, and
whether to generate release notes automatically.

## Workflow File

The workflow file is located at `.github/workflows/createRelease.yml`.

## Inputs

- **preRelease**:
    - **Type**: boolean
    - **Description**: Indicates if the release is a pre-release.
    - **Required**: No
    - **Default**: `false`

- **generateReleaseNotes**:
    - **Type**: boolean
    - **Description**: Indicates if release notes should be generated automatically.
    - **Required**: No
    - **Default**: `true`

- **isDraft**:
    - **Type**: boolean
    - **Description**: Indicates if the release is a draft.
    - **Required**: No
    - **Default**: `false`

## Secrets

- **token**:
    - **Description**: GitHub token used for authentication.
    - **Required**: Yes

## Usage

To use this workflow, you can call it from another workflow. Below is an example of how to call this workflow:

```yaml
name: "Create Release"

on:
    push:
        tags:
            - 'v[0-9]+.[0-9]+.[0-9]+'
            - '!v[0-9]+.[0-9]+.[0-9]+-rc[0-9]+'
            - '!v[0-9]'

permissions:
    pull-requests: write
    contents: write

jobs:
    CreateRelease:
        name: Create Release
        uses: EncoreDigitalGroup/ci-workflows/.github/workflows/createRelease.yml
        with:
            generateReleaseNotes: true
            isDraft: false
        secrets:
            token: ${{ secrets.GITHUB_TOKEN }}
```

This workflow will create a release based on the provided inputs and secrets. Adjust the inputs as needed to fit your release process.