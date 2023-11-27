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
	"strings"
	"sync"
)

type DFileInfo struct {
	path       string
	downloadOK bool
}

func DownloadNormalFile(ctx context.Context, url string, path string) error {
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

func DownloadPixivFile(ctx context.Context, url string, path string) error {
	utils2.Logger.Info("FileDownloading",
		zap.String("url", url),
		zap.String("path", path),
		zap.String("type", "pixiv"))

	globalDFileStatus := ctx.Value("FileStatus").(*sync.Map)
	globalDFileStatus.Store(url, DFileInfo{
		path:       path,
		downloadOK: false,
	})

	// Get the data
	header := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "zh-CN,zh;q=0.9",
		"Cache-Control":             "max-age=0",
		"Dnt":                       "1",
		"If-Modified-Since":         "Mon, 09 Sep 2019 23:00:01 GMT",
		"Referer":                   "https://www.pixiv.net/artworks/76712185",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
	}
	resp, err := utils2.GetHttpWithHeaderCLimit(url, header)
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

func DownloadFile(ctx context.Context, url string, path string) error {
	if strings.HasPrefix(url, "phttp") {
		return DownloadPixivFile(ctx, url[1:], path)
	} else {
		return DownloadNormalFile(ctx, url, path)
	}
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
