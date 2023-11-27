package pixiv

import (
	"encoding/json"
	"errors"
	"fmt"
	"kdownloader/pkg/utils"
	"net/http"
	"os"
)

func SearchPoster(userid int64) []string {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/user/%d/profile/all?lang=zh", userid)

	cookie, _ := os.LookupEnv("PIXIV_COOKIE")
	headers := map[string]string{
		"cookie":     cookie,
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
	}
	response, err := utils.GetHttpWithHeaderCLimit(url, headers)
	if err != nil {
		panic(err)
	}
	defer response.Close()

	// Check server response
	if response.Resp.StatusCode != http.StatusOK {
		panic(errors.New("bad status: " + response.Resp.Status))
	}

	// 解码 JSON 响应体到 map
	var result map[string]interface{}
	err = json.NewDecoder(response.Resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	resultBody := result["body"].(map[string]interface{})

	illusts := resultBody["illusts"].(map[string]interface{})

	var postsList []string
	for k, _ := range illusts {
		postsList = append(postsList, k)
	}

	return postsList
}
