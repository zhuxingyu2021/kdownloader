package downloader

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	utils2 "kdownloader/pkg/utils"
	"os"
	"path/filepath"
)

type ZWorkerArg struct {
	postID string
	files  []string
}

func zipFiles(config OSSConfig, files []string, postsID string) error {
	zipFileName := postsID
	zipDirectoryPath := DownloadZipPath + zipFileName + "/"

	err := os.MkdirAll(zipDirectoryPath, 0755)
	if err != nil {
		return err
	}

	for i, file := range files {
		oldPath := file
		ext := filepath.Ext(oldPath)
		newPath := fmt.Sprintf("%s%04x%s", zipDirectoryPath, i, ext)

		err := os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
	}

	// 打包 ZIP 文件
	err = utils2.ZipDirectoryToOSS(zipDirectoryPath, config.BucketName, zipFileName+".zip",
		config.EndPoint, config.AccessKeyID, config.AccessKeySecret)

	return err
}

func ZWorker(ctx context.Context, zchan <-chan ZWorkerArg, config *GlobalConfig) {
	for {
		select {
		case <-ctx.Done():
			utils2.Logger.Info("ZWorker",
				zap.String("action", "Done"))
			return
		case work, ok := <-zchan:
			if !ok {
				utils2.Logger.Info("ZWorker",
					zap.String("action", "Done"))
				return
			}
			utils2.Logger.Info("ZWorker",
				zap.String("action", "Task"),
				zap.String("postID", work.postID))

			err := zipFiles(config.OSS, work.files, work.postID)
			if err != nil {
				utils2.Logger.Error("ZWorker",
					zap.String("action", "Error"),
					zap.Error(err))
			}
		}
	}
}
