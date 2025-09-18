package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestExtractJSONKeys(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected map[string]bool
	}{
		{
			name:     "source.json keys",
			filePath: "test_files/source.json",
			expected: map[string]bool{
				"rootKey1":         true,
				"rootKey1.subKey":  true,
				"rootKey2":         true,
				"rootKey2.subKey1": true,
				"rootKey2.subKey2": true,
			},
		},
		{
			name:     "destination.json keys",
			filePath: "test_files/destination.json",
			expected: map[string]bool{
				"rootKey1":         true,
				"rootKey1.subKey":  true,
				"rootKey2":         true,
				"rootKey2.subKey1": true,
				"rootKey3":         true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := extractJSONKeys(tt.filePath)
			if err != nil {
				t.Fatalf("extractJSONKeys() error = %v", err)
			}

			if !reflect.DeepEqual(keys, tt.expected) {
				t.Errorf("extractJSONKeys() = %v, want %v", keys, tt.expected)
			}
		})
	}
}

func TestCompareJSONKeys(t *testing.T) {
	sourceFile := "test_files/source.json"
	destFiles := []string{"test_files/destination.json"}
	rootDirectory := "."

	newInSource, missingFromSource, warnings, err := compareJSONKeys(sourceFile, destFiles, rootDirectory)
	if err != nil {
		t.Fatalf("compareJSONKeys() error = %v", err)
	}

	// Check for expected new keys in source (missing from destination)
	expectedNewInSource := []KeyDifference{
		{
			Key:                    "rootKey2.subKey2",
			MissingFromDestination: []string{filepath.Join("test_files", "destination.json")},
		},
	}

	if !reflect.DeepEqual(newInSource, expectedNewInSource) {
		t.Errorf("newInSource = %v, want %v", newInSource, expectedNewInSource)
	}

	// Check for expected missing from source (present in destination)
	expectedMissingFromSource := []MissingFromSource{
		{
			File: filepath.Join("test_files", "destination.json"),
			Key:  "rootKey3",
		},
	}

	if !reflect.DeepEqual(missingFromSource, expectedMissingFromSource) {
		t.Errorf("missingFromSource = %v, want %v", missingFromSource, expectedMissingFromSource)
	}

	// Should have no warnings for valid files
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}
}

func TestBuildComment(t *testing.T) {
	sourceFile := "test_files/source.json"
	newInSource := []KeyDifference{
		{
			Key:                    "rootKey2.subKey2",
			MissingFromDestination: []string{"test_files/destination.json"},
		},
	}
	missingFromSource := []MissingFromSource{
		{
			File: "test_files/destination.json",
			Key:  "rootKey3",
		},
	}
	warnings := []ValidationWarning{}
	rootDirectory := "."

	comment := buildComment(sourceFile, newInSource, missingFromSource, warnings, rootDirectory)

	// Check that comment contains expected sections
	if !contains(comment, "JSON Key Differences Report") {
		t.Error("comment should contain 'JSON Key Differences Report'")
	}
	if !contains(comment, "New keys in") {
		t.Error("comment should contain 'New keys in'")
	}
	if !contains(comment, "rootKey2.subKey2") {
		t.Error("comment should contain 'rootKey2.subKey2'")
	}
	if !contains(comment, "rootKey3") {
		t.Error("comment should contain 'rootKey3'")
	}
	if !contains(comment, "test_files/destination.json") {
		t.Error("comment should contain 'test_files/destination.json'")
	}
}

func TestValidateJSONFile(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		shouldErr bool
	}{
		{
			name:      "valid JSON file",
			filePath:  "test_files/source.json",
			shouldErr: false,
		},
		{
			name:      "another valid JSON file",
			filePath:  "test_files/destination.json",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJSONFile(tt.filePath)
			if (err != nil) != tt.shouldErr {
				t.Errorf("validateJSONFile() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}

func TestValidateJSONFileNonExistent(t *testing.T) {
	err := validateJSONFile("non_existent_file.json")
	if err == nil {
		t.Error("validateJSONFile() should return error for non-existent file")
	}
}

func TestIntegrationWithTestFiles(t *testing.T) {
	// Test the core comparison logic
	sourceFile := filepath.Join(".", "test_files/source.json")
	destFiles := []string{filepath.Join(".", "test_files/destination.json")}

	newInSource, missingFromSource, warnings, err := compareJSONKeys(sourceFile, destFiles, ".")
	if err != nil {
		t.Fatalf("compareJSONKeys() error = %v", err)
	}

	// Verify the differences are detected correctly
	if len(newInSource) != 1 {
		t.Errorf("expected 1 new key in source, got %d", len(newInSource))
	}
	if len(missingFromSource) != 1 {
		t.Errorf("expected 1 missing key from source, got %d", len(missingFromSource))
	}
	if len(warnings) != 0 {
		t.Errorf("expected 0 warnings, got %d", len(warnings))
	}

	// Test comment generation and write to markdown file
	comment := buildComment(sourceFile, newInSource, missingFromSource, warnings, ".")
	if comment == "" {
		t.Error("comment should not be empty")
	}

	// Write comment to markdown file for verification
	outputFile := filepath.Join("test_files", "actual_output.md")
	err = os.WriteFile(outputFile, []byte(comment), 0644)
	if err != nil {
		t.Fatalf("failed to write output file: %v", err)
	}

	// Read expected output
	expectedFile := filepath.Join("test_files", "expected_output.md")
	expectedContent, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read expected output file: %v", err)
	}

	// Read actual output
	actualContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read actual output file: %v", err)
	}

	// Normalize line separators and compare
	expectedNormalized := normalizeLineEndings(string(expectedContent))
	actualNormalized := normalizeLineEndings(string(actualContent))

	if expectedNormalized != actualNormalized {
		t.Errorf("output does not match expected content.\nExpected:\n%s\nActual:\n%s", expectedNormalized, actualNormalized)
		return
	}

	// Clean up actual output file after successful test
	os.Remove(outputFile)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// normalizeLineEndings converts all line endings to \n for comparison
func normalizeLineEndings(s string) string {
	// Replace Windows line endings (\r\n) with Unix line endings (\n)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	// Replace old Mac line endings (\r) with Unix line endings (\n)
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}
