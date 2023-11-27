package utils

import (
	"archive/zip"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func GetUrlBody(url string) []byte {
	// 发送GET请求
	response, err := GetHttpCLimit(url)
	if err != nil {
		panic(err)
	}
	defer response.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

// 定义信号量大小
const concurrentLimit = 4

// 创建一个带缓冲的 channel 作为信号量
var semaphore = make(chan struct{}, concurrentLimit)

type ResponseClimit struct {
	Resp *http.Response
}

func (r *ResponseClimit) Close() error {
	err := r.Resp.Body.Close()

	<-semaphore

	return err
}

// GetHttpCLimit 封装了 http.Get 调用，使用信号量来限制并发数量
func GetHttpCLimit(url string) (response *ResponseClimit, err error) {
	semaphore <- struct{}{} // 获取信号量的一个插槽，如果信号量满了就会阻塞

	Logger.Info("httpRequest",
		zap.String("method", "GET"),
		zap.String("url", url))

	// 执行 http.Get 调用
	resp, err := http.Get(url)
	if err != nil {
		<-semaphore
		return nil, err
	}
	return &ResponseClimit{
		Resp: resp,
	}, nil
}

func GetHttpWithHeaderCLimit(url string, headers map[string]string) (response *ResponseClimit, err error) {
	semaphore <- struct{}{} // 获取信号量的一个插槽，如果信号量满了就会阻塞
	client := &http.Client{}

	Logger.Info("httpRequest",
		zap.String("method", "GET"),
		zap.String("url", url))

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		<-semaphore
		return nil, err
	}

	// 添加头部
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		<-semaphore
		return nil, err
	}

	return &ResponseClimit{
		Resp: resp,
	}, nil
}

var timeRegex = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
var timeRegex2 = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

func ExtractTime(text string) (time.Time, error) {
	subStr1 := timeRegex.FindString(text)
	if subStr1 != "" {
		const layout1 = "2006-01-02 15:04:05"
		return time.Parse(layout1, subStr1)
	} else {
		subStr2 := timeRegex2.FindString(text)
		if subStr2 != "" {
			const layout2 = "2006-01-02"
			return time.Parse(layout2, subStr2)
		} else {
			return time.Unix(0, 0), fmt.Errorf("Extract time failed for text: %s", text)
		}
	}
}

// ZipDirectoryToOSS 打包并上传目录到 OSS
func ZipDirectoryToOSS(dir, bucketName, objectName, endpoint, accessKeyID, accessKeySecret string) error {
	// 创建 OSS 客户端。
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return err
	}

	// 获取 OSS 存储桶。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 创建一个管道，用于传输 ZIP 文件。
	pr, pw := io.Pipe()
	// 在一个新的 goroutine 中进行 ZIP 打包。
	go func() {
		// 创建 zip writer。
		zw := zip.NewWriter(pw)
		// 遍历目录，添加文件到 ZIP。
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil // 忽略目录本身。
			}
			// 创建 ZIP 文件头。
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name, err = filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			header.Method = zip.Deflate // 设置压缩算法。
			w, err := zw.CreateHeader(header)
			if err != nil {
				return err
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(w, f)
			return err
		})
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		// 关闭 zip writer 来完成压缩过程。
		if err := zw.Close(); err != nil {
			pw.CloseWithError(err)
			return
		}
		pw.Close()
	}()

	// 使用管道读取器作为数据源进行 OSS 上传。
	err = bucket.PutObject(objectName, pr)
	if err != nil {
		pr.CloseWithError(err)
		return err
	}

	return nil
}
