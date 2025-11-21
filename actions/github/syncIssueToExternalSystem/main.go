package main

import (
    "context"
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/EncoreDigitalGroup/golib/logger"
    "github.com/ctreminiom/go-atlassian/jira/v3"
    "github.com/ctreminiom/go-atlassian/pkg/infra/models"
    "github.com/google/go-github/v70/github"
    "golang.org/x/oauth2"
)

type Config struct {
    GitHubToken              string
    IssueNumber              int
    IssueTitle               string
    IssueBody                string
    IssueURL                 string
    Repository               string
    Assignees                []string
    Labels                   []string
    ExternalSystemURL        string
    ExternalSystemEmail      string
    ExternalSystemAuthToken  string
    ExternalSystemProjectKey string
}

func main() {
    config, err := loadConfig()
    if err != nil {
        logger.Errorf("Failed to load configuration: %v", err)
        os.Exit(1)
    }

    err = syncIssueToJira(config)
    if err != nil {
        logger.Errorf("Failed to sync issue to Jira: %v", err)
        os.Exit(1)
    }

    logger.Info("Successfully synced GitHub issue to Jira")
}

func loadConfig() (*Config, error) {
    issueNumber, err := strconv.Atoi(os.Getenv("ISSUE_NUMBER"))
    if err != nil {
        return nil, fmt.Errorf("invalid issue number: %v", err)
    }

    assignees := strings.Split(os.Getenv("ASSIGNEES"), ",")
    if len(assignees) == 1 && assignees[0] == "" {
        assignees = []string{}
    }

    labels := strings.Split(os.Getenv("LABELS"), ",")
    if len(labels) == 1 && labels[0] == "" {
        labels = []string{}
    }

    config := &Config{
        GitHubToken:              os.Getenv("GITHUB_TOKEN"),
        IssueNumber:              issueNumber,
        IssueTitle:               os.Getenv("ISSUE_TITLE"),
        IssueBody:                os.Getenv("ISSUE_BODY"),
        IssueURL:                 os.Getenv("ISSUE_URL"),
        Repository:               os.Getenv("REPOSITORY"),
        Assignees:                assignees,
        Labels:                   labels,
        ExternalSystemURL:        os.Getenv("OPT_ENDPOINT"),
        ExternalSystemEmail:      os.Getenv("OPT_EMAIL"),
        ExternalSystemAuthToken:  os.Getenv("OPT_AUTH_TOKEN"),
        ExternalSystemProjectKey: os.Getenv("OPT_PROJECT_KEY"),
    }

    if config.GitHubToken == "" {
        return nil, fmt.Errorf("GITHUB_TOKEN is required")
    }
    if config.ExternalSystemURL == "" {
        return nil, fmt.Errorf("OPT_ENDPOINT is required")
    }
    if config.ExternalSystemEmail == "" {
        return nil, fmt.Errorf("OPT_EMAIL is required")
    }
    if config.ExternalSystemAuthToken == "" {
        return nil, fmt.Errorf("OPT_AUTH_TOKEN is required")
    }
    if config.ExternalSystemProjectKey == "" {
        return nil, fmt.Errorf("OPT_PROJECT_KEY is required")
    }

    return config, nil
}

func syncIssueToJira(config *Config) error {
    ctx := context.Background()

    // Initialize Jira client
    jiraClient, err := v3.New(nil, config.ExternalSystemURL)
    if err != nil {
        return fmt.Errorf("failed to create Jira client: %v", err)
    }

    jiraClient.Auth.SetBasicAuth(config.ExternalSystemEmail, config.ExternalSystemAuthToken)

    // Prepare issue type - default to "Task"
    issueType := "Task"

    // Check if any label indicates this should be a "Bug"
    for _, label := range config.Labels {
        if strings.ToLower(label) == "bug" {
            issueType = "Bug"
            break
        }
    }

    // Create Jira issue
    payload := &models.IssueScheme{
        Fields: &models.IssueFieldsScheme{
            Summary: config.IssueTitle,
            Project: &models.ProjectScheme{
                Key: config.ExternalSystemProjectKey,
            },
            IssueType: &models.IssueTypeScheme{
                Name: issueType,
            },
        },
    }

    // Add labels if provided
    if len(config.Labels) > 0 {
        payload.Fields.Labels = config.Labels
    }

    // Create the issue
    createdIssue, response, err := jiraClient.Issue.Create(ctx, payload, nil)
    if err != nil {
        return fmt.Errorf("failed to create Jira issue: %v, response: %v", err, response)
    }

    logger.Infof("Created Jira issue: %s", createdIssue.Key)

    // Add comment to GitHub issue with Jira link
    err = addGitHubComment(config, createdIssue.Key)
    if err != nil {
        logger.Warnf("Failed to add comment to GitHub issue: %v", err)
        // Don't fail the entire operation if comment fails
    }

    return nil
}

func addGitHubComment(config *Config, jiraKey string) error {
    ctx := context.Background()

    // Create GitHub client
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: config.GitHubToken},
    )
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    // Parse repository
    repoParts := strings.Split(config.Repository, "/")
    if len(repoParts) != 2 {
        return fmt.Errorf("invalid repository format: %s", config.Repository)
    }

    owner := repoParts[0]
    repo := repoParts[1]

    // Create comment
    jiraURL := fmt.Sprintf("%s/browse/%s", strings.TrimSuffix(config.ExternalSystemURL, "/"), jiraKey)
    commentBody := fmt.Sprintf("🔗 **Jira Issue Created:** [%s](%s)", jiraKey, jiraURL)

    comment := &github.IssueComment{
        Body: &commentBody,
    }

    _, _, err := client.Issues.CreateComment(ctx, owner, repo, config.IssueNumber, comment)
    if err != nil {
        return fmt.Errorf("failed to create GitHub comment: %v", err)
    }

    logger.Infof("Added comment to GitHub issue #%d", config.IssueNumber)
    return nil
}
