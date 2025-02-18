#!/bin/bash

githubToken=$GH_TOKEN
repoOwner=$(echo "$GH_REPOSITORY" | cut -d'/' -f1)
repoName=$(echo "$GH_REPOSITORY" | cut -d'/' -f2)
prNumber=$PR_NUMBER
branchName=$BRANCH_NAME
prTitle=''
useSemanticCommits=${USE_SEMANTIC_COMMITS:-false} # Default to false if not set

# Add regex for semantic commit message types
semanticPrefix="^(feat|fix|chore|refactor|docs|test)(\(.+\))?:"
regexWithIssueType="^(epic|feature|bugfix|hotfix)/([A-Z]+-[0-9]+)-(.+)$"
regexWithoutIssueType="^([A-Z]+-[0-9]+)-(.+)$"

formatTitle() {
  local title="$1"
  local formattedIssueName=""
  local fullFormattedName=""

  # Check if semantic commit messages are enabled
  if [[ "$useSemanticCommits" != "true" ]]; then
    echo "Semantic commit messages are disabled. Using default format."
    if [[ $title =~ $regexWithIssueType ]]; then
      issueKey="${BASH_REMATCH[2]}"
      issueName="${BASH_REMATCH[3]}"
    elif [[ $title =~ $regexWithoutIssueType ]]; then
      issueKey="${BASH_REMATCH[1]}"
      issueName="${BASH_REMATCH[2]}"
    else
      echo "Title does not match expected format"
      echo "$title"
      return
    fi

    formattedIssueName=$(echo "$issueName" | sed -e 's/-/ /g' -e 's/\b\w/\u&/g')
    fullFormattedName="[$issueKey] $formattedIssueName"
    echo "$fullFormattedName"
    return
  fi

  # Semantic commits are enabled, check for existing prefix
  if [[ $title =~ $semanticPrefix ]]; then
    echo "Branch name already follows semantic commit message format."
    echo "$title"
    return
  fi

  # Match branch name with issue type
  if [[ $title =~ $regexWithIssueType ]]; then
    issueKey="${BASH_REMATCH[2]}"
    issueName="${BASH_REMATCH[3]}"
    typePrefix="${BASH_REMATCH[1]}"
  # Match branch name without issue type
  elif [[ $title =~ $regexWithoutIssueType ]]; then
    issueKey="${BASH_REMATCH[1]}"
    issueName="${BASH_REMATCH[2]}"
    typePrefix="feat"  # Default to 'feat' if no type is specified
  else
    echo "Title does not match expected format"
    echo "$title"
    return
  fi

  # Format the issue name for the PR title
  formattedIssueName=$(echo "$issueName" | sed -e 's/-/ /g' -e 's/\b\w/\u&/g')
  
  # Construct the full formatted title with semantic prefix
  fullFormattedName="$typePrefix: [$issueKey] $formattedIssueName"
  echo "$fullFormattedName"
}

updatePullRequestTitle() {
  local title
  title=$(formatTitle "$1")
  echo "Attempting to Update Pull Request Title to: $title"

  curl -X PATCH \
    -H "Authorization: token $githubToken" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"title\":\"$title\"}" \
    "https://api.github.com/repos/$repoOwner/$repoName/pulls/$prNumber"
}

branchNameMatches() {
  local pullRequestTitle
  pullRequestTitle=$(curl -s \
    -H "Authorization: token $githubToken" \
    -H "Accept: application/vnd.github.v3+json" \
    "https://api.github.com/repos/$repoOwner/$repoName/pulls/$prNumber" | jq -r '.title')

  echo "Pull Request Title is: $pullRequestTitle"
  prTitle="$pullRequestTitle"

  if [[ "$prTitle" == "$(formatTitle "$branchName")" ]]; then
    echo "Pull Request Titles Match"
    return 0
  else
    echo "Pull Request Titles Do Not Match"
    return 1
  fi
}

act() {
  if ! branchNameMatches "$branchName"; then
    echo "Pull Request Title Should Be Updated."
    updatePullRequestTitle "$branchName"
  fi
}

act
