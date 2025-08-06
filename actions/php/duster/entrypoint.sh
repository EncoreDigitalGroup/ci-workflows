#!/bin/bash

run_action() {
    echo "Starting Code Style Enforcement"
    echo "DEBUG: Repository - $GH_REPO"
    echo "DEBUG: Branch - $GH_BRANCH"

    git config --global user.name "EncoreBot"
    git config --global user.email "ghbot@encoredigitalgroup.com"

    composer install --no-ansi --no-interaction --no-progress --prefer-dist --ignore-platform-reqs

    duster_run
}

#region duster
duster_run() {
    echo "Running Duster"

    "$GITHUB_WORKSPACE"/vendor/bin/duster fix

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