package common

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func DownloadFile(ctx context.Context, url, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("Failed to create directory: %v", err)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Failed to create file %s: %v", destPath, err)
	}

	defer out.Close()
	err = SendRequestNoJson(ctx, APIRequest{
		Method:  "GET",
		URL:     url,
		Timeout: 60 * time.Second,
	}, func(body io.Reader) error {
		_, err = io.Copy(out, body)
		if err != nil {
			return fmt.Errorf("Failed to write body to file: %v", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failed to download file from: %s: %v", url, err)
	}

	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("Can't make the downloaded file executable: %v", err)
	}

	return nil
}
