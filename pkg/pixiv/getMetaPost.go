package pixiv

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"kdownloader/pkg/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Poster struct {
	Username string

	Platform string
	Userid   int64
}

type PostInfo struct {
	PostTitle string
	PostId    string

	PostPublished time.Time
}

type PostMeta struct {
	PostsInfoID string

	Url string

	PosterInfo   Poster
	PostInfoMeta PostInfo

	PostContent string

	PostFiles     []string
	PostDownloads []string

	Tags []string
}

func getRequestPixiv(url string) (map[string]interface{}, error) {
	// 发送GET请求
	cookie, _ := os.LookupEnv("PIXIV_COOKIE")
	headers := map[string]string{
		"cookie":     cookie,
		"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
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

func getPostMetadata(postID string, postMeta *PostMeta) (error, int) {
	postMeta.Url = "https://www.pixiv.net/artworks/" + postID

	postMeta.PostInfoMeta.PostId = postID

	reqUrl := "https://www.pixiv.net/ajax/illust/" + postID

	result, err := getRequestPixiv(reqUrl)
	if err != nil {
		return err, 0
	}

	resultBody := result["body"].(map[string]interface{})
	postMeta.PostInfoMeta.PostTitle = resultBody["title"].(string)
	postMeta.PostInfoMeta.PostPublished, err = time.Parse(time.RFC3339, resultBody["uploadDate"].(string))
	if err != nil {
		return err, 0
	}

	postMeta.PosterInfo.Platform = PLATFORM

	userIDRaw := resultBody["userId"].(string)
	postMeta.PosterInfo.Userid, err = strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		return err, 0
	}

	postMeta.PosterInfo.Username = resultBody["userName"].(string)

	resultTags1 := resultBody["tags"].(map[string]interface{})
	resultTags2 := resultTags1["tags"].([]interface{})

	for _, ritem := range resultTags2 {
		item := ritem.(map[string]interface{})
		postMeta.Tags = append(postMeta.Tags, item["tag"].(string))
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resultBody["description"].(string)))
	if err != nil {
		return err, 0
	}

	postMeta.PostContent = strings.TrimSpace(doc.Text())

	count := int(resultBody["pageCount"].(float64))
	return nil, count
}

func getPostFiles(postID string, postFiles *[]string) error {
	reqUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/pages?lang=zh", postID)

	result, err := getRequestPixiv(reqUrl)
	if err != nil {
		return err
	}

	resultBody := result["body"].([]interface{})

	urls := []string{}
	for _, itemR := range resultBody {
		item := itemR.(map[string]interface{})
		urlsRaw := item["urls"].(map[string]interface{})
		ori := urlsRaw["original"].(string)

		urls = append(urls, "p"+ori)
	}

	*postFiles = append(*postFiles, urls...)

	return nil
}

func GetMetaPost(postID string) *PostMeta {
	ret := new(PostMeta)

	err, count := getPostMetadata(postID, ret)

	if err != nil {
		panic(err)
	}

	err = getPostFiles(postID, &ret.PostFiles)

	if err != nil {
		panic(err)
	}

	if count != len(ret.PostFiles) {
		panic("file count not equal")
	}

	return ret
}
