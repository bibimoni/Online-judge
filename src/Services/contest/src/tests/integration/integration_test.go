package integration

import (
	"context"
	"testing"
	"time"

	"contest/src/domain/entity"
	"contest/src/infrastructure/database"
	"contest/src/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestContainerSetup holds test containers for integration tests
type TestContainerSetup struct {
	MongoContainer *mongodb.MongoDBContainer
	RedisContainer *redis.RedisContainer
	MongoClient    *mongo.Client
	RedisClient    *redis.Client
	Cleanup        func()
}

// SetupTestContainers initializes MongoDB and Redis test containers
func SetupTestContainers(t *testing.T) *TestContainerSetup {
	ctx := context.Background()

	// Start MongoDB container
	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:7"),
		mongodb.WithUsername(""),
		mongodb.WithPassword(""),
	)
	require.NoError(t, err, "Failed to start MongoDB container")

	mongoURI, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get MongoDB connection string")

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err, "Failed to connect to MongoDB")

	// Start Redis container
	redisContainer, err := redis.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
	)
	require.NoError(t, err, "Failed to start Redis container")

	redisURI, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get Redis connection string")

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURI,
	})

	cleanup := func() {
		if mongoClient != nil {
			mongoClient.Disconnect(ctx)
		}
		if redisClient != nil {
			redisClient.Close()
		}
		if mongoContainer != nil {
			mongoContainer.Terminate(ctx)
		}
		if redisContainer != nil {
			redisContainer.Terminate(ctx)
		}
	}

	return &TestContainerSetup{
		MongoContainer: mongoContainer,
		RedisContainer: redisContainer,
		MongoClient:    mongoClient,
		RedisClient:    redisClient,
		Cleanup:        cleanup,
	}
}

// TestIntegration_CreateContest_EndToEnd tests full contest creation flow
func TestIntegration_CreateContest_EndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup
	setup := SetupTestContainers(t)
	defer setup.Cleanup()

	ctx := context.Background()
	db := setup.MongoClient.Database("contest-test-db")

	// Initialize repository
	contestRepo := database.NewMongoContestRepository(db)

	// Create a test contest
	contest := testutil.NewContestBuilder().
		WithName("Integration Test Contest").
		WithDescription("Testing full flow").
		Build()

	// Act - Create contest
	contestID, err := contestRepo.Create(ctx, contest)

	// Assert
	require.NoError(t, err, "Failed to create contest")
	assert.NotEmpty(t, contestID, "Contest ID should not be empty")

	// Verify contest was saved
	retrieved, err := contestRepo.GetByID(ctx, contestID)
	require.NoError(t, err, "Failed to retrieve contest")
	assert.Equal(t, contest.Name, retrieved.Name)
	assert.Equal(t, contest.Description, retrieved.Description)
}

// TestIntegration_SubmissionIngestion_WithScoreboardUpdate tests the full submission flow
func TestIntegration_SubmissionIngestion_WithScoreboardUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	setup := SetupTestContainers(t)
	defer setup.Cleanup()

	ctx := context.Background()
	db := setup.MongoClient.Database("contest-test-db")

	// Initialize repositories
	contestRepo := database.NewMongoContestRepository(db)
	submissionRepo := database.NewMongoContestSubmissionRepository(db)

	// Create a contest
	contest := testutil.NewContestBuilder().
		AsOngoing().
		Build()

	contestID, err := contestRepo.Create(ctx, contest)
	require.NoError(t, err)

	// Submit a solution
	submission := testutil.NewContestSubmissionBuilder().
		WithContestID(contestID).
		WithUsername("testuser").
		WithProblemID("A").
		AsAccepted().
		Build()

	err = submissionRepo.UpsertFromJudgeEvent(ctx, *submission)
	require.NoError(t, err, "Failed to ingest submission")

	// Verify submission was saved
	submissions, err := submissionRepo.ListByContest(ctx, contestID, map[string]interface{}{})
	require.NoError(t, err)
	assert.Len(t, submissions, 1, "Should have 1 submission")
	assert.Equal(t, "testuser", submissions[0].Username)
	assert.Equal(t, entity.VerdictAccepted, submissions[0].Verdict)
}

// TestIntegration_Rejudge_Flow tests the rejudge flow end-to-end
func TestIntegration_Rejudge_Flow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	setup := SetupTestContainers(t)
	defer setup.Cleanup()

	ctx := context.Background()
	db := setup.MongoClient.Database("contest-test-db")

	contestRepo := database.NewMongoContestRepository(db)
	submissionRepo := database.NewMongoContestSubmissionRepository(db)
	rejudgeRepo := database.NewMongoRejudgeJobRepository(db)

	// Create contest
	contest := testutil.NewContestBuilder().AsEnded().Build()
	contestID, err := contestRepo.Create(ctx, contest)
	require.NoError(t, err)

	// Create initial submissions
	sub1 := testutil.NewContestSubmissionBuilder().
		WithContestID(contestID).
		WithSubmissionID("sub-1").
		AsAccepted().
		Build()

	err = submissionRepo.UpsertFromJudgeEvent(ctx, *sub1)
	require.NoError(t, err)

	// Create rejudge job
	rejudgeJob := testutil.NewRejudgeJobBuilder().
		WithContestID(contestID).
		WithSubmissionIDs("sub-1").
		Build()

	jobID, err := rejudgeRepo.Create(ctx, rejudgeJob)
	require.NoError(t, err)

	// Simulate rejudge result - verdict changes to WA
	err = submissionRepo.UpdateVerdictPointsForRejudge(ctx, contestID, "sub-1", entity.VerdictWrongAnswer, 0)
	require.NoError(t, err)

	// Verify updated verdict
	submissions, err := submissionRepo.ListByContest(ctx, contestID, map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, entity.VerdictWrongAnswer, submissions[0].Verdict)

	// Update rejudge job status
	err = rejudgeRepo.UpdateStatus(ctx, jobID, entity.RejudgeStatusCompleted)
	require.NoError(t, err)
}

// TestIntegration_Moderation_SkipUser tests skipping a user's submissions
func TestIntegration_Moderation_SkipUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	setup := SetupTestContainers(t)
	defer setup.Cleanup()

	ctx := context.Background()
	db := setup.MongoClient.Database("contest-test-db")

	contestRepo := database.NewMongoContestRepository(db)
	submissionRepo := database.NewMongoContestSubmissionRepository(db)

	// Create contest
	contest := testutil.NewContestBuilder().Build()
	contestID, err := contestRepo.Create(ctx, contest)
	require.NoError(t, err)

	// Create submissions for a cheater
	for i := 0; i < 3; i++ {
		sub := testutil.NewContestSubmissionBuilder().
			WithContestID(contestID).
			WithUsername("cheater").
			AsAccepted().
			Build()
		err = submissionRepo.UpsertFromJudgeEvent(ctx, *sub)
		require.NoError(t, err)
	}

	// Skip all submissions by user
	err = submissionRepo.SetIgnoredByUser(ctx, contestID, "cheater", true)
	require.NoError(t, err)

	// Verify all submissions are ignored
	submissions, err := submissionRepo.ListByContest(ctx, contestID, map[string]interface{}{
		"username": "cheater",
	})
	require.NoError(t, err)

	for _, sub := range submissions {
		assert.True(t, sub.Ignored, "Submission should be ignored")
	}
}

// TestIntegration_ScoreboardSnapshot_LiveAndFinal tests scoreboard snapshots
func TestIntegration_ScoreboardSnapshot_LiveAndFinal(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	setup := SetupTestContainers(t)
	defer setup.Cleanup()

	ctx := context.Background()
	
	// Use Redis for live scoreboard
	scoreboardRepo := database.NewRedisScoreboardRepository(setup.RedisClient)
	
	// Create a scoreboard snapshot
	scoreboard := testutil.NewScoreboardBuilder().
		WithContestID("contest-1").
		AsLive().
		Build()

	// Save live snapshot
	err := scoreboardRepo.SaveSnapshot(ctx, "contest-1", "live", scoreboard)
	require.NoError(t, err)

	// Retrieve live snapshot
	retrieved, err := scoreboardRepo.GetSnapshot(ctx, "contest-1", "live")
	require.NoError(t, err)
	assert.Equal(t, entity.SnapshotKindLive, retrieved.Kind)

	// Save final snapshot
	scoreboard.Kind = entity.SnapshotKindFinal
	err = scoreboardRepo.SaveSnapshot(ctx, "contest-1", "final", scoreboard)
	require.NoError(t, err)

	// Both should exist
	liveSnapshot, err := scoreboardRepo.GetSnapshot(ctx, "contest-1", "live")
	require.NoError(t, err)
	assert.Equal(t, entity.SnapshotKindLive, liveSnapshot.Kind)

	finalSnapshot, err := scoreboardRepo.GetSnapshot(ctx, "contest-1", "final")
	require.NoError(t, err)
	assert.Equal(t, entity.SnapshotKindFinal, finalSnapshot.Kind)
}

// BenchmarkScoreboardCalculation benchmarks scoreboard calculation performance
func BenchmarkScoreboardCalculation(b *testing.B) {
	ctx := context.Background()
	
	// Create a contest with many submissions
	contest := testutil.NewContestBuilder().AsICPC().Build()
	
	// Generate 1000 submissions from 100 users
	submissions := make([]entity.ContestSubmission, 1000)
	for i := 0; i < 1000; i++ {
		submissions[i] = *testutil.NewContestSubmissionBuilder().
			WithUsername("user" + string(rune(i%100))).
			WithProblemID("problem" + string(rune(i%10))).
			Build()
	}

	calculator := NewICPCScoreboardCalculator(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calculator.BuildScoreboard(ctx, contest, submissions)
	}
}
