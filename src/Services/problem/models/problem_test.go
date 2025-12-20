package models

import (
	"testing"
	"time"
)

// TestProblem_Validation tests problem validation
func TestProblem_Validation(t *testing.T) {
	tests := []struct {
		name    string
		problem Problem
		wantErr bool
	}{
		{
			name: "valid_problem",
			problem: Problem{
				Title:       "Two Sum",
				Description: "Find two numbers that add up to target",
				Difficulty:  "easy",
				TimeLimit:   1000,
				MemoryLimit: 256,
			},
			wantErr: false,
		},
		{
			name: "empty_title",
			problem: Problem{
				Title:       "",
				Description: "Description",
				Difficulty:  "easy",
			},
			wantErr: true,
		},
		{
			name: "invalid_difficulty",
			problem: Problem{
				Title:       "Problem",
				Description: "Description",
				Difficulty:  "invalid",
			},
			wantErr: true,
		},
		{
			name: "negative_time_limit",
			problem: Problem{
				Title:       "Problem",
				Description: "Description",
				Difficulty:  "easy",
				TimeLimit:   -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.problem.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTestCase_Validation tests test case validation
func TestTestCase_Validation(t *testing.T) {
	tests := []struct {
		name     string
		testCase TestCase
		wantErr  bool
	}{
		{
			name: "valid_test_case",
			testCase: TestCase{
				Input:  "1 2\n",
				Output: "3\n",
				Points: 10,
			},
			wantErr: false,
		},
		{
			name: "empty_input",
			testCase: TestCase{
				Input:  "",
				Output: "output",
			},
			wantErr: true,
		},
		{
			name: "empty_output",
			testCase: TestCase{
				Input:  "input",
				Output: "",
			},
			wantErr: true,
		},
		{
			name: "negative_points",
			testCase: TestCase{
				Input:  "input",
				Output: "output",
				Points: -5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testCase.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestProblemPackage_Parse tests problem package parsing
func TestProblemPackage_Parse(t *testing.T) {
	// Test parsing a valid problem package structure
	pkg := &ProblemPackage{
		ProblemID: "problem-123",
		Version:   "1.0",
		TestCases: []TestCase{
			{Input: "1 2", Output: "3", Points: 10},
			{Input: "2 3", Output: "5", Points: 10},
		},
	}

	if pkg.TotalPoints() != 20 {
		t.Errorf("Expected total points 20, got %d", pkg.TotalPoints())
	}

	if pkg.TestCaseCount() != 2 {
		t.Errorf("Expected 2 test cases, got %d", pkg.TestCaseCount())
	}
}

// TestDifficulty_Levels tests difficulty validation
func TestDifficulty_Levels(t *testing.T) {
	validDifficulties := []string{"easy", "medium", "hard"}
	
	for _, diff := range validDifficulties {
		if !IsValidDifficulty(diff) {
			t.Errorf("Difficulty %s should be valid", diff)
		}
	}

	if IsValidDifficulty("invalid") {
		t.Error("'invalid' should not be a valid difficulty")
	}
}

// TestProblem_Statistics tests problem statistics
func TestProblem_Statistics(t *testing.T) {
	problem := &Problem{
		TotalSubmissions:    100,
		AcceptedSubmissions: 40,
	}

	acceptanceRate := problem.AcceptanceRate()
	expected := 40.0
	if acceptanceRate != expected {
		t.Errorf("Expected acceptance rate %.2f%%, got %.2f%%", expected, acceptanceRate)
	}
}

// TestProblemTag_Management tests tag operations
func TestProblemTag_Management(t *testing.T) {
	problem := &Problem{
		Tags: []string{"array", "hash-table"},
	}

	// Test adding tag
	problem.AddTag("dynamic-programming")
	if len(problem.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(problem.Tags))
	}

	// Test removing tag
	problem.RemoveTag("hash-table")
	if len(problem.Tags) != 2 {
		t.Errorf("Expected 2 tags after removal, got %d", len(problem.Tags))
	}

	// Test has tag
	if !problem.HasTag("array") {
		t.Error("Should have 'array' tag")
	}
}
