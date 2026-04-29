package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bibimoni/Online-judge/submission-judge/src/common"
	"github.com/bibimoni/Online-judge/submission-judge/src/infrastructure/config"
)

func FileExsits(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("An error occured when trying to verify if a file exists")
	}
}

func EnsureProblemDirectory(stringAddr string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.Judge.IsMain {
		stat, err := FileExsits(stringAddr)
		if err != nil {
			return err
		}
		if !stat {
			return fmt.Errorf("This directory isn't available")
		}
	} else {
		if err := os.MkdirAll(stringAddr, 0755); err != nil {
			return err
		}
	}
	return nil
}

func CheckIfRemoteIsNewer(
	ctx context.Context,
	remoteURL string,
	localModTime time.Time,
) (bool, error) {
	resq, err := common.SendRequestNoJson(ctx, common.APIRequest{
		Method:  "HEAD",
		URL:     remoteURL,
		Timeout: 30 * time.Second,
	})

	lastModifiedStr := resq.Header.Get("Last-Modified")
	config.GetLogger().Debug().Msgf("Last modified value: %s, current time: %v", lastModifiedStr, localModTime)

	// Assume not newer just to be safe
	if lastModifiedStr == "" {
		config.GetLogger().Warn().Msgf("Can't find Last-Modified field in the header")
		return false, nil
	}

	remoteModTime, err := http.ParseTime(lastModifiedStr)
	if err != nil {
		return false, err
	}

	return remoteModTime.After(localModTime), nil
}

func GetFileWithCache(
	ctx context.Context,
	localPath string,
	remoteURL string,
) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}

	info, err := os.Stat(localPath)
	exists := err == nil
	if cfg.Judge.IsMain {
		if exists {
			return localPath, nil
		}
		return "", fmt.Errorf("file not found locally: %s", localPath)
	}

	if exists {
		isNewer, err := CheckIfRemoteIsNewer(ctx, remoteURL, info.ModTime())
		config.GetLogger().Debug().Msgf("Is newer ?: %t", isNewer)
		if err != nil {
			config.GetLogger().Warn().Err(err).Msgf("Failed to check for freshness for %s, using cached version", localPath)
			return localPath, nil
		}
		if !isNewer {
			return localPath, nil
		}

		config.GetLogger().Info().Msgf("Remote file %s is newer, re-downloading...", remoteURL)
	}

	config.GetLogger().Info().Msgf("Downloading file from %s to %s", remoteURL, localPath)
	if err := common.DownloadFile(ctx, remoteURL, localPath); err != nil {
		return "", err
	}
	return localPath, nil
}
