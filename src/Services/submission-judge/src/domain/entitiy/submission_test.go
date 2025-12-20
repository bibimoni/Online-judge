package domain

import (
	"testing"
	"time"
)

// TestSubmission_Validation tests submission validation rules
func TestSubmission_Validation(t *testing.T) {
	tests := []struct {
		name        string
		submission  Submission
		wantErr     bool
		errContains string
	}{
		{
			name: "valid_submission",
			submission: Submission{
				ProblemID: "problem-123",
				Language:  "cpp",
				Code:      "int main() { return 0; }",
				UserID:    "user-123",
			},
			wantErr: false,
		},
		{
			name: "empty_problem_id",
			submission: Submission{
				ProblemID: "",
				Language:  "cpp",
				Code:      "code",
				UserID:    "user-123",
			},
			wantErr:     true,
			errContains: "problem ID is required",
		},
		{
			name: "invalid_language",
			submission: Submission{
				ProblemID: "problem-123",
				Language:  "invalid-lang",
				Code:      "code",
				UserID:    "user-123",
			},
			wantErr:     true,
			errContains: "unsupported language",
		},
		{
			name: "empty_code",
			submission: Submission{
				ProblemID: "problem-123",
				Language:  "cpp",
				Code:      "",
				UserID:    "user-123",
			},
			wantErr:     true,
			errContains: "code is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.submission.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestJudgeResult_CalculateScore tests score calculation
func TestJudgeResult_CalculateScore(t *testing.T) {
	tests := []struct {
		name           string
		testCasesPassed int
		totalTestCases int
		wantScore      float64
	}{
		{"all_passed", 10, 10, 100.0},
		{"half_passed", 5, 10, 50.0},
		{"none_passed", 0, 10, 0.0},
		{"partial", 7, 10, 70.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &JudgeResult{
				TestCasesPassed: tt.testCasesPassed,
				TotalTestCases:  tt.totalTestCases,
			}
			score := result.CalculateScore()
			if score != tt.wantScore {
				t.Errorf("CalculateScore() = %v, want %v", score, tt.wantScore)
			}
		})
	}
}

// TestEvaluationStatus_Transitions tests status transitions
func TestEvaluationStatus_Transitions(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus EvalStatus
		toStatus   EvalStatus
		wantErr    bool
	}{
		{"pending_to_judging", StatusPending, StatusJudging, false},
		{"judging_to_finished", StatusJudging, StatusFinished, false},
		{"invalid_finished_to_pending", StatusFinished, StatusPending, true},
		{"invalid_judging_to_pending", StatusJudging, StatusPending, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := &Evaluation{Status: tt.fromStatus}
			err := eval.TransitionTo(tt.toStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitionTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestIsolatePool_Management tests isolate pool operations
func TestIsolatePool_Management(t *testing.T) {
	pool := NewIsolatePool(3)

	// Test acquiring isolates
	for i := 0; i < 3; i++ {
		isolate, err := pool.Acquire()
		if err != nil {
			t.Fatalf("Failed to acquire isolate %d: %v", i, err)
		}
		if isolate == nil {
			t.Error("Acquired isolate is nil")
		}
	}

	// Pool should be exhausted
	_, err := pool.AcquireWithTimeout(100 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error when pool is exhausted")
	}

	// Release one isolate
	isolate, _ := pool.Acquire()
	pool.Release(isolate)

	// Should be able to acquire again
	_, err = pool.Acquire()
	if err != nil {
		t.Error("Failed to acquire after release")
	}
}

// TestBatchEvaluation tests batch evaluation logic
func TestBatchEvaluation(t *testing.T) {
	batch := &Batch{
		Submissions: []string{"sub1", "sub2", "sub3"},
		Status:      BatchStatusPending,
	}

	// Test batch processing
	if batch.TotalSubmissions() != 3 {
		t.Errorf("Expected 3 submissions, got %d", batch.TotalSubmissions())
	}

	// Process submissions
	batch.MarkProcessed("sub1")
	batch.MarkProcessed("sub2")

	if batch.ProcessedCount() != 2 {
		t.Errorf("Expected 2 processed, got %d", batch.ProcessedCount())
	}

	if batch.IsComplete() {
		t.Error("Batch should not be complete yet")
	}

	batch.MarkProcessed("sub3")
	if !batch.IsComplete() {
		t.Error("Batch should be complete")
	}
}
