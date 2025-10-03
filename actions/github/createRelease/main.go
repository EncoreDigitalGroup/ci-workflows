package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/EncoreDigitalGroup/golib/logger"
    "github.com/google/go-github/v70/github"
    "golang.org/x/oauth2"
)

func main() {
    githubToken := getEnv("GH_TOKEN")
    repo := getEnv("GH_REPOSITORY")
    tagName := getEnv("TAG_NAME")
    preReleaseStr := getEnv("PRE_RELEASE")
    generateReleaseNotesStr := getEnv("GENERATE_RELEASE_NOTES")
    isDraftStr := getEnv("IS_DRAFT")
    includeDependabotStr := getEnv("INCLUDE_DEPENDABOT")

    preRelease := parseBool(preReleaseStr)
    generateReleaseNotes := parseBool(generateReleaseNotesStr)
    isDraft := parseBool(isDraftStr)
    includeDependabot := parseBool(includeDependabotStr)

    repoParts := strings.Split(repo, "/")
    if len(repoParts) != 2 {
        log.Fatal("Invalid repository format. Expected owner/repo")
    }
    repoOwner := repoParts[0]
    repoName := repoParts[1]

    // Initialize GitHub client
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    // Create the release
    release := &github.RepositoryRelease{
        TagName:    &tagName,
        Name:       &tagName,
        Draft:      &isDraft,
        Prerelease: &preRelease,
    }

    // Handle release notes
    if generateReleaseNotes {
        releaseBody, err := generateCustomReleaseNotes(ctx, client, repoOwner, repoName, includeDependabot)
        if err != nil {
            log.Printf("Warning: Failed to generate custom release notes: %v", err)
            // Fall back to GitHub's auto-generated release notes
            release.GenerateReleaseNotes = &generateReleaseNotes
        } else {
            release.Body = &releaseBody
        }
    }

    // Create the release
    createdRelease, _, err := client.Repositories.CreateRelease(ctx, repoOwner, repoName, release)
    if err != nil {
        log.Fatalf("Failed to create release: %v", err)
    }

    fmt.Printf("Release created successfully: %s\n", *createdRelease.HTMLURL)
}

func getEnv(key string) string {
    value := os.Getenv(key)

    if value == "" {
        logger.Errorf("Environment variable %s is required", key)
        os.Exit(1)
    }

    return value
}

func parseBool(value string) bool {
    val, err := strconv.ParseBool(value)

    if err != nil {
        logger.Errorf("Failed to parse bool value %s: %v", value, err)
        os.Exit(1)
    }

    return val
}

func generateCustomReleaseNotes(ctx context.Context, client *github.Client, owner, repo string, includeDependabot bool) (string, error) {
    // Get the latest release to determine the comparison point
    latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
    var sinceTime *time.Time
    if err == nil {
        releaseTime := latestRelease.CreatedAt.GetTime().Add(time.Second)
        adjustedTime := releaseTime.Add(time.Second)
        sinceTime = &adjustedTime
    }

    // Get merged PRs since the last release
    allPRs, err := getMergedPRsSince(ctx, client, owner, repo, sinceTime)
    if err != nil {
        return "", err
    }

    // Filter and build release notes
    var releaseNotes strings.Builder
    releaseNotes.WriteString("## What's Changed\n\n")

    for _, pr := range allPRs {
        // Handle Dependabot PRs
        isDependabotPR := strings.Contains(strings.ToLower(*pr.Title), "dependabot") ||
            (pr.User != nil && strings.Contains(strings.ToLower(*pr.User.Login), "dependabot"))

        if isDependabotPR && !includeDependabot {
            continue
        }

        // Add PR to release notes
        if pr.User != nil {
            releaseNotes.WriteString(fmt.Sprintf("* %s by @%s in #%d\n", *pr.Title, *pr.User.Login, *pr.Number))
        } else {
            releaseNotes.WriteString(fmt.Sprintf("* %s in #%d \n", *pr.Title, *pr.Number))
        }
    }

    if releaseNotes.Len() == len("## What's Changed\n\n") {
        releaseNotes.WriteString("No changes in this release.\n")
    }

    return releaseNotes.String(), nil
}

func getMergedPRsSince(ctx context.Context, client *github.Client, owner string, repo string, sinceTime *time.Time) ([]*github.PullRequest, error) {
    var allPRs []*github.PullRequest
    page := 1

    for {
        listOptions := &github.PullRequestListOptions{
            State:       "closed",
            Sort:        "updated",
            Direction:   "desc",
            ListOptions: github.ListOptions{Page: page, PerPage: 100},
        }

        prs, resp, err := client.PullRequests.List(ctx, owner, repo, listOptions)
        if err != nil {
            return nil, fmt.Errorf("failed to get pull requests: %v", err)
        }

        foundOldPR := false
        // Filter PRs that were merged after the last release (or all merged PRs if no previous release)
        for _, pr := range prs {
            if pr.MergedAt == nil {
                continue // Skip unmerged PRs
            }

            // If we have a sinceTime and this PR was merged before it, we can stop
            if sinceTime != nil && pr.MergedAt.Before(*sinceTime) {
                foundOldPR = true
                break
            }

            allPRs = append(allPRs, pr)
        }

        // Stop if we found an old PR or reached the end
        if resp.NextPage == 0 || foundOldPR {
            break
        }
        page = resp.NextPage
    }

    return allPRs, nil
}
