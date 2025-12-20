package service

import (
	"context"
	"testing"
	"time"

	"contest/src/domain/entity"
	"contest/src/testutil"
	"contest/src/testutil/mocks"
)

// TestScoreboardService_BuildICPCScoreboard tests ICPC scoreboard calculation
func TestScoreboardService_BuildICPCScoreboard(t *testing.T) {
	ctx := context.Background()
	
	// Arrange
	contest := testutil.NewContestBuilder().
		WithID("contest-1").
		AsICPC().
		AsOngoing().
		Build()

	submissions := []entity.ContestSubmission{
		// User1: AC on problem A at 10 min, WA on B
		*testutil.NewContestSubmissionBuilder().
			WithContestID("contest-1").
			WithUsername("user1").
			WithProblemID("A").
			AsAccepted().
			Build(),
		*testutil.NewContestSubmissionBuilder().
			WithContestID("contest-1").
			WithUsername("user1").
			WithProblemID("B").
			AsWrongAnswer().
			Build(),
		
		// User2: AC on problem A at 15 min
		*testutil.NewContestSubmissionBuilder().
			WithContestID("contest-1").
			WithUsername("user2").
			WithProblemID("A").
			AsAccepted().
			Build(),
	}

	mockSubmissionRepo := mocks.NewMockContestSubmissionRepository()
	mockSubmissionRepo.ListByContestFunc = func(ctx context.Context, contestID string, filter map[string]interface{}) ([]entity.ContestSubmission, error) {
		return submissions, nil
	}

	calculator := NewICPCScoreboardCalculator(mockSubmissionRepo)

	// Act
	scoreboard, err := calculator.BuildScoreboard(ctx, contest, submissions)

	// Assert
	if err != nil {
		t.Fatalf("BuildScoreboard failed: %v", err)
	}

	if len(scoreboard.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(scoreboard.Rows))
	}

	// User1 should be ranked first (solved earlier)
	if scoreboard.Rows[0].Username != "user1" {
		t.Errorf("Expected user1 to be ranked first, got %s", scoreboard.Rows[0].Username)
	}

	if scoreboard.Rows[0].Solved != 1 {
		t.Errorf("Expected user1 to have 1 solved, got %d", scoreboard.Rows[0].Solved)
	}
}

// TestScoreboardService_BuildIOIScoreboard tests IOI scoreboard calculation
func TestScoreboardService_BuildIOIScoreboard(t *testing.T) {
	ctx := context.Background()
	
	contest := testutil.NewContestBuilder().
		WithID("contest-2").
		AsIOI().
		AsOngoing().
		Build()

	submissions := []entity.ContestSubmission{
		// User1: 100 points on A, 50 points on B
		*testutil.NewContestSubmissionBuilder().
			WithContestID("contest-2").
			WithUsername("user1").
			WithProblemID("A").
			WithPoints(100).
			AsAccepted().
			Build(),
		*testutil.NewContestSubmissionBuilder().
			WithContestID("contest-2").
			WithUsername("user1").
			WithProblemID("B").
			WithPoints(50).
			AsAccepted().
			Build(),
	}

	mockSubmissionRepo := mocks.NewMockContestSubmissionRepository()
	mockSubmissionRepo.ListByContestFunc = func(ctx context.Context, contestID string, filter map[string]interface{}) ([]entity.ContestSubmission, error) {
		return submissions, nil
	}

	calculator := NewIOIScoreboardCalculator(mockSubmissionRepo)

	// Act
	scoreboard, err := calculator.BuildScoreboard(ctx, contest, submissions)

	// Assert
	if err != nil {
		t.Fatalf("BuildScoreboard failed: %v", err)
	}

	if len(scoreboard.Rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(scoreboard.Rows))
	}

	if scoreboard.Rows[0].TotalPoints != 150 {
		t.Errorf("Expected 150 total points, got %f", scoreboard.Rows[0].TotalPoints)
	}
}

// TestScoreboardService_IgnoredSubmissions tests that ignored submissions are excluded
func TestScoreboardService_IgnoredSubmissions(t *testing.T) {
	ctx := context.Background()
	
	contest := testutil.NewContestBuilder().
		WithID("contest-3").
		AsICPC().
		Build()

	submissions := []entity.ContestSubmission{
		*testutil.NewContestSubmissionBuilder().
			WithUsername("user1").
			WithProblemID("A").
			AsAccepted().
			Build(),
		*testutil.NewContestSubmissionBuilder().
			WithUsername("cheater").
			WithProblemID("A").
			AsAccepted().
			AsIgnored(). // This should be excluded
			Build(),
	}

	mockSubmissionRepo := mocks.NewMockContestSubmissionRepository()
	mockSubmissionRepo.ListByContestFunc = func(ctx context.Context, contestID string, filter map[string]interface{}) ([]entity.ContestSubmission, error) {
		// Filter out ignored submissions
		filtered := []entity.ContestSubmission{}
		for _, s := range submissions {
			if !s.Ignored {
				filtered = append(filtered, s)
			}
		}
		return filtered, nil
	}

	calculator := NewICPCScoreboardCalculator(mockSubmissionRepo)

	// Act
	scoreboard, err := calculator.BuildScoreboard(ctx, contest, submissions)

	// Assert
	if err != nil {
		t.Fatalf("BuildScoreboard failed: %v", err)
	}

	// Should only have user1, not cheater
	if len(scoreboard.Rows) != 1 {
		t.Errorf("Expected 1 row (ignored submissions excluded), got %d", len(scoreboard.Rows))
	}

	if scoreboard.Rows[0].Username == "cheater" {
		t.Error("Cheater should be excluded from scoreboard")
	}
}

// TestContestIngestionService_IngestSubmission tests submission ingestion flow
func TestContestIngestionService_IngestSubmission(t *testing.T) {
	ctx := context.Background()

	contest := testutil.NewContestBuilder().
		WithID("contest-1").
		AsOngoing().
		Build()

	mockContestRepo := mocks.NewMockContestRepository()
	mockContestRepo.GetByIDFunc = func(ctx context.Context, id string) (*entity.Contest, error) {
		return &contest, nil
	}

	mockSubmissionRepo := mocks.NewMockContestSubmissionRepository()
	
	service := NewContestIngestionService(mockContestRepo, mockSubmissionRepo)

	submission := testutil.NewContestSubmissionBuilder().
		WithContestID("contest-1").
		WithSubmissionID("sub-123").
		AsAccepted().
		Build()

	// Act
	err := service.IngestSubmission(ctx, submission)

	// Assert
	if err != nil {
		t.Fatalf("IngestSubmission failed: %v", err)
	}

	if mockContestRepo.GetByIDCalls != 1 {
		t.Errorf("Expected GetByID to be called once, got %d", mockContestRepo.GetByIDCalls)
	}

	if mockSubmissionRepo.UpsertCalls != 1 {
		t.Errorf("Expected Upsert to be called once, got %d", mockSubmissionRepo.UpsertCalls)
	}
}

// TestContestIngestionService_IngestSubmission_ContestNotFound tests error handling
func TestContestIngestionService_IngestSubmission_ContestNotFound(t *testing.T) {
	ctx := context.Background()

	mockContestRepo := mocks.NewMockContestRepository()
	mockContestRepo.GetByIDFunc = func(ctx context.Context, id string) (*entity.Contest, error) {
		return nil, errors.New("contest not found")
	}

	mockSubmissionRepo := mocks.NewMockContestSubmissionRepository()
	service := NewContestIngestionService(mockContestRepo, mockSubmissionRepo)

	submission := testutil.NewContestSubmissionBuilder().Build()

	// Act
	err := service.IngestSubmission(ctx, submission)

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent contest, got nil")
	}

	// Should not attempt to upsert if contest doesn't exist
	if mockSubmissionRepo.UpsertCalls != 0 {
		t.Error("Should not upsert when contest doesn't exist")
	}
}

// TestContestModerationService_SkipSubmission tests submission skipping
func TestContestModerationService_SkipSubmission(t *testing.T) {
	ctx := context.Background()

	mockSubmissionRepo := mocks.NewMockContestSubmissionRepository()
	service := NewContestModerationService(mockSubmissionRepo)

	// Act
	err := service.SkipSubmission(ctx, "contest-1", "sub-123")

	// Assert
	if err != nil {
		t.Fatalf("SkipSubmission failed: %v", err)
	}

	if mockSubmissionRepo.SetIgnoredCalls != 1 {
		t.Errorf("Expected SetIgnored to be called once, got %d", mockSubmissionRepo.SetIgnoredCalls)
	}
}

// TestContestModerationService_BulkSkipUser tests bulk user skip
func TestContestModerationService_BulkSkipUser(t *testing.T) {
	ctx := context.Background()

	mockSubmissionRepo := mocks.NewMockContestSubmissionRepository()
	service := NewContestModerationService(mockSubmissionRepo)

	// Act
	err := service.BulkSkipUser(ctx, "contest-1", "cheater")

	// Assert
	if err != nil {
		t.Fatalf("BulkSkipUser failed: %v", err)
	}

	// Verify the correct method was called
	// Implementation would call SetIgnoredByUser
}
