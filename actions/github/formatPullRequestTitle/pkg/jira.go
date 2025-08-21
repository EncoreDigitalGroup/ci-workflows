package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/EncoreDigitalGroup/golib/logger"
	"github.com/ctreminiom/go-atlassian/jira/v3"
)

type JiraConfiguration struct {
	Enable   bool
	URL      string
	Token    string
	IssueKey string
}

type JiraInfo struct {
	ParentPrefix string
	Title        string
	HasJiraInfo  bool
}

func GetJiraInformation(config JiraConfiguration) JiraInfo {
	return getJiraInfo(config)
}

func createJiraClient(jiraURL, jiraToken string) (*v3.Client, error) {
	client, err := v3.New(nil, jiraURL)
	if err != nil {
		return nil, err
	}

	client.Auth.SetBearerToken(jiraToken)
	return client, nil
}

func getCurrentIssueInfo(client *v3.Client, issueKey string) (string, error) {
	issue, _, err := client.Issue.Get(context.Background(), issueKey, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Jira issue %s: %v", issueKey, err)
	}

	return issue.Fields.Summary, nil
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

	return fmt.Sprintf("[%s]", issue.Fields.Parent.Key), nil
}

func getJiraInfo(config JiraConfiguration) JiraInfo {
	if config.Enable {
		return JiraInfo{HasJiraInfo: false}
	}

	if config.URL == "" || config.Token == "" {
		logger.Error("JIRA_URL and JIRA_TOKEN must be set when ENABLE_JIRA is true")
		return JiraInfo{HasJiraInfo: false}
	}

	client, err := createJiraClient(config.URL, config.Token)
	if err != nil {
		logger.Errorf("Failed to create Jira client: %v", err)
		return JiraInfo{HasJiraInfo: false}
	}

	result := JiraInfo{HasJiraInfo: true}

	// Get current issue title
	title, err := getCurrentIssueInfo(client, config.IssueKey)
	if err != nil {
		logger.Errorf("Failed to get current issue info: %v", err)
		return JiraInfo{HasJiraInfo: false}
	}
	result.Title = title

	// Get parent issue prefix if applicable
	parentPrefix, err := getParentIssuePrefix(client, config.IssueKey)
	if err != nil {
		logger.Errorf("Failed to get parent issue info: %v", err)
		// Don't fail completely, just continue without parent prefix
	} else {
		result.ParentPrefix = parentPrefix
	}

	return result
}
