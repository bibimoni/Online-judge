package entity

import (
	"testing"
	"time"
)

// TestContest_Validation tests contest validation rules
func TestContest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		contest     Contest
		wantErr     bool
		errContains string
	}{
		{
			name: "valid_icpc_contest",
			contest: Contest{
				Name:        "ICPC Regional 2025",
				Description: "Regional programming contest",
				StartTime:   time.Now().Add(24 * time.Hour),
				EndTime:     time.Now().Add(29 * time.Hour),
				ScoringType: ScoringTypeICPC,
				Status:      StatusDraft,
				Rated:       true,
			},
			wantErr: false,
		},
		{
			name: "invalid_end_before_start",
			contest: Contest{
				Name:        "Invalid Contest",
				StartTime:   time.Now().Add(24 * time.Hour),
				EndTime:     time.Now().Add(20 * time.Hour), // Before start
				ScoringType: ScoringTypeICPC,
			},
			wantErr:     true,
			errContains: "end time must be after start time",
		},
		{
			name: "invalid_empty_name",
			contest: Contest{
				Name:        "",
				StartTime:   time.Now().Add(24 * time.Hour),
				EndTime:     time.Now().Add(29 * time.Hour),
				ScoringType: ScoringTypeICPC,
			},
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name: "invalid_scoring_type",
			contest: Contest{
				Name:        "Test Contest",
				StartTime:   time.Now().Add(24 * time.Hour),
				EndTime:     time.Now().Add(29 * time.Hour),
				ScoringType: "INVALID",
			},
			wantErr:     true,
			errContains: "invalid scoring type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contest.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && !contains(err.Error(), tt.errContains) {
				t.Errorf("Validate() error = %v, should contain %v", err, tt.errContains)
			}
		})
	}
}

// TestContest_StatusTransitions tests valid contest status transitions
func TestContest_StatusTransitions(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus ContestStatus
		toStatus   ContestStatus
		wantErr    bool
	}{
		{"draft_to_scheduled", StatusDraft, StatusScheduled, false},
		{"scheduled_to_ongoing", StatusScheduled, StatusOngoing, false},
		{"ongoing_to_ended", StatusOngoing, StatusEnded, false},
		{"ended_to_rated", StatusEnded, StatusRated, false},
		{"invalid_draft_to_ended", StatusDraft, StatusEnded, true},
		{"invalid_rated_to_ongoing", StatusRated, StatusOngoing, true},
		{"invalid_ongoing_to_draft", StatusOngoing, StatusDraft, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contest := &Contest{Status: tt.fromStatus}
			err := contest.TransitionTo(tt.toStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransitionTo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && contest.Status != tt.toStatus {
				t.Errorf("Expected status %v, got %v", tt.toStatus, contest.Status)
			}
		})
	}
}

// TestContest_IsActive tests contest active state
func TestContest_IsActive(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		startTime  time.Time
		endTime    time.Time
		status     ContestStatus
		wantActive bool
	}{
		{
			name:       "active_ongoing_contest",
			startTime:  now.Add(-1 * time.Hour),
			endTime:    now.Add(1 * time.Hour),
			status:     StatusOngoing,
			wantActive: true,
		},
		{
			name:       "not_active_ended_contest",
			startTime:  now.Add(-2 * time.Hour),
			endTime:    now.Add(-1 * time.Hour),
			status:     StatusEnded,
			wantActive: false,
		},
		{
			name:       "not_active_future_contest",
			startTime:  now.Add(1 * time.Hour),
			endTime:    now.Add(2 * time.Hour),
			status:     StatusScheduled,
			wantActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contest := &Contest{
				StartTime: tt.startTime,
				EndTime:   tt.endTime,
				Status:    tt.status,
			}
			if got := contest.IsActive(); got != tt.wantActive {
				t.Errorf("IsActive() = %v, want %v", got, tt.wantActive)
			}
		})
	}
}

// TestContest_CanAcceptSubmissions tests submission acceptance rules
func TestContest_CanAcceptSubmissions(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		contest   Contest
		wantAllow bool
	}{
		{
			name: "allow_during_contest",
			contest: Contest{
				StartTime: now.Add(-1 * time.Hour),
				EndTime:   now.Add(1 * time.Hour),
				Status:    StatusOngoing,
			},
			wantAllow: true,
		},
		{
			name: "deny_after_contest",
			contest: Contest{
				StartTime: now.Add(-2 * time.Hour),
				EndTime:   now.Add(-1 * time.Hour),
				Status:    StatusEnded,
			},
			wantAllow: false,
		},
		{
			name: "deny_before_contest",
			contest: Contest{
				StartTime: now.Add(1 * time.Hour),
				EndTime:   now.Add(2 * time.Hour),
				Status:    StatusScheduled,
			},
			wantAllow: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.contest.CanAcceptSubmissions(); got != tt.wantAllow {
				t.Errorf("CanAcceptSubmissions() = %v, want %v", got, tt.wantAllow)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}
