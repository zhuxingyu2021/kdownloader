package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

func GetUrlBody(url string) []byte {
	// 发送GET请求
	response, err := GetHttpCLimit(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return body
}

// 定义信号量大小
const concurrentLimit = 5

// 创建一个带缓冲的 channel 作为信号量
var semaphore = make(chan struct{}, concurrentLimit)

// GetHttpCLimit 封装了 http.Get 调用，使用信号量来限制并发数量
func GetHttpCLimit(url string) (response *http.Response, err error) {
	semaphore <- struct{}{}        // 获取信号量的一个插槽，如果信号量满了就会阻塞
	defer func() { <-semaphore }() // 函数返回前释放插槽

	fmt.Printf("Get %s\n", url)
	// 执行 http.Get 调用
	response, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	return response, nil
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
