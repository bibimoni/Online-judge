package ratingserviceimpl

import (
	domain "contest/src/domain/entity"
	"math"
	"testing"
)

func TestLogistic(t *testing.T) {
	// Equal ratings → 50% win probability.
	p := logistic(1500, 1500, 200)
	if math.Abs(p-0.5) > 1e-9 {
		t.Errorf("expected 0.5, got %f", p)
	}
	// Higher rating should have > 50% win chance.
	p2 := logistic(1700, 1500, 200)
	if p2 <= 0.5 {
		t.Errorf("expected > 0.5, got %f", p2)
	}
}

func TestComputeRatings_BasicOrdering(t *testing.T) {
	svc := NewEloMMRService(domain.DefaultEloMMRConfig())

	contest := &domain.Contest{}
	contest.EndTime = contest.StartTime.Add(2 * 3600e9)

	standings := []domain.StandingRow{
		{Username: "alice", Rank: 1, Score: 5},
		{Username: "bob", Rank: 2, Score: 3},
		{Username: "charlie", Rank: 3, Score: 1},
	}

	oldStates := map[string]domain.UserRatingState{} // all new players

	results, newStates, err := svc.ComputeRatings(nil, contest, standings, oldStates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// The winner should have a higher new rating than the loser.
	if newStates["alice"].Mu <= newStates["charlie"].Mu {
		t.Errorf("alice (rank 1) should have higher mu than charlie (rank 3): alice=%f charlie=%f",
			newStates["alice"].Mu, newStates["charlie"].Mu)
	}

	// All deltas should be reasonable (new players can have large first-contest swings
	// since display_rating starts at mu - 2*sigma = 800).
	for _, r := range results {
		if math.Abs(r.Delta) > 1500 {
			t.Errorf("unreasonable delta for %s: %f", r.Username, r.Delta)
		}
	}
}

func TestComputeRatings_EmptyStandings(t *testing.T) {
	svc := NewEloMMRService(domain.DefaultEloMMRConfig())
	_, _, err := svc.ComputeRatings(nil, &domain.Contest{}, nil, nil)
	if err == nil {
		t.Error("expected error for empty standings")
	}
}

func TestComputeRatings_ExperiencedPlayers(t *testing.T) {
	svc := NewEloMMRService(domain.DefaultEloMMRConfig())

	contest := &domain.Contest{}

	standings := []domain.StandingRow{
		{Username: "pro", Rank: 1},
		{Username: "newbie", Rank: 2},
	}

	cfg := domain.DefaultEloMMRConfig()
	oldStates := map[string]domain.UserRatingState{
		"pro": {
			Username:      "pro",
			Mu:            2000,
			Sigma:         100,
			DisplayRating: 1800,
			NumContests:   50,
		},
		"newbie": domain.NewUserRatingState("newbie", cfg),
	}

	results, newStates, err := svc.ComputeRatings(nil, contest, standings, oldStates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Pro should still be rated higher.
	if newStates["pro"].Mu <= newStates["newbie"].Mu {
		t.Errorf("pro should remain higher rated")
	}

	// Pro's sigma should be smaller than newbie's.
	if newStates["pro"].Sigma >= newStates["newbie"].Sigma {
		t.Errorf("pro sigma (%f) should be less than newbie sigma (%f)",
			newStates["pro"].Sigma, newStates["newbie"].Sigma)
	}

	_ = results
}

func TestStandingsFromScoreboard(t *testing.T) {
	contest := &domain.Contest{
		Contestants: []domain.Contestant{
			{Username: "alice", ParticipantType: domain.RatedParticipant},
			{Username: "bob", ParticipantType: domain.RatedParticipant},
			{Username: "virtual_user", ParticipantType: domain.VirtualParticipant},
		},
	}

	snapshot := &domain.ScoreboardSnapshot{
		Rows: []domain.ScoreboardRow{
			{Username: "alice", Rank: 1, Score: 5, Penalty: 100},
			{Username: "bob", Rank: 2, Score: 3, Penalty: 200},
			{Username: "virtual_user", Rank: 3, Score: 1, Penalty: 300},
		},
	}

	standings := StandingsFromScoreboard(snapshot, contest)

	if len(standings) != 2 {
		t.Fatalf("expected 2 rated standings, got %d", len(standings))
	}
	if standings[0].Username != "alice" {
		t.Errorf("expected alice at rank 1, got %s", standings[0].Username)
	}
}
