package downloader

import (
	"os"
	"time"
)

const DownloadFilePath string = "/tmp/kdl/download/"
const DownloadCompletePath string = "/tmp/kdl/complete/"
const DownloadZipPath string = "/tmp/kdl/zip/"

const DownloadRetryCount int = 8

func InitDownload() error {
	os.RemoveAll(DownloadFilePath)
	os.RemoveAll(DownloadCompletePath)
	os.RemoveAll(DownloadZipPath)

	time.Sleep(time.Second)

	err := os.MkdirAll(DownloadFilePath, 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(DownloadCompletePath, 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(DownloadZipPath, 0755)
	if err != nil {
		return err
	}

	return nil
}
