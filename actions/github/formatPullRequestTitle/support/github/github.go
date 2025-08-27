package github

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/EncoreDigitalGroup/golib/logger"
	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const envGHToken = "GH_TOKEN"
const envBranchName = "BRANCH_NAME"

// GitHub interface defines the contract for GitHub operations
type GitHub interface {
	GetBranchName() (string, error)
	BranchNameMatchesPRTitle(currentPRTitle string) bool
	UpdatePRTitle(newPRTitle string)
	ApplyFormatting(issueKey string, issueName string) string
}

// GitHubClient implements the GitHub interface
type GitHubClient struct {
	client            *github.Client
	repositoryOwner   string
	repositoryName    string
	pullRequestNumber int
}

var (
	client *github.Client
	once   sync.Once
)

func New(repoOwner string, repoName string, prNumber int) GitHub {
	once.Do(func() {
		githubToken := os.Getenv(envGHToken)

		if githubToken == "" {
			logger.Error(envGHToken + " is not set")
			os.Exit(1)
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	})

	return &GitHubClient{
		client:            client,
		repositoryOwner:   repoOwner,
		repositoryName:    repoName,
		pullRequestNumber: prNumber,
	}
}

func (gh *GitHubClient) GetBranchName() (string, error) {
	branchName := os.Getenv(envBranchName)

	if branchName == "" {
		logger.Error(envBranchName + " is not set")
		return "", fmt.Errorf("%s is not set", envBranchName)
	}

	return branchName, nil
}

func (gh *GitHubClient) BranchNameMatchesPRTitle(currentPRTitle string) bool {
	pullRequest, _, err := gh.client.PullRequests.Get(context.Background(), gh.repositoryOwner, gh.repositoryName, gh.pullRequestNumber)
	if err != nil {
		logger.Errorf("Failed to get pullRequest request: %v", err)
	}

	if currentPRTitle == *pullRequest.Title {
		logger.Info("Pull Request Titles Match; No Need to Update.")
		return true
	}

	logger.Info("Pull Request Titles Do Not Match; Update Needed.")
	return false
}

func (gh *GitHubClient) UpdatePRTitle(newPRTitle string) {
	fmt.Println("Attempting to Update Pull Request Title to:", newPRTitle)

	_, _, err := gh.client.PullRequests.Edit(context.Background(), gh.repositoryOwner, gh.repositoryName, gh.pullRequestNumber, &github.PullRequest{
		Title: &newPRTitle,
	})

	if err != nil {
		logger.Errorf("Failed to update pull request prTitle: %v", err)
	}

	logger.Infof("Updated Pull Request Title to: %s", newPRTitle)
}

func (gh *GitHubClient) ApplyFormatting(issueKey string, issueName string) string {
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
	words := strings.Fields(issueName)
	for i, word := range words {
		if val, ok := defaultExceptions[word]; ok {
			words[i] = val
		}
	}

	return fmt.Sprintf("[%s] %s", issueKey, formattedIssueName)
}
