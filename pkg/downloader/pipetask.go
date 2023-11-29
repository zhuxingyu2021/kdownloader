package downloader

import (
	"context"
	"go.uber.org/zap"
	db2 "kdownloader/pkg/db"
	"kdownloader/pkg/utils"
	"os"
	"path/filepath"
	"time"
)

type OSSConfig struct {
	BucketName string
	EndPoint   string

	AccessKeyID     string
	AccessKeySecret string
}

type GlobalConfig struct {
	URI    string
	DBName string

	OSS OSSConfig
}

func removeIndex(s []*db2.DBLinkQueryResult, idx []int) []*db2.DBLinkQueryResult {
	result := make([]*db2.DBLinkQueryResult, 0, len(s))
	deleteIdx := 0
	for i, v := range s {
		if deleteIdx < len(idx) && i == idx[deleteIdx] {
			// Skip this index as it's marked for deletion
			deleteIdx++
		} else {
			result = append(result, v)
		}
	}
	return result
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func copyFilesToComplete(path string) (string, error) {
	filename := filepath.Base(path)

	newPath := DownloadCompletePath + filename

	utils.Logger.Info("Rename",
		zap.String("oldpath", path),
		zap.String("newpath", newPath),
		zap.Bool("exists", exists(path)))

	return newPath, os.Rename(path, newPath)
}

func checkCompleteDownloading(ctx context.Context, qResult []*db2.DBLinkQueryResult, zchan chan ZWorkerArg) error {
	okUrls := ListOKUrls(ctx)
	var okPostID []int
	var okPostUrls []string

	// 计算已经下载完成的postID
	for i, v := range qResult {
		var notOK = false
		for _, url := range v.PostFiles {
			_, exists := okUrls[url]
			if !exists {
				notOK = true
				break
			}
		}
		if notOK {
			continue
		}
		for _, url := range v.PostDownloads {
			_, exists := okUrls[url]
			if !exists {
				notOK = true
				break
			}
		}
		if notOK {
			continue
		}

		okPostID = append(okPostID, i)
	}

	if len(okPostID) > 0 {
		// 统计下载完成的url
		for _, postID := range okPostID {
			v := qResult[postID]

			for _, url := range v.PostFiles {
				okPostUrls = append(okPostUrls, url)
			}

			for _, url := range v.PostDownloads {
				okPostUrls = append(okPostUrls, url)
			}
		}

		DeleteUrls(ctx, okPostUrls)

		// 移动所有下载完成的文件到完成目录
		for _, postID := range okPostID {
			v := qResult[postID]
			files := map[string]bool{}

			for _, url := range v.PostFiles {
				files[okUrls[url]] = true
			}

			for _, url := range v.PostDownloads {
				files[okUrls[url]] = true
			}

			var zArg ZWorkerArg
			for file := range files {
				newPath, err := copyFilesToComplete(file)
				if err != nil {
					return err
				}

				zArg.files = append(zArg.files, newPath)
			}
			zArg.postID = v.DBQueryID
			zchan <- zArg
		}

		qResult = removeIndex(qResult, okPostID)
	}

	return nil
}

func PipeTask(config *GlobalConfig) error {
	err := InitDownload()
	if err != nil {
		return err
	}

	ctx_, cancel := context.WithCancel(context.Background())
	ctx := DownloadContext(ctx_)
	defer cancel()

	cli, err := db2.InitMongo(ctx, config.URI, config.DBName)
	if err != nil {
		return err
	}
	defer cli.Close()

	urlchan := make(chan string)
	zchan := make(chan ZWorkerArg)
	go DWorker(ctx, urlchan)
	go ZWorker(ctx, zchan, config, cli)

	qResult, err := cli.LinkQuery()
	if err != nil {
		return err
	}

	for _, v := range qResult {
		for _, url := range v.PostFiles {
			urlchan <- url
		}
		for _, url := range v.PostDownloads {
			urlchan <- url
		}
	}

	for {
		time.Sleep(time.Second * 5)

		err = checkCompleteDownloading(ctx, qResult, zchan)
		if err != nil {
			return err
		}

		// 所有Url都下载完成
		if len(ListNOKUrls(ctx)) == 0 {
			err = checkCompleteDownloading(ctx, qResult, zchan)
			if err != nil {
				return err
			}

			break
		}
	}

	zchan <- ZWorkerArg{
		postID: "zzzclosehaha",
		files:  nil,
	}

	return nil
}
