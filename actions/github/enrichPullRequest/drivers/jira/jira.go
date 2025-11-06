package jira

import (
    "context"
    "fmt"
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
    Token    string
    IssueKey string
}

type Information struct {
    ParentPrefix string
    Title        string
    HasJiraInfo  bool
}

func Format(gh github.GitHub) {
    if gh.HasLabel("jira-sync-complete") {
        logger.Info("PR already has 'jira-sync-complete' label, skipping Jira sync")
        return
    }

    branchName, err := gh.GetBranchName()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    envUrl := os.Getenv("JIRA_URL")
    if envUrl == "" {
        logger.Error("JIRA_URL is not set")
        os.Exit(1)
    }

    envToken := os.Getenv("JIRA_TOKEN")
    if envToken == "" {
        logger.Error("JIRA_TOKEN is not set")
        os.Exit(1)
    }

    issueKey, err := branchname.GetIssueKeyFromBranchName(branchName)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    config := Configuration{
        Enable:   true,
        URL:      envUrl,
        Token:    envToken,
        IssueKey: issueKey,
    }

    jira := getJiraInfo(config)

    newPRTitle := gh.ApplyFormatting(issueKey, jira.Title)

    if jira.ParentPrefix != "" {
        newPRTitle = fmt.Sprintf("[%s]%s", jira.ParentPrefix, newPRTitle)
    }

    gh.UpdatePRTitle(newPRTitle)

    // Ensure the jira-sync-complete label exists and add it to the PR
    gh.EnsureLabelExists("jira-sync-complete", "Indicates that Jira synchronization has been completed for this PR", "0e8a16")
    gh.AddLabelToPR("jira-sync-complete")
}

func createJiraClient(jiraURL, jiraToken string) (*v3.Client, error) {
    client, err := v3.New(nil, jiraURL)
    if err != nil {
        return nil, err
    }

    client.Auth.SetBearerToken(jiraToken)
    return client, nil
}

func getCurrentIssueInfo(client *v3.Client, issueKey string) (*models.IssueScheme, error) {
    issue, _, err := client.Issue.Get(context.Background(), issueKey, nil, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch Jira issue %s: %v", issueKey, err)
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

    return fmt.Sprintf("[%s]", issue.Fields.Parent.Key), nil
}

func getJiraInfo(config Configuration) Information {
    if config.Enable {
        return Information{HasJiraInfo: false}
    }

    if config.URL == "" || config.Token == "" {
        logger.Error("JIRA_URL and JIRA_TOKEN must be set when ENABLE_JIRA is true")
        return Information{HasJiraInfo: false}
    }

    client, err := createJiraClient(config.URL, config.Token)
    if err != nil {
        logger.Errorf("Failed to create Jira client: %v", err)
        return Information{HasJiraInfo: false}
    }

    result := Information{HasJiraInfo: true}

    // Get current issue jiraIssue
    jiraIssue, err := getCurrentIssueInfo(client, config.IssueKey)
    if err != nil {
        logger.Errorf("Failed to get current issue info: %v", err)
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

    return result
}
