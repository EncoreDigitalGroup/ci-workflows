#!/bin/bash

echo "Checking for changes in specific files..."

# Fetch the pull request's merge commit
pr_number=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")
git fetch origin +refs/pull/"${pr_number}"/merge
diff_files=$(git diff --name-only origin/"$BRANCH_NAME" FETCH_HEAD)

# Output the diff files for debugging
echo "Changed files:"
echo "$diff_files"

# Initialize the flag
should_run_entire_workflow=false

# Check each file in the diff
for file in $diff_files; do
    if [[ "$file" == *.php ]] ||
       [[ "$file" == *composer.json ]] ||
       [[ "$file" == *composer.lock ]] ||
       [[ "$file" == *package.json ]] ||
       [[ "$file" == *package-lock.json ]]; then
        echo "Triggering file is $file"
        should_run_entire_workflow=true
        break
    fi
done

# Set the output based on the flag
if [ "$should_run_entire_workflow" = true ]; then
    echo "Setting shouldRun to $should_run_entire_workflow"
    echo "shouldRun=true" >> "$GITHUB_OUTPUT"
else
    echo "Setting shouldRun to $should_run_entire_workflow"
    echo "shouldRun=false" >> "$GITHUB_OUTPUT"
fi