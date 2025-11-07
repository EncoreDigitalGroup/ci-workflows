package main

import (
    "os"
    "strconv"
    "strings"

    "github.com/EncoreDigitalGroup/golib/logger"

    "github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/drivers"
    "github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/drivers/branch_name"
    "github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/drivers/jira"
    "github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/support/github"
)

var gh github.GitHub

const envGHRepository = "GH_REPOSITORY"
const envPRNumber = "PR_NUMBER"
const envStrategy = "OPT_FMT_STRATEGY"

// Retrieve environment variables
var strategy = os.Getenv(envStrategy)
var repo = os.Getenv(envGHRepository)
var prNumberStr = os.Getenv(envPRNumber)
var parts = strings.Split(repo, "/")

// Main function to execute the program
func main() {
    checkEnvVars()
    repoOwner := parts[0]
    repoName := parts[1]

    // Convert PR_NUMBER to integer
    prNumber, err := strconv.Atoi(prNumberStr)
    if err != nil {
        logger.Errorf(envPRNumber+" is not a valid integer: %v", err)
    }

    gh = github.New(repoOwner, repoName, prNumber)

    if strategy == drivers.BranchName {
        branchname.Format(gh)
    }

    if strategy == drivers.Jira {
        jira.Format(gh)
    }
}

func checkEnvVars() {
    isMissingVar := false
    if strategy == "" {
        logger.Error(envStrategy + " environment variable is not set")
        isMissingVar = true
    }

    if repo == "" {
        logger.Error(envGHRepository + " environment variable is not set")
        isMissingVar = true
    }

    if prNumberStr == "" {
        logger.Error(envPRNumber + " environment variable is not set")
        isMissingVar = true
    }

    if len(parts) != 2 {
        logger.Error(envGHRepository + " must be in the format owner/repo")
        isMissingVar = true
    }

    if isMissingVar {
        os.Exit(1)
    }
}
