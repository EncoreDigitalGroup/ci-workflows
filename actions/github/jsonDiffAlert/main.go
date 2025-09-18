package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
)

const envGHToken = "GH_TOKEN"
const envGHRepository = "GH_REPOSITORY"
const envPRNumber = "PR_NUMBER"
const envSourceFile = "SOURCE_FILE"
const envDestinationFiles = "DESTINATION_FILES"
const envRootDirectory = "ROOT_DIRECTORY"
const envSuppressNoChanges = "SUPPRESS_NO_CHANGES"

type KeyDifference struct {
	Key                    string
	MissingFromDestination []string
}

type MissingFromSource struct {
	File string
	Key  string
}

type ValidationWarning struct {
	File   string
	Reason string
}

func main() {
	githubToken := os.Getenv(envGHToken)
	repo := os.Getenv(envGHRepository)
	prNumberStr := os.Getenv(envPRNumber)
	sourceFile := os.Getenv(envSourceFile)
	destinationFiles := os.Getenv(envDestinationFiles)
	rootDirectory := os.Getenv(envRootDirectory)
	suppressNoChanges := os.Getenv(envSuppressNoChanges) == "true"

	// Default to GITHUB_WORKSPACE if no root directory is provided
	if rootDirectory == "" {
		rootDirectory = os.Getenv("GITHUB_WORKSPACE")
		if rootDirectory == "" {
			rootDirectory = "."
		}
	}

	validateEnv(envGHToken, githubToken)
	validateEnv(envGHRepository, repo)
	validateEnv(envPRNumber, prNumberStr)
	validateEnv(envSourceFile, destinationFiles)
	validateEnv(envDestinationFiles, destinationFiles)

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		fmt.Printf("::error::%s must be in the format owner/repo\n", envGHRepository)
		os.Exit(1)
	}
	repoOwner := parts[0]
	repoName := parts[1]

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		fmt.Printf("::error::%s is not a valid integer: %v\n", envPRNumber, err)
		os.Exit(1)
	}

	destFiles := strings.Split(destinationFiles, ",")
	for i, file := range destFiles {
		destFiles[i] = strings.TrimSpace(file)
	}

	// Resolve file paths relative to root directory
	resolvedSourceFile := filepath.Join(rootDirectory, sourceFile)
	var resolvedDestFiles []string
	for _, file := range destFiles {
		resolvedDestFiles = append(resolvedDestFiles, filepath.Join(rootDirectory, file))
	}

	newInSource, missingFromSource, warnings, err := compareJSONKeys(resolvedSourceFile, resolvedDestFiles, rootDirectory)
	if err != nil {
		fmt.Printf("::error::Failed to compare JSON keys: %v\n", err)
		os.Exit(1)
	}

	if len(newInSource) == 0 && len(missingFromSource) == 0 && len(warnings) == 0 {
		fmt.Println("No JSON key differences found")
		return
	}

	if suppressNoChanges && len(newInSource) == 0 && len(missingFromSource) == 0 && len(warnings) == 0 {
		fmt.Println("No JSON key differences found - comment suppressed")
		return
	}

	comment := buildComment(sourceFile, newInSource, missingFromSource, warnings, rootDirectory)
	err = postComment(githubToken, repoOwner, repoName, prNumber, comment)
	if err != nil {
		fmt.Printf("::error::Failed to post comment: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("JSON diff comment posted successfully")
}

func validateEnv(envVar string, envVal string) {
	if envVal == "" {
		fmt.Printf("::error::%s environment variable is not set\n", envVar)
		os.Exit(1)
	}
}

func isJSONFile(filePath string) bool {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".json"
}

func validateJSONFile(filePath string) error {
	if !isJSONFile(filePath) {
		return fmt.Errorf("file is not a JSON file (extension: %s)", filepath.Ext(filePath))
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

func extractJSONKeys(filePath string) (map[string]bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON in %s: %w", filePath, err)
	}

	keys := make(map[string]bool)
	extractKeysRecursive("", jsonData, keys)
	return keys, nil
}

func extractKeysRecursive(prefix string, data interface{}, keys map[string]bool) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}
			keys[fullKey] = true
			extractKeysRecursive(fullKey, value, keys)
		}
	case []interface{}:
		for i, item := range v {
			indexKey := fmt.Sprintf("%s[%d]", prefix, i)
			extractKeysRecursive(indexKey, item, keys)
		}
	}
}

func compareJSONKeys(sourceFile string, destFiles []string, rootDirectory string) ([]KeyDifference, []MissingFromSource, []ValidationWarning, error) {
	var warnings []ValidationWarning

	// Helper function to get relative path for display
	getDisplayPath := func(fullPath string) string {
		relPath, err := filepath.Rel(rootDirectory, fullPath)
		if err != nil {
			return filepath.Base(fullPath) // fallback to basename
		}
		return relPath
	}

	// Validate and extract source file keys
	var sourceKeys map[string]bool
	if err := validateJSONFile(sourceFile); err != nil {
		warnings = append(warnings, ValidationWarning{
			File:   getDisplayPath(sourceFile),
			Reason: fmt.Sprintf("Source file validation failed: %s", err.Error()),
		})
		sourceKeys = make(map[string]bool) // Use empty keys if validation fails
	} else {
		var err error
		sourceKeys, err = extractJSONKeys(sourceFile)
		if err != nil {
			warnings = append(warnings, ValidationWarning{
				File:   getDisplayPath(sourceFile),
				Reason: fmt.Sprintf("Failed to extract keys from source: %s", err.Error()),
			})
			sourceKeys = make(map[string]bool) // Use empty keys if extraction fails
		}
	}

	destKeysMap := make(map[string]map[string]bool)
	for _, destFile := range destFiles {
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			warnings = append(warnings, ValidationWarning{
				File:   getDisplayPath(destFile),
				Reason: "File does not exist",
			})
			continue
		}

		if err := validateJSONFile(destFile); err != nil {
			warnings = append(warnings, ValidationWarning{
				File:   getDisplayPath(destFile),
				Reason: fmt.Sprintf("Invalid JSON file: %s", err.Error()),
			})
			continue
		}

		destKeys, err := extractJSONKeys(destFile)
		if err != nil {
			warnings = append(warnings, ValidationWarning{
				File:   getDisplayPath(destFile),
				Reason: fmt.Sprintf("Failed to extract keys: %s", err.Error()),
			})
			continue
		}
		destKeysMap[destFile] = destKeys
	}

	var newInSource []KeyDifference
	var missingFromSource []MissingFromSource

	// Find keys in source that are missing from destination files
	for sourceKey := range sourceKeys {
		missingFromDest := []string{}
		for destFile, destKeys := range destKeysMap {
			if !destKeys[sourceKey] {
				missingFromDest = append(missingFromDest, getDisplayPath(destFile))
			}
		}
		if len(missingFromDest) > 0 {
			sort.Strings(missingFromDest)
			newInSource = append(newInSource, KeyDifference{
				Key:                    sourceKey,
				MissingFromDestination: missingFromDest,
			})
		}
	}

	// Find keys in destination files that are missing from source
	for destFile, destKeys := range destKeysMap {
		for destKey := range destKeys {
			if !sourceKeys[destKey] {
				missingFromSource = append(missingFromSource, MissingFromSource{
					File: getDisplayPath(destFile),
					Key:  destKey,
				})
			}
		}
	}

	// Sort results for consistent output
	sort.Slice(newInSource, func(i, j int) bool {
		return newInSource[i].Key < newInSource[j].Key
	})
	sort.Slice(missingFromSource, func(i, j int) bool {
		if missingFromSource[i].File == missingFromSource[j].File {
			return missingFromSource[i].Key < missingFromSource[j].Key
		}
		return missingFromSource[i].File < missingFromSource[j].File
	})

	return newInSource, missingFromSource, warnings, nil
}

func buildComment(sourceFile string, newInSource []KeyDifference, missingFromSource []MissingFromSource, warnings []ValidationWarning, rootDirectory string) string {
	var comment strings.Builder

	comment.WriteString("## JSON Key Differences Report\n\n")

	// Get relative path for source file display
	sourceFileName, err := filepath.Rel(rootDirectory, sourceFile)
	if err != nil {
		sourceFileName = filepath.Base(sourceFile) // fallback to basename
	}

	if len(newInSource) > 0 {
		comment.WriteString(fmt.Sprintf("### ⚠️ New keys in `%s`:\n\n", sourceFileName))

		// Group by key for better readability
		for _, diff := range newInSource {
			comment.WriteString(fmt.Sprintf("**%s:**\n", diff.Key))
			for _, file := range diff.MissingFromDestination {
				comment.WriteString(fmt.Sprintf("- Missing from: %s\n", file))
			}
			comment.WriteString("\n")
		}
	}

	if len(missingFromSource) > 0 {
		comment.WriteString(fmt.Sprintf("###⚠️ Keys present in destination files but missing from `%s`:\n\n", sourceFileName))

		// Group by file for better readability
		fileGroups := make(map[string][]string)
		for _, missing := range missingFromSource {
			fileGroups[missing.File] = append(fileGroups[missing.File], missing.Key)
		}

		var files []string
		for file := range fileGroups {
			files = append(files, file)
		}
		sort.Strings(files)

		for _, file := range files {
			comment.WriteString(fmt.Sprintf("**%s:**\n", file))
			sort.Strings(fileGroups[file])
			for _, key := range fileGroups[file] {
				comment.WriteString(fmt.Sprintf("- `%s`\n", key))
			}
			comment.WriteString("\n")
		}
	}

	if len(warnings) > 0 {
		comment.WriteString("### ⚠️ File Processing Warnings:\n\n")
		for _, warning := range warnings {
			comment.WriteString(fmt.Sprintf("- **%s**: %s\n", warning.File, warning.Reason))
		}
		comment.WriteString("\n")
	}

	if len(newInSource) == 0 && len(missingFromSource) == 0 {
		if len(warnings) == 0 {
			comment.WriteString("✅ No key differences found between source and destination files.")
		} else {
			comment.WriteString("❌ Analysis completed with warnings.")
		}
	}

	return comment.String()
}

func postComment(token, owner, repo string, prNumber int, body string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	comment := &github.IssueComment{
		Body: &body,
	}

	_, _, err := client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	return err
}
