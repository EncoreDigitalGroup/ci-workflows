package general

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "time"

    "github.com/EncoreDigitalGroup/golib/logger"

    branchname "github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/drivers/branch_name"
    "github.com/EncoreDigitalGroup/ci-workflows/actions/github/enrichPullRequest/support/github"
)

type Configuration struct {
    Enable    bool
    Endpoint  string
    AuthToken string
    TicketID  string
}

type Label struct {
    Title       string  `json:"title"`
    Description *string `json:"description"`
    Color       *string `json:"color"`
}

type APIResponse struct {
    Title       string   `json:"title"`
    Description *string  `json:"description"`
    Assignee    *string  `json:"assignee"`
    Labels      *[]Label `json:"labels"`
}

type HTTPError struct {
    StatusCode int
    Message    string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func Format(gh github.GitHub) {
    branchName, err := gh.GetBranchName()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    endpoint := os.Getenv("OPT_HTTP_ENDPOINT")
    if endpoint == "" {
        logger.Error("OPT_HTTP_ENDPOINT is not set")
        os.Exit(1)
    }

    authToken := os.Getenv("OPT_AUTH_TOKEN")
    if authToken == "" {
        logger.Error("OPT_AUTH_TOKEN is not set")
        os.Exit(1)
    }

    ticketID, err := branchname.GetIssueKeyFromBranchName(branchName)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    if ticketID == "" {
        logger.Error("Ticket ID is empty")
        return
    }

    config := Configuration{
        Enable:    true,
        Endpoint:  endpoint,
        AuthToken: authToken,
        TicketID:  ticketID,
    }

    apiResponse, err := getTicketInfo(config)
    if err != nil {
        var httpErr *HTTPError
        if errors.As(err, &httpErr) {
            if httpErr.StatusCode == 404 {
                logger.Errorf("Ticket not found: %s", ticketID)
                comment := fmt.Sprintf("**Ticket Not Found**\n\n"+
                    "Unable to find ticket information for: `%s`\n\n"+
                    "Please verify that the ticket ID in your branch name is correct.", ticketID)
                gh.AddPRComment(comment)
                return
            }

            logger.Errorf("HTTP API error: %v", httpErr)
            comment := "**API Error**\n\n" +
                "Failed to fetch ticket information due to an API error. " +
                "Please check the GitHub Action logs for specific error information."
            gh.AddPRComment(comment)
            return
        }

        logger.Errorf("Failed to get ticket info: %v", err)
        comment := "Failed to get information from the API.\n\n" +
            "Please check the GitHub Action logs for specific error information."
        gh.AddPRComment(comment)
        return
    }

    newPRTitle := apiResponse.Title
    if newPRTitle == "" {
        logger.Error("API response missing required 'title' field")
        comment := "**Invalid API Response**\n\n" +
            "The API response is missing the required 'title' field."
        gh.AddPRComment(comment)
        return
    }

    var newPRDescription string
    if apiResponse.Description != nil {
        newPRDescription = *apiResponse.Description
        gh.UpdatePR(newPRTitle, newPRDescription)
    } else {
        gh.UpdatePRTitle(newPRTitle)
    }

    if apiResponse.Assignee != nil && *apiResponse.Assignee != "" {
        logger.Infof("Assignee information received: %s (Note: GitHub assignee setting not implemented)", *apiResponse.Assignee)
    }

    if apiResponse.Labels != nil {
        for _, label := range *apiResponse.Labels {
            if label.Title == "" {
                continue
            }

            description := ""
            if label.Description != nil {
                description = *label.Description
            }

            color := ""
            if label.Color != nil {
                color = *label.Color
            }

            gh.EnsureLabelExists(label.Title, description, color)
            gh.AddLabelToPR(label.Title)
        }
    }
}

func getTicketInfo(config Configuration) (*APIResponse, error) {
    if !config.Enable {
        return nil, errors.New("general driver is not enabled")
    }

    if config.Endpoint == "" || config.AuthToken == "" || config.TicketID == "" {
        return nil, errors.New("missing required configuration: endpoint, auth token, or ticket ID")
    }

    reqURL, err := url.Parse(config.Endpoint)
    if err != nil {
        return nil, fmt.Errorf("invalid endpoint URL: %v", err)
    }

    query := reqURL.Query()
    query.Set("id", config.TicketID)
    reqURL.RawQuery = query.Encode()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Authorization", "Bearer "+config.AuthToken)
    req.Header.Set("Accept", "application/json")

    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make HTTP request: %v", err)
    }
    defer func(Body io.ReadCloser) {
        err := Body.Close()
        if err != nil {
            logger.Errorf("Failed to close response body: %v", err)
        }
    }(resp.Body)

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    if resp.StatusCode == 404 {
        return nil, &HTTPError{
            StatusCode: resp.StatusCode,
            Message:    "ticket not found",
        }
    }

    if resp.StatusCode != 200 {
        return nil, &HTTPError{
            StatusCode: resp.StatusCode,
            Message:    string(body),
        }
    }

    var apiResponse APIResponse
    if err := json.Unmarshal(body, &apiResponse); err != nil {
        return nil, fmt.Errorf("failed to parse JSON response: %v", err)
    }

    return &apiResponse, nil
}
