# Directory Import Workflow

This GitHub Actions workflow automates the process of importing a directory from a source repository to a target repository.

## Workflow File

The workflow file is located at `.github/workflows/directoryImport.yml`.

## Inputs

- **sourceRepository**:
    - **Type**: string
    - **Description**: The repository from which the directory will be imported.
    - **Required**: Yes

- **sourceDirectory**:
    - **Type**: string
    - **Description**: The directory in the source repository to be imported.
    - **Required**: Yes

- **targetRepository**:
    - **Type**: string
    - **Description**: The repository to which the directory will be imported.
    - **Required**: Yes

- **targetDirectory**:
    - **Type**: string
    - **Description**: The directory in the target repository where the imported directory will be placed.
    - **Required**: Yes

- **targetDirectoryName**:
    - **Type**: string
    - **Description**: The name of the directory in the target repository. Defaults to `empty_value`.
    - **Required**: No
    - **Default**: `empty_value`

- **commitMessage**:
    - **Type**: string
    - **Description**: The commit message for the import operation. Defaults to `[Automated] DirectoryImport`.
    - **Required**: No
    - **Default**: `[Automated] DirectoryImport`

## Secrets

- **token**:
    - **Description**: GitHub token used for authentication.
    - **Required**: Yes

## Usage

To use this workflow, you can call it from another workflow. Below is an example of how to call this workflow:

```yaml
name: Directory Import

on:
    workflow_dispatch:

jobs:
    ImportDirectory:
        uses: EncoreDigitalGroup/ci-pipelines/.github/workflows/directoryImport.yml@v1
        with:
            sourceRepository: 'sourceOrganization/sourceRepository'
            sourceDirectory: 'sourceDirectory/' #Trailing Slash is Required
            targetRepository: 'targetOrganization/targetRepository' #Trailing Slash is Required
            targetDirectory: 'targetDirectory' #Trailing Slash Should Be Omitted
        secrets:
            token: ${{ secrets.GITHUB_TOKEN }}
```

This workflow will import a directory from the specified source repository to the specified target repository based on the provided inputs and secrets. Adjust the inputs
as needed to fit your import process.