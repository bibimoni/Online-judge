package common

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
)

func DownloadFile(ctx context.Context, url, destPath string) error {
	config.GetLogger().Debug().Msgf("url: %s, path: %s", url, destPath)
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Failed to create file %s: %v", destPath, err)
	}

	defer out.Close()
	resp, err := SendRequestNoJson(ctx, APIRequest{
		Method:  "GET",
		URL:     url,
		Timeout: 60 * time.Second,
	})

	if err != nil {
		os.Remove(destPath)
		return fmt.Errorf("Failed to write body to file: %v", err)
	}
	_, err = io.Copy(out, resp.Body)
	defer resp.Body.Close()

	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("Can't make the downloaded file executable: %v", err)
	}

	return nil
}
