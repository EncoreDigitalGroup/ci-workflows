package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/EncoreDigitalGroup/golib/logger"
	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const envGHToken = "GH_TOKEN"
const envGHRepository = "GH_REPOSITORY"
const envPRNumber = "PR_NUMBER"
const envBranchName = "BRANCH_NAME"

// Retrieve environment variables
var githubToken = os.Getenv(envGHToken)
var repo = os.Getenv(envGHRepository)
var prNumberStr = os.Getenv(envPRNumber)
var branchName = os.Getenv(envBranchName)
var parts = strings.Split(repo, "/")

// Initialize GitHub client with OAuth2 authentication
var ctx = context.Background()
var ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
var tc = oauth2.NewClient(ctx, ts)
var client = github.NewClient(tc)

// Define regular expressions for title formatting
var regexWithIssueType = regexp.MustCompile(`^(epic|feature|bugfix|hotfix)/([A-Z]+-[0-9]+)-(.+)$`)
var regexWithoutIssueType = regexp.MustCompile(`^([A-Z]+-[0-9]+)-(.+)$`)
var pullRequestTitle = ""

// Main function to execute the program
func main() {
	if githubToken == "" {
		logger.Error(envGHToken + " environment variable is not set")
	}

	if repo == "" {
		logger.Error(envGHRepository + " environment variable is not set")
	}

	if prNumberStr == "" {
		logger.Error(envPRNumber + " environment variable is not set")
	}

	if branchName == "" {
		logger.Error(envBranchName + " environment variable is not set")
	}

	if len(parts) != 2 {
		logger.Error(envGHRepository + " must be in the format owner/repo")
	}
	repoOwner := parts[0]
	repoName := parts[1]

	// Convert PR_NUMBER to integer
	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		logger.Errorf(envPRNumber+" is not a valid integer: %v", err)
	}

	// Main logic: update title if it doesn't match
	if !branchNameMatches(repoOwner, repoName, prNumber) {
		fmt.Println("Pull Request Title Should Be Updated.")
		updatePullRequestTitle(repoOwner, repoName, prNumber, branchName)
	}
}

func formatTitle(title string) string {
	var issueKey, issueName string
	if matches := regexWithIssueType.FindStringSubmatch(title); matches != nil {
		issueKey = matches[2]
		issueName = matches[3]
	} else if matches := regexWithoutIssueType.FindStringSubmatch(title); matches != nil {
		issueKey = matches[1]
		issueName = matches[2]
	} else {
		fmt.Println("Title does not match expected format")
		logger.Info(title)
		return pullRequestTitle
	}

	// Replace hyphens with spaces and capitalize each word
	formattedIssueName := strings.ReplaceAll(issueName, "-", " ")
	titleCaser := cases.Title(language.English)
	formattedIssueName = titleCaser.String(formattedIssueName)

	defaultExceptions := map[string]string{
		"Api":          "API",
		"Css":          "CSS",
		"Db":           "DB",
		"Html":         "HTML",
		"Rest":         "REST",
		"Rockrms":      "RockRMS",
		"Mpc":          "MPC",
		"Myportal":     "MyPortal",
		"Pco":          "PCO",
		"Php":          "PHP",
		"Phpstan":      "PHPStan",
		"Servicepoint": "ServicePoint",
		"Themekit":     "ThemeKit",
		"Uri":          "URI",
		"Webcms":       "WebCMS",
		"Webui":        "WebUI",
	}

	if userDefinedExceptions := os.Getenv("CI_FMT_WORDS"); userDefinedExceptions != "" {
		pairs := strings.Split(userDefinedExceptions, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				defaultExceptions[key] = value
			}
		}
	}
	words := strings.Fields(formattedIssueName)
	for i, word := range words {
		if val, ok := defaultExceptions[word]; ok {
			words[i] = val
		}
	}
	formattedIssueName = strings.Join(words, " ")

	return fmt.Sprintf("[%s] %s", issueKey, formattedIssueName)
}

func updatePullRequestTitle(repoOwner string, repoName string, prNumber int, prTitle string) {
	formattedTitle := formatTitle(prTitle)
	fmt.Println("Attempting to Update Pull Request Title to:", formattedTitle)

	_, _, err := client.PullRequests.Edit(ctx, repoOwner, repoName, prNumber, &github.PullRequest{
		Title: &formattedTitle,
	})
	if err != nil {
		logger.Errorf("Failed to update pull request prTitle: %v", err)
	}
	logger.Infof("Updated Pull Request Title to: %s", formattedTitle)
}

func branchNameMatches(repoOwner string, repoName string, prNumber int) bool {
	pullRequest, _, err := client.PullRequests.Get(ctx, repoOwner, repoName, prNumber)
	if err != nil {
		logger.Errorf("Failed to get pullRequest request: %v", err)
	}
	pullRequestTitle = *pullRequest.Title
	logger.Info("Pull Request Title is:", pullRequestTitle)

	formattedBranchName := formatTitle(branchName)
	if pullRequestTitle == formattedBranchName {
		logger.Info("Pull Request Titles Match; No Need to Update.")
		return true
	}
	logger.Info("Pull Request Titles Do Not Match; Update Needed.")
	return false
}
