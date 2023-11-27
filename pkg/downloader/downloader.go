package downloader

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	utils2 "kdownloader/pkg/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

type DFileInfo struct {
	path       string
	downloadOK bool
}

func DownloadFile(ctx context.Context, url string, path string) error {
	utils2.Logger.Info("FileDownloading",
		zap.String("url", url),
		zap.String("path", path))

	globalDFileStatus := ctx.Value("FileStatus").(*sync.Map)
	globalDFileStatus.Store(url, DFileInfo{
		path:       path,
		downloadOK: false,
	})

	// Get the data
	resp, err := utils2.GetHttpCLimit(url)
	if err != nil {
		return err
	}
	defer resp.Close()

	// Check server response
	if resp.Resp.StatusCode != http.StatusOK {
		return errors.New("bad status: " + resp.Resp.Status)
	}

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Resp.Body)
	return err
}

func getUrlExt(rawUrl string) (string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	// 获取路径部分的扩展名
	ext := filepath.Ext(parsedURL.Path)

	return ext, nil
}

func DownloadContext(parent context.Context) context.Context {
	return context.WithValue(parent, "FileStatus", new(sync.Map))
}

// DWorker listens for URLs on the channel and downloads them.
// It also listens for a context cancellation to stop the worker.
func DWorker(ctx context.Context, urls <-chan string) {
	var fileID int64 = 0
	var wg sync.WaitGroup

	globalDFileStatus := ctx.Value("FileStatus").(*sync.Map)
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			utils2.Logger.Info("DWorker",
				zap.String("action", "Done"))
			return
		case url, ok := <-urls:
			if !ok {
				wg.Wait()
				utils2.Logger.Info("DWorker",
					zap.String("action", "Done"))
				return
			}

			ext, err := getUrlExt(url)
			if err != nil {
				utils2.Logger.Error("DWorker",
					zap.String("action", "ErrorUrl"),
					zap.String("url", url),
					zap.Error(err))
			} else {
				path := fmt.Sprintf("%s%016x%s", DownloadFilePath, fileID, ext)
				fileID++

				wg.Add(1)
				go func(url string, path string) {
					defer wg.Done()
					for retry := 0; retry < DownloadRetryCount; retry++ {
						err := DownloadFile(ctx, url, path)
						if err != nil {
							utils2.Logger.Error("DWorker",
								zap.String("action", "ErrorDownload"),
								zap.String("url", url),
								zap.String("path", path),
								zap.Error(err))
						} else {
							globalDFileStatus.Store(url, DFileInfo{
								path:       path,
								downloadOK: true,
							})
							utils2.Logger.Info("FileDownloadOK",
								zap.String("url", url),
								zap.String("path", path))
							return
						}
					}
					globalDFileStatus.Delete(url)
					utils2.Logger.Fatal("DWorker",
						zap.String("action", "FatalDownload"),
						zap.String("url", url),
						zap.String("path", path),
						zap.Error(err))
				}(url, path)
			}
		}
	}
}

func GetUndownloadUrls(ctx context.Context) []string {
	var ret []string

	globalDFileStatus := ctx.Value("FileStatus").(*sync.Map)
	globalDFileStatus.Range(func(key, value interface{}) bool {
		dFileInfo := value.(DFileInfo)
		if !dFileInfo.downloadOK {
			ret = append(ret, key.(string))
		}
		return true
	})

	return ret
}

func ListOKUrls(ctx context.Context) map[string]string {
	ret := map[string]string{}

	globalDFileStatus := ctx.Value("FileStatus").(*sync.Map)
	globalDFileStatus.Range(func(key, value interface{}) bool {
		dFileInfo := value.(DFileInfo)
		if dFileInfo.downloadOK {
			ret[key.(string)] = dFileInfo.path
		}
		return true
	})

	return ret
}

func DeleteUrls(ctx context.Context, urls []string) {
	globalDFileStatus := ctx.Value("FileStatus").(*sync.Map)
	for _, u := range urls {
		globalDFileStatus.Delete(u)
	}
}
