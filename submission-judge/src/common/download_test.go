package common_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
)

func TestDownloadFile(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		config.GetLogger().Panic().Err(err)
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		url      string
		destPath string
		wantErr  bool
	}{
		{
			"test download",
			fmt.Sprintf("%sget/445985/checker", cfg.ProblemServerAddr),
			"/storage/checker",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := common.DownloadFile(context.Background(), tt.url, tt.destPath)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DownloadFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DownloadFile() succeeded unexpectedly")
			}
		})
	}
}
