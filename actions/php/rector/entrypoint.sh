#!/bin/bash

run_action() {
    echo "Running Rector"

    composer install --no-ansi --no-interaction --no-progress --prefer-dist --ignore-platform-reqs

    "$GITHUB_WORKSPACE"/vendor/bin/rector process
}


auto_commit_changes() {
    cd "$GITHUB_WORKSPACE"

    git config --global user.name "EncoreBot"
    git config --global user.email "ghbot@encoredigitalgroup.com"

    if [ -z "$(git status --porcelain)" ]; then
      # Working directory clean
      echo "Working Tree is Clean! Nothing to commit."
    else
      # Add all changes to staging
      git add .

      # Commit changes
      commit_message="Rectifying"
      git commit -m "$commit_message"

      ignore_rector_commit

      # Push changes to origin
      git push origin main --force
      # Uncommitted changes
    fi
}

ignore_rector_commit() {
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

    commit_updated_blame_file
}

commit_updated_blame_file() {
    git add .
    git commit -m "Ignore Rector Commit in Git Blame"
}

run_action