package downloader

import (
	"os"
)

const DownloadFilePath string = "/tmp/kdl/download/"
const DownloadCompletePath string = "/tmp/kdl/complete/"
const DownloadZipPath string = "/tmp/kdl/zip/"

const DownloadRetryCount int = 5

func init() {
	err := os.MkdirAll(DownloadFilePath, 0755)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(DownloadCompletePath, 0755)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(DownloadZipPath, 0755)
	if err != nil {
		panic(err)
	}
}
