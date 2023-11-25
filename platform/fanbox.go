package platform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetTagFanbox(postId string) ([]string, error) {
	url := fmt.Sprintf("https://api.fanbox.cc/post.info?postId=%s", postId)

	// 发送GET请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头部
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("Origin", "https://www.fanbox.cc")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respJson map[string]map[string]interface{}
	err = json.Unmarshal(body, &respJson)

	if err != nil {
		return nil, err
	}

	retRaw, exists := respJson["body"]["tags"]
	if !exists {
		return nil, fmt.Errorf("Error when get fanbox tags")
	}

	retRawA, isArray := retRaw.([]interface{})
	if isArray {
		ret := []string{}
		for _, strRaw := range retRawA {
			str, isString := strRaw.(string)
			if !isString {
				return nil, fmt.Errorf("Error when get fanbox tags")
			}

			ret = append(ret, str)
		}

		return ret, nil
	} else {
		return nil, fmt.Errorf("Error when get fanbox tags")
	}
}
