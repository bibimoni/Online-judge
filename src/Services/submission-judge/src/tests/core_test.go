package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestEndToEndJudgingFlow tests the complete judging flow
// This test requires all services to be running:
// - Problem service (port 3000)
// - Auth service (port 50051) via Gateway (port 81)
// - Submission service via Gateway (port 81)
func TestEndToEndJudgingFlow(t *testing.T) {
	// Configuration
	const (
		problemServiceURL = "http://127.0.0.1:3000"
		gatewayURL        = "http://127.0.0.1:81"
		judgeWaitTime     = 60 * time.Second // Wait 1 minute for judging to complete
	)

	// Test data - 4 submissions as specified
	testSubmissions := []TestSubmission{
		{
			Name:      "AC Solution for 445985",
			ProblemID: "445985",
			Code: `/**
 * Author: distiled
 */
#include <bits/stdc++.h>
using namespace std;
 
#ifdef DEBUG
#include </Users/distiled/codeStuff/templates/debug.h>
#else
#define dbg(x...)
#endif
#define int int64_t
 
signed main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int tt;
  cin >> tt;
  while (tt--) {
    int n, r, c;
    cin >> n >> r >> c;
    vector<int> h(n);
    for (int i = 0; i < n; i++)
      cin >> h[i];
    vector<int> w(n);
    for (int i = 0; i < n; i++)
      cin >> w[i];
 
    int ans = 0;
    for (int i = 0; i < n; i++) {
      ans += ((w[i] + r - 1) / r) * ((h[i] + c - 1) / c);
    }
    cout << ans << '\n';
  }
}`,
			Language:        "cpp14",
			SubmissionType:  "ICPC",
			ExpectedVerdict: "ACCEPTED",
		},
		{
			Name:      "MLE Solution for 445985",
			ProblemID: "445985",
			Code: `#include <bits/stdc++.h>
using namespace std;

signed main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  
  // Allocate a huge array to trigger MLE
  // Allocating ~500MB (500 million integers * 8 bytes = 4GB)
  vector<long long> huge_array(500000000, 0);
  
  int tt;
  cin >> tt;
  while (tt--) {
    int n, r, c;
    cin >> n >> r >> c;
    vector<int> h(n);
    for (int i = 0; i < n; i++)
      cin >> h[i];
    vector<int> w(n);
    for (int i = 0; i < n; i++)
      cin >> w[i];
    
    long long ans = 0;
    for (int i = 0; i < n; i++) {
      ans += ((w[i] + r - 1) / r) * ((h[i] + c - 1) / c);
      // Use the huge array to prevent compiler optimization
      huge_array[i % 1000] = ans;
    }
    cout << ans << '\n';
  }
  return 0;
}`,
			Language:        "cpp14",
			SubmissionType:  "ICPC",
			ExpectedVerdict: "MEMORY_LIMIT_EXCEEDED",
		},
		{
			Name:      "TLE Solution for 445985",
			ProblemID: "445985",
			Code: `/**
 * Author: distiled
 */
#include <bits/stdc++.h>
using namespace std;
 
#ifdef DEBUG
#include </Users/distiled/codeStuff/templates/debug.h>
#else
#define dbg(x...)
#endif
#define int int64_t
 
signed main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int tt;
  cin >> tt;
  while (tt--) {
    int n, r, c;
 while(1) {}
   cin >> n >> r >> c;
    vector<int> h(n);
    for (int i = 0; i < n; i++)
      cin >> h[i];
    vector<int> w(n);
    for (int i = 0; i < n; i++)
      cin >> w[i];
 
    int ans = 0;
    for (int i = 0; i < n; i++) {
      ans += ((w[i] + r - 1) / r) * ((h[i] + c - 1) / c);
    }
    cout << ans << '\n';
  }
}`,
			Language:        "cpp14",
			SubmissionType:  "ICPC",
			ExpectedVerdict: "TIME_LIMIT_EXCEEDED",
		},
		{
			Name:      "AC Solution for 440176",
			ProblemID: "440176",
			Code: `/**
 * Author: distiled
 */
#include <bits/stdc++.h>
using namespace std;
 
#ifdef DEBUG
#include </Users/distiled/codeStuff/templates/debug.h>
#else
#define dbg(x...)
#endif
#define int int64_t
 
signed main() {
  ios::sync_with_stdio(false);
  cin.tie(0);
  int tt;
  cin >> tt;
  while (tt--) {
    int n, r, c;
    cin >> n >> r >> c;
    vector<int> h(n);
    for (int i = 0; i < n; i++)
      cin >> h[i];
    vector<int> w(n);
    for (int i = 0; i < n; i++)
      cin >> w[i];
 
    int ans = 0;
    for (int i = 0; i < n; i++) {
      ans += ((w[i] + r - 1) / r) * ((h[i] + c - 1) / c);
    }
    cout << ans << '\n';
  }
}`,
			Language:        "cpp14",
			SubmissionType:  "ICPC",
			ExpectedVerdict: "ACCEPTED",
		},
	}

	var accessToken string
	submissionIDs := make([]string, 0)

	t.Run("1_InitializeProblems", func(t *testing.T) {
		t.Log("Step 1: Initializing problems 445985 and 440176...")
		
		// Initialize problem 445985
		if err := initializeProblem(problemServiceURL, "445985"); err != nil {
			t.Logf("Warning: Failed to initialize problem 445985: %v (may already exist)", err)
		} else {
			t.Log("✓ Problem 445985 initialized successfully")
		}

		// Initialize problem 440176
		if err := initializeProblem(problemServiceURL, "440176"); err != nil {
			t.Logf("Warning: Failed to initialize problem 440176: %v (may already exist)", err)
		} else {
			t.Log("✓ Problem 440176 initialized successfully")
		}

		// Wait a bit for problems to be fully ready
		time.Sleep(3 * time.Second)
	})

	t.Run("2_AuthenticateAdmin", func(t *testing.T) {
		t.Log("Step 2: Authenticating with admin/bkacbkac...")
		
		token, err := authenticate(gatewayURL, "admin", "bkacbkac")
		if err != nil {
			t.Fatalf("Failed to authenticate: %v", err)
		}
		
		accessToken = token
		t.Logf("✓ Authentication successful, token received (length: %d)", len(token))
	})

	t.Run("3_SubmitAllSolutions", func(t *testing.T) {
		t.Log("Step 3: Submitting all 4 test solutions...")
		
		for i, sub := range testSubmissions {
			t.Run(sub.Name, func(t *testing.T) {
				submissionID, err := submitSolution(gatewayURL, sub, accessToken)
				if err != nil {
					t.Fatalf("Failed to submit %s: %v", sub.Name, err)
				}
				
				submissionIDs = append(submissionIDs, submissionID)
				t.Logf("✓ Submitted %s, ID: %s", sub.Name, submissionID)
			})

			// Small delay between submissions
			if i < len(testSubmissions)-1 {
				time.Sleep(500 * time.Millisecond)
			}
		}

		t.Logf("✓ All %d submissions created successfully", len(submissionIDs))
	})

	t.Run("4_WaitForJudging", func(t *testing.T) {
		t.Logf("Step 4: Waiting %v for all submissions to be judged...", judgeWaitTime)
		time.Sleep(judgeWaitTime)
		t.Log("✓ Wait completed")
	})

	t.Run("5_VerifySubmissions", func(t *testing.T) {
		t.Log("Step 5: Verifying all submission results...")
		
		if len(submissionIDs) != len(testSubmissions) {
			t.Fatalf("Mismatch: expected %d submissions but got %d IDs", len(testSubmissions), len(submissionIDs))
		}

		allPassed := true
		for i, submissionID := range submissionIDs {
			sub := testSubmissions[i]
			t.Run(sub.Name, func(t *testing.T) {
				result, err := getSubmissionResult(gatewayURL, submissionID, accessToken)
				if err != nil {
					t.Fatalf("Failed to get result for %s (ID: %s): %v", sub.Name, submissionID, err)
				}

				t.Logf("Submission Details:")
				t.Logf("  Name: %s", sub.Name)
				t.Logf("  ID: %s", submissionID)
				t.Logf("  Problem: %s", sub.ProblemID)
				t.Logf("  Verdict: %s", result.Verdict)
				t.Logf("  Status: %s", result.Status)
				t.Logf("  Message: %s", result.Message)
				
				if result.Time > 0 {
					t.Logf("  Time: %.3fs", result.Time)
				}
				if result.Memory != "" {
					t.Logf("  Memory: %s", result.Memory)
				}

				// Verify the submission was judged (not pending)
				if result.Status == "PENDING" || result.Status == "JUDGING" {
					t.Errorf("❌ Submission is still %s after wait period", result.Status)
					allPassed = false
					return
				}

				// Check for system errors (like cgroup issues)
				if result.Verdict == "COMPILATION_ERROR" && strings.Contains(result.Message, "cgroup") {
					t.Errorf("❌ System error detected - cgroup issue: %s", result.Message)
					t.Errorf("   This indicates the submission-judge container has cgroup configuration problems")
					allPassed = false
					return
				}

				// Verify verdict matches expected (if provided)
				if sub.ExpectedVerdict != "" {
					if result.Verdict != sub.ExpectedVerdict {
						t.Errorf("❌ Verdict mismatch!")
						t.Errorf("   Expected: %s", sub.ExpectedVerdict)
						t.Errorf("   Got:      %s", result.Verdict)
						if result.Message != "" {
							t.Errorf("   Message:  %s", result.Message)
						}
						allPassed = false
					} else {
						t.Logf("✓ Verdict matches expected: %s", result.Verdict)
					}
				}

				// Additional verification based on verdict
				switch sub.ExpectedVerdict {
				case "ACCEPTED":
					if result.Points != result.TotalPoints {
						t.Errorf("❌ AC submission should have full points, got %d/%d", result.Points, result.TotalPoints)
						allPassed = false
					}
					if result.PassedTests != result.TotalTests {
						t.Errorf("❌ AC submission should pass all tests, got %d/%d", result.PassedTests, result.TotalTests)
						allPassed = false
					}
				case "MEMORY_LIMIT_EXCEEDED", "TIME_LIMIT_EXCEEDED":
					if result.Points != 0 {
						t.Errorf("❌ MLE/TLE submission should have 0 points, got %d", result.Points)
						allPassed = false
					}
				}

				if allPassed {
					t.Logf("✓ Submission %s verified successfully", submissionID)
				}
			})
		}

		if allPassed {
			t.Logf("✓ All %d submissions verified successfully", len(submissionIDs))
		} else {
			t.Errorf("❌ Some submissions failed verification")
		}
	})
}

// Helper structures
type TestSubmission struct {
	Name            string
	ProblemID       string
	Code            string
	Language        string
	SubmissionType  string
	ExpectedVerdict string
}

type SubmissionResult struct {
	ID          string  `json:"id"`
	ProblemID   string  `json:"problem_id"`
	UserID      string  `json:"user_id"`
	Code        string  `json:"code"`
	Language    string  `json:"language"`
	Verdict     string  `json:"verdict"`
	Status      string  `json:"status"`
	Time        float64 `json:"time"`
	Memory      string  `json:"memory"`
	Points      int     `json:"points"`
	TotalPoints int     `json:"total_points"`
	PassedTests int     `json:"passed_tests"`
	TotalTests  int     `json:"total_tests"`
	Message     string  `json:"message"`
	CreatedAt   string  `json:"created_at"`
}

type SubmissionResponse struct {
	Data struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	} `json:"data"`
	Success bool `json:"success"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Helper functions

func initializeProblem(baseURL, problemID string) error {
	url := fmt.Sprintf("%s/problem/add?problemId=%s", baseURL, problemID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func authenticate(baseURL, username, password string) (string, error) {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login data: %w", err)
	}

	url := fmt.Sprintf("%s/auth/login", baseURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("authentication failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if authResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in response: %s", string(body))
	}

	return authResp.AccessToken, nil
}

func submitSolution(baseURL string, sub TestSubmission, token string) (string, error) {
	submitData := map[string]string{
		"problem_id":      sub.ProblemID,
		"code":            sub.Code,
		"language":        sub.Language,
		"submission_type": sub.SubmissionType,
	}

	jsonData, err := json.Marshal(submitData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal submission data: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/submission/submit", baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("submission failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var submitResp SubmissionResponse
	if err := json.Unmarshal(body, &submitResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if !submitResp.Success {
		return "", fmt.Errorf("submission unsuccessful: %s", submitResp.Data.Message)
	}

	if submitResp.Data.ID == "" {
		return "", fmt.Errorf("no submission ID in response: %s", string(body))
	}

	return submitResp.Data.ID, nil
}

func getSubmissionResult(baseURL, submissionID, token string) (*SubmissionResult, error) {
	url := fmt.Sprintf("%s/api/v1/submission/view/%s", baseURL, submissionID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get submission: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result SubmissionResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	return &result, nil
}
