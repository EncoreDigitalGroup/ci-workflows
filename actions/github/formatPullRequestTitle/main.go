package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"
)

const envGHToken = "GH_TOKEN"
const envGHRepository = "GH_REPOSITORY"
const envPRNumber = "PR_NUMBER"
const envBranchName = "BRANCH_NAME"

// Main function to execute the program
func main() {
	// Retrieve environment variables
	githubToken := os.Getenv(envGHToken)
	if githubToken == "" {
		log.Fatal(envGHToken + " environment variable is not set")
	}

	repo := os.Getenv(envGHRepository)
	if repo == "" {
		log.Fatal(envGHRepository + " environment variable is not set")
	}

	prNumberStr := os.Getenv(envPRNumber)
	if prNumberStr == "" {
		log.Fatal(envPRNumber + " environment variable is not set")
	}

	branchName := os.Getenv(envBranchName)
	if branchName == "" {
		log.Fatal(envBranchName + " environment variable is not set")
	}

	// Split GH_REPOSITORY into owner and name
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		log.Fatal(envGHRepository + " must be in the format owner/repo")
	}
	repoOwner := parts[0]
	repoName := parts[1]

	// Convert PR_NUMBER to integer
	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		log.Fatalf(envPRNumber+" is not a valid integer: %v", err)
	}

	// Initialize GitHub client with OAuth2 authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Define regular expressions for title formatting
	regexWithIssueType := regexp.MustCompile(`^(epic|feature|bugfix|hotfix)/([A-Z]+-[0-9]+)-(.+)$`)
	regexWithoutIssueType := regexp.MustCompile(`^([A-Z]+-[0-9]+)-(.+)$`)

	// formatTitle formats the title based on regex patterns
	formatTitle := func(title string) string {
		var issueKey, issueName string
		if matches := regexWithIssueType.FindStringSubmatch(title); matches != nil {
			issueKey = matches[2]
			issueName = matches[3]
		} else if matches := regexWithoutIssueType.FindStringSubmatch(title); matches != nil {
			issueKey = matches[1]
			issueName = matches[2]
		} else {
			fmt.Println("Title does not match expected format")
			return title
		}

		// Replace hyphens with spaces and capitalize each word
		formattedIssueName := strings.ReplaceAll(issueName, "-", " ")
		String.SetLanguage(language.English)
		formattedIssueName = String.TitleCase(formattedIssueName)
		return fmt.Sprintf("[%s] %s", issueKey, formattedIssueName)
	}

	// updatePullRequestTitle updates the PR title using the GitHub SDK
	updatePullRequestTitle := func(title string) {
		formattedTitle := formatTitle(title)
		fmt.Println("Attempting to Update Pull Request Title to:", formattedTitle)

		_, _, err := client.PullRequests.Edit(ctx, repoOwner, repoName, prNumber, &github.PullRequest{
			Title: &formattedTitle,
		})
		if err != nil {
			log.Fatalf("Failed to update pull request title: %v", err)
		}
		fmt.Println("Updated Pull Request Title to:", formattedTitle)
	}

	// branchNameMatches checks if the current PR title matches the formatted branch name
	branchNameMatches := func() bool {
		pr, _, err := client.PullRequests.Get(ctx, repoOwner, repoName, prNumber)
		if err != nil {
			log.Fatalf("Failed to get pull request: %v", err)
		}
		pullRequestTitle := *pr.Title
		fmt.Println("Pull Request Title is:", pullRequestTitle)

		formattedBranchName := formatTitle(branchName)
		if pullRequestTitle == formattedBranchName {
			fmt.Println("Pull Request Titles Match")
			return true
		}
		fmt.Println("Pull Request Titles Do Not Match")
		return false
	}

	// Main logic: update title if it doesn't match
	if !branchNameMatches() {
		fmt.Println("Pull Request Title Should Be Updated.")
		updatePullRequestTitle(branchName)
	}
}
