package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/EncoreDigitalGroup/ci-workflows/actions/github/formatPullRequestTitle/drivers"
	"github.com/EncoreDigitalGroup/ci-workflows/actions/github/formatPullRequestTitle/drivers/branch_name"
	"github.com/EncoreDigitalGroup/ci-workflows/actions/github/formatPullRequestTitle/drivers/jira"
	"github.com/EncoreDigitalGroup/ci-workflows/actions/github/formatPullRequestTitle/support/github"
	"github.com/EncoreDigitalGroup/golib/logger"
)

var gh github.GitHub

const envGHRepository = "GH_REPOSITORY"
const envPRNumber = "PR_NUMBER"
const envBranchName = "BRANCH_NAME"
const envEnableExperiments = "ENABLE_EXPERIMENTS"
const envStrategy = "CI_FMT_STRATEGY"

// Retrieve environment variables
var strategy = os.Getenv(envStrategy)
var repo = os.Getenv(envGHRepository)
var prNumberStr = os.Getenv(envPRNumber)
var branchName = os.Getenv(envBranchName)
var enableExperiments = getEnableExperiments()
var parts = strings.Split(repo, "/")

var pullRequestTitle = ""

// getEnableExperiments retrieves the ENABLE_EXPERIMENTS environment variable as a boolean
func getEnableExperiments() bool {
	envValue := os.Getenv(envEnableExperiments)
	return strings.ToLower(envValue) == "true"
}

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

	if branchName == "" {
		logger.Error(envBranchName + " environment variable is not set")
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
