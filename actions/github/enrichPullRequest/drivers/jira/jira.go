package jira

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/EncoreDigitalGroup/golib/logger"
	"github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"

	branchname "github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/drivers/branch_name"
	"github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/support/github"
)

type Configuration struct {
	Enable   bool
	URL      string
	Email    string
	Token    string
	IssueKey string
}

type Information struct {
	ParentPrefix string
	Title        string
	Description  string
	HasJiraInfo  bool
	AuthFailure  bool
}

type JiraError struct {
	IsAuthFailure bool
	OriginalError error
}

func (e *JiraError) Error() string {
	return e.OriginalError.Error()
}

func Format(gh github.GitHub) {
	if jiraLabelSyncEnabled() && gh.HasLabel(jiraLabelSyncName()) {
		logger.Info("PR already has '" + jiraLabelSyncName() + "' label, skipping Jira sync")
		return
	}

	branchName, err := gh.GetBranchName()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	envUrl := os.Getenv("OPT_JIRA_URL")
	if envUrl == "" {
		logger.Error("OPT_JIRA_URL is not set")
		os.Exit(1)
	}

	envEmail := os.Getenv("OPT_JIRA_EMAIL")
	if envEmail == "" {
		logger.Error("OPT_JIRA_EMAIL is not set")
		os.Exit(1)
	}

	envToken := os.Getenv("OPT_JIRA_TOKEN")
	if envToken == "" {
		logger.Error("OPT_JIRA_TOKEN is not set")
		os.Exit(1)
	}

	issueKey, err := branchname.GetIssueKeyFromBranchName(branchName)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if issueKey == "" {
		logger.Error("Issue key is empty")
		return
	}

	config := Configuration{
		Enable:   true,
		URL:      envUrl,
		Email:    envEmail,
		Token:    envToken,
		IssueKey: issueKey,
	}

	jira := getJiraInfo(config)

	if jira.AuthFailure {
		logger.Errorf("Jira authentication failed")

		comment := "**Jira Authentication Failed**\n\n" +
			"Unable to authenticate with Jira to fetch issue information. " +
			"Please verify that the Jira credentials (URL, email, and token) are correctly configured and that the token has not expired.\n\n" +
			"**Possible solutions:**\n" +
			"- Check that `OPT_JIRA_URL`, `OPT_JIRA_EMAIL`, and `OPT_JIRA_TOKEN` environment variables are set correctly\n" +
			"- Verify that the Jira API token is still valid\n" +
			"- Ensure the Jira user has permission to access the issue: `" + issueKey + "`"

		gh.AddPRComment(comment)
		return
	}

	if config.Enable && !jira.HasJiraInfo {
		logger.Errorf("Failed to get Jira info")
		comment := "Failed to get information from Jira.\n\n" +
			"Please check the GitHub Action logs for specific error information."

		gh.AddPRComment(comment)

		return
	}

	newPRTitle := gh.ApplyFormatting(issueKey, jira.Title)

	if jira.ParentPrefix != "" {
		newPRTitle = fmt.Sprintf("[%s]%s", jira.ParentPrefix, newPRTitle)
	}

	if jiraDescriptionSyncEnabled() {
		gh.UpdatePR(newPRTitle, jira.Description)
	} else {
		gh.UpdatePRTitle(newPRTitle)
	}

	if jiraLabelSyncEnabled() {
		gh.EnsureLabelExists(jiraLabelSyncName(), "Indicates that Jira synchronization has been completed for this PR", "0052cc")
		gh.AddLabelToPR(jiraLabelSyncName())
	}
}

func createJiraClient(jiraURL, jiraEmail, jiraToken string) (*v3.Client, error) {
	client, err := v3.New(nil, jiraURL)
	if err != nil {
		return nil, err
	}

	client.Auth.SetBasicAuth(jiraEmail, jiraToken)
	return client, nil
}

func getCurrentIssueInfo(client *v3.Client, issueKey string) (*models.IssueScheme, error) {
	issue, response, err := client.Issue.Get(context.Background(), issueKey, nil, nil)
	if err != nil {
		isAuthFailure := response != nil && response.StatusCode == http.StatusUnauthorized
		return nil, &JiraError{
			IsAuthFailure: isAuthFailure,
			OriginalError: fmt.Errorf("failed to fetch Jira issue %s: %v", issueKey, err),
		}
	}

	return issue, nil
}

func getParentIssuePrefix(client *v3.Client, issueKey string) (string, error) {
	issue, _, err := client.Issue.Get(context.Background(), issueKey, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Jira issue %s: %v", issueKey, err)
	}

	if issue.Fields.Parent == nil {
		return "", nil
	}

	parentIssue, _, err := client.Issue.Get(context.Background(), issue.Fields.Parent.Key, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch parent Jira issue %s: %v", issue.Fields.Parent.Key, err)
	}

	if strings.ToLower(parentIssue.Fields.IssueType.Name) == "epic" {
		return "", nil
	}

	return fmt.Sprintf("%s", issue.Fields.Parent.Key), nil
}

func getJiraInfo(config Configuration) Information {
	if !config.Enable {
		return Information{HasJiraInfo: false}
	}

	if config.URL == "" || config.Email == "" || config.Token == "" {
		logger.Error("OPT_JIRA_URL, OPT_JIRA_EMAIL, and OPT_JIRA_TOKEN must be set when configured strategy is 'jira'.")
		return Information{HasJiraInfo: false}
	}

	client, err := createJiraClient(config.URL, config.Email, config.Token)
	if err != nil {
		logger.Errorf("Failed to create Jira client: %v", err)
		return Information{HasJiraInfo: false}
	}

	result := Information{HasJiraInfo: true}

	jiraIssue, err := getCurrentIssueInfo(client, config.IssueKey)
	if err != nil {
		logger.Errorf("Failed to get current issue info: %v", err)
		var jiraErr *JiraError
		if errors.As(err, &jiraErr) && jiraErr.IsAuthFailure {
			return Information{HasJiraInfo: false, AuthFailure: true}
		}

		return Information{HasJiraInfo: false}
	}
	result.Title = jiraIssue.Fields.Summary

	// Get parent issue prefix if applicable
	parentPrefix, err := getParentIssuePrefix(client, config.IssueKey)
	if err != nil {
		logger.Errorf("Failed to get parent issue info: %v", err)
		// Don't fail completely, just continue without parent prefix
	} else {
		result.ParentPrefix = parentPrefix
	}

	if jiraIssue.Fields.Description != nil {
		result.Description = jiraIssue.Fields.Description.Text
	}

	return result
}

func jiraLabelSyncEnabled() bool {
	return strings.ToLower(os.Getenv("OPT_ENABLE_JIRA_SYNC_LABEL")) == "true"
}

func jiraDescriptionSyncEnabled() bool {
	return strings.ToLower(os.Getenv("OPT_ENABLE_JIRA_SYNC_DESCRIPTION")) == "true"
}

func jiraLabelSyncName() string {
	label := os.Getenv("OPT_JIRA_SYNC_LABEL")

	if label == "" {
		return "jira-sync-complete"
	}

	return label
}
