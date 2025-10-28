#!/bin/bash

run_action() {
    echo "Starting Code Style Enforcement"
    echo "DEBUG: Repository - $GH_REPO"
    echo "DEBUG: Branch - $GH_BRANCH"
    echo "DEBUG: Directory - $GH_DIRECTORY"

    # Check if GH_DIRECTORY is different from GITHUB_WORKSPACE
    if [ "$GITHUB_WORKSPACE" != "$GH_DIRECTORY" ]; then
        echo "Directory is different from workspace, changing to $GH_DIRECTORY"
        GH_DIRECTORY="$GITHUB_WORKSPACE/$GH_DIRECTORY"
        # Ensure directory exists
        if [ ! -d "$GH_DIRECTORY" ]; then
            echo "ERROR: Directory $GH_DIRECTORY does not exist"
            exit 1
        fi
        cd "$GH_DIRECTORY"
    fi

    git config --global user.name "EncoreBot"
    git config --global user.email "ghbot@encoredigitalgroup.com"

    composer install --no-ansi --no-interaction --no-progress --prefer-dist --ignore-platform-reqs
    rm -f "$GH_DIRECTORY"/auth.json

    rector_run
    duster_run
}


#region rector
rector_run() {
    echo "Running Rector"

    echo "GH_DIRECTORY: $GH_DIRECTORY"

    if [ "$GH_ONLY_DIFF" = "true" ]; then
        # Get changed PHP files and run rector only on them
        changed_files=$(git diff --name-only --diff-filter=ACMR origin/main...HEAD | grep '\.php$' || true)

        if [ -n "$changed_files" ]; then
            echo "Running Rector on changed files"
            "$GH_DIRECTORY"/vendor/bin/rector process "$changed_files"
        else
            echo "No PHP files changed, skipping Rector"
        fi
    else
        # Run rector on all files (original behavior)
        echo "Running Rector on all files"
        "$GH_DIRECTORY"/vendor/bin/rector process
    fi

    rector_auto_commit_changes
}

rector_auto_commit_changes() {
    cd "$GITHUB_WORKSPACE"

    if [ -z "$(git status --porcelain)" ]; then
      # Working directory clean
      echo "Working Tree is Clean! Nothing to commit."
    else
      # Add all changes to staging
      git add .

      # Commit changes
      commit_message="Rectifying"
      git commit -m "$commit_message"

      rector_ignore_commit

      # Push changes to origin
      git push origin --force
      # Uncommitted changes
    fi
}

rector_ignore_commit() {
    cd "$GITHUB_WORKSPACE"
    # Get the most recent commit hash
    latest_commit=$(git rev-parse HEAD)

    # Check if the commit hash was retrieved successfully
    if [ -z "$latest_commit" ]; then
      echo "Error: Could not retrieve the latest commit hash."
      exit 1
    fi

    # Define the ignore file path
    ignore_file=".git-blame-ignore-revs"

    # Append the latest commit hash to the file
    echo "$latest_commit" >> "$ignore_file"

    # Confirm the operation
    echo "Commit hash $latest_commit has been added to $ignore_file."

    rector_commit_blame_file
}

rector_commit_blame_file() {
    git add .
    git commit -m "Ignore Rector Commit in Git Blame"
}
#endregion

#region duster
duster_run() {
    echo "Running Duster"

    echo "GH_DIRECTORY: $GH_DIRECTORY"

    if [ "$GH_ONLY_DIFF" = "true" ]; then
        # Run duster only on changed files
        echo "Running Duster on changed files only"
        "$GH_DIRECTORY"/vendor/bin/duster fix --diff=main
    else
        # Run duster on all files (original behavior)
        echo "Running Duster on all files"
        "$GH_DIRECTORY"/vendor/bin/duster fix
    fi

    duster_auto_commit_changes
}

duster_auto_commit_changes() {
    cd "$GITHUB_WORKSPACE"

    if [ -z "$(git status --porcelain)" ]; then
      # Working directory clean
      echo "Working Tree is Clean! Nothing to commit."
    else
      # Add all changes to staging
      git add .

      # Commit changes
      commit_message="Dusting"
      git commit -m "$commit_message"

      duster_ignore_commit

      # Push changes to origin
      git push origin --force
      # Uncommitted changes
    fi
}

duster_ignore_commit() {
    cd "$GITHUB_WORKSPACE"
    # Get the most recent commit hash
    latest_commit=$(git rev-parse HEAD)

    # Check if the commit hash was retrieved successfully
    if [ -z "$latest_commit" ]; then
      echo "Error: Could not retrieve the latest commit hash."
      exit 1
    fi

    # Define the ignore file path
    ignore_file=".git-blame-ignore-revs"

    # Append the latest commit hash to the file
    echo "$latest_commit" >> "$ignore_file"

    # Confirm the operation
    echo "Commit hash $latest_commit has been added to $ignore_file."

    duster_commit_blame_file
}

duster_commit_blame_file() {
    git add .
    git commit -m "Ignore Duster Commit in Git Blame"
}
#endregion

run_action
