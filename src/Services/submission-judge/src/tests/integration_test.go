package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	domain "github.com/bibimoni/Online-judge/submission-judge/src/domain/entitiy"
	isolateservice "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate"
	isolateutils "github.com/bibimoni/Online-judge/submission-judge/src/service/isolate/utils"
	problemutils "github.com/bibimoni/Online-judge/submission-judge/src/service/problem/utils"
)

// TestDownloadFile verifies that DownloadFile downloads content correctly,
// sets permissions to 0755, and cleans up on failure.
func TestDownloadFile(t *testing.T) {
	// 1. Setup Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success content"))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	successPath := filepath.Join(tmpDir, "success_file")
	failPath := filepath.Join(tmpDir, "fail_file")

	// 2. Test Success Case
	err := common.DownloadFile(context.Background(), ts.URL+"/success", successPath)
	if err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}

	// Verify content
	content, err := os.ReadFile(successPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if string(content) != "success content" {
		t.Errorf("Expected 'success content', got '%s'", string(content))
	}

	// Verify permissions (0755)
	info, err := os.Stat(successPath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	// Check executable bit (0111) is set for user/group/others (or at least user)
	// 0755 is -rwxr-xr-x.
	if info.Mode()&0111 == 0 {
		t.Errorf("File is not executable. Mode: %v", info.Mode())
	}

	// 3. Test Failure Cleanup Case
	err = common.DownloadFile(context.Background(), ts.URL+"/fail", failPath)
	if err == nil {
		t.Fatal("Expected error for 404 download, got nil")
	}

	// Verify file does NOT exist
	if _, err := os.Stat(failPath); !os.IsNotExist(err) {
		t.Errorf("File should have been cleaned up after failed download, but it exists")
	}
}

// TestGetFileWithCache verifies the caching logic:
// - Downloads if missing
// - Uses cache if present and fresh
// - Re-downloads if remote is newer
func TestGetFileWithCache(t *testing.T) {
	// Setup Env for Distributed Judging
	os.Setenv("SUBMISSION_IS_MAIN_JUDGE", "false") // Distributed mode
	defer os.Unsetenv("SUBMISSION_IS_MAIN_JUDGE")

	lastModified := time.Now().Add(-1 * time.Hour)
	content := "version 1"

	// Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Last-Modified", lastModified.UTC().Format(http.TimeFormat))
		w.Write([]byte(content))
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	localPath := filepath.Join(tmpDir, "cached_file")

	// 1. First Call: Should Download
	path, err := problemutils.GetFileWithCache(context.Background(), localPath, ts.URL)
	if err != nil {
		t.Fatalf("First GetFileWithCache failed: %v", err)
	}
	if path != localPath {
		t.Errorf("Expected path %s, got %s", localPath, path)
	}

	// Verify content
	data, _ := os.ReadFile(localPath)
	if string(data) != "version 1" {
		t.Errorf("Expected 'version 1', got '%s'", string(data))
	}

	// 2. Second Call: Should Use Cache (No change in Last-Modified)
	// We can't easily verify network calls without spying, but we can verify file mod time or content
	// Let's change content on server but keep Last-Modified same -> Should NOT update
	content = "version 2 (should not download)"

	_, err = problemutils.GetFileWithCache(context.Background(), localPath, ts.URL)
	if err != nil {
		t.Fatalf("Second GetFileWithCache failed: %v", err)
	}

	data, _ = os.ReadFile(localPath)
	if string(data) != "version 1" {
		t.Errorf("Cache was ignored! Got '%s'", string(data))
	}

	// 3. Third Call: Remote is Newer -> Should Download
	lastModified = time.Now().Add(1 * time.Hour) // Newer than local file
	content = "version 3 (new)"

	_, err = problemutils.GetFileWithCache(context.Background(), localPath, ts.URL)
	if err != nil {
		t.Fatalf("Third GetFileWithCache failed: %v", err)
	}

	data, _ = os.ReadFile(localPath)
	if string(data) != "version 3 (new)" {
		t.Errorf("Failed to update cache! Got '%s'", string(data))
	}
}

// TestCopyCheckerPermissions verifies that CopyChecker forces 0755 permissions
// even if the source file is not executable.
func TestCopyCheckerPermissions(t *testing.T) {
	// Override IsolateRoot for testing
	tmpIsolateRoot := t.TempDir() + "/"
	isolateservice.IsolateRoot = tmpIsolateRoot

	// Setup Dummy Isolate
	isolate := &domain.Isolate{ID: 1, Inited: true}
	submissionId := "sub123"

	// Create Isolate Directory Structure
	isolateDir := isolateutils.GetIsolateDir(isolate)
	subDir := filepath.Join(isolateDir, submissionId)
	os.MkdirAll(subDir, 0755)

	// Create Source Checker (Non-executable 0644)
	tmpSourceDir := t.TempDir()
	checkerPath := filepath.Join(tmpSourceDir, "checker")
	os.WriteFile(checkerPath, []byte("checker binary"), 0644)

	// Verify source is NOT executable
	info, _ := os.Stat(checkerPath)
	if info.Mode()&0111 != 0 {
		t.Fatalf("Setup failed: source checker is already executable")
	}

	// Mock Config for CheckerBinName (since CopyChecker uses config.Load)
	// config.Load reads env or defaults. Default is "checker".
	// We rely on default or set env if needed.

	// Run CopyChecker
	err := isolateutils.CopyChecker(isolate, submissionId, checkerPath)
	if err != nil {
		t.Fatalf("CopyChecker failed: %v", err)
	}

	// Verify Destination Permissions
	destPath := filepath.Join(subDir, "checker") // Default name
	destInfo, err := os.Stat(destPath)
	if err != nil {
		t.Fatalf("Failed to stat destination: %v", err)
	}

	// Check if executable (0755)
	if destInfo.Mode()&0111 == 0 {
		t.Errorf("Destination checker is NOT executable! Mode: %v", destInfo.Mode())
	}
}
