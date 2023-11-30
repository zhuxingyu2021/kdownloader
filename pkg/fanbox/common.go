package fanbox

import (
	"encoding/json"
	"errors"
	"kdownloader/pkg/utils"
	"net/http"
	"os"
	"time"
)

const PLATFORM string = "ffanbox"

func getRequestFanbox(url string) (map[string]interface{}, error) {
	// 发送GET请求
	cookie, _ := os.LookupEnv("FANBOX_COOKIE")
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Origin":     "https://www.fanbox.cc",
		"Accept":     "application/json, text/plain, */*",
		"Cookie":     cookie,
	}

	var response *utils.ResponseClimit
	var err error
	var retryCount = 1
	for {
		response, err = utils.GetHttpWithHeaderCLimit(url,
			headers)
		if err != nil {
			return nil, err
		}

		// Check server response
		if response.Resp.StatusCode == http.StatusOK {
			break
		} else if response.Resp.StatusCode == http.StatusTooManyRequests {
			// Retry
			time.Sleep(time.Second * time.Duration(retryCount*2))
			retryCount++
			response.Close()
			continue
		} else {
			response.Close()
			return nil, errors.New("bad status: " + response.Resp.Status)
		}
	}

	defer response.Close()
	// 解码 JSON 响应体到 map
	var result map[string]interface{}
	err = json.NewDecoder(response.Resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
