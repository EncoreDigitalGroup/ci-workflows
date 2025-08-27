package branchname

import (
	"fmt"
	"os"
	"regexp"

	"github.com/EncoreDigitalGroup/ci-workflows/actions/github/formatPullRequestTitle/support/github"
	"github.com/EncoreDigitalGroup/golib/logger"
)

var regexWithIssueType = regexp.MustCompile(`^(epic|feature|bugfix|hotfix)/([A-Z]+-[0-9]+)-(.+)$`)
var regexWithoutIssueType = regexp.MustCompile(`^([A-Z]+-[0-9]+)-(.+)$`)
var pullRequestTitle string

func Format(gh github.GitHub) {
	branchName, err := gh.GetBranchName()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if !gh.BranchNameMatchesPRTitle(branchName) {
		formattedTitle := formatTitle(gh, branchName)
		gh.UpdatePRTitle(formattedTitle)
	}
}

func GetIssueKeyFromBranchName(branchName string) (string, error) {
	if matches := regexWithIssueType.FindStringSubmatch(branchName); matches != nil {
		return matches[2], nil
	} else if matches := regexWithoutIssueType.FindStringSubmatch(branchName); matches != nil {
		return matches[1], nil
	} else {
		fmt.Println("Title does not match expected format")
		logger.Info(pullRequestTitle)
		return "", nil
	}
}

func GetIssueNameFromBranchName(branchName string) (string, error) {
	if matches := regexWithIssueType.FindStringSubmatch(branchName); matches != nil {
		return matches[3], nil
	} else if matches := regexWithoutIssueType.FindStringSubmatch(branchName); matches != nil {
		return matches[2], nil
	} else {
		fmt.Println("Title does not match expected format")
		logger.Info(pullRequestTitle)
		return "", nil
	}
}

func formatTitle(gh github.GitHub, branchName string) string {
	issueKey, err := GetIssueKeyFromBranchName(branchName)
	issueName, err := GetIssueNameFromBranchName(branchName)

	if err != nil {
		fmt.Println("Title does not match expected format")
		logger.Error(err.Error())
		return pullRequestTitle
	}

	return gh.ApplyFormatting(issueKey, issueName)
}
