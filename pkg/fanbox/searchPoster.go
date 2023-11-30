package fanbox

import (
	"fmt"
	"go.uber.org/zap"
	"kdownloader/pkg/utils"
	"strconv"
	"time"
)

type PostRefType struct {
	PostId      string
	Title       string
	Url         string
	PageId      int
	InterPageId int

	CoverUrl string
}

type PostsInfo struct {
	ID string

	PosterInfo Poster

	FetchTime time.Time

	PostRef []PostRefType
}

func getPageUrl(username string) ([]string, error) {
	var ret []string

	respJson, err := getRequestFanbox(fmt.Sprintf("https://api.fanbox.cc/post.paginateCreator?creatorId=%s", username))

	if err != nil {
		return nil, err
	}

	resultRaw := respJson["body"].([]interface{})
	for _, result := range resultRaw {
		ret = append(ret, result.(string))
	}

	return ret, nil
}

func SearchPoster(username string) *PostsInfo {
	var ret PostsInfo
	ret.FetchTime = time.Now()
	ret.PosterInfo.Platform = PLATFORM
	ret.PosterInfo.Username = username

	pageUrls, err := getPageUrl(username)

	if err != nil {
		panic(err)
	}

	for i, pageUrl := range pageUrls {
		respJson, err := getRequestFanbox(pageUrl)

		if err != nil {
			panic(err)
		}

		resultBody := respJson["body"].(map[string]interface{})
		resultItems := resultBody["items"].([]interface{})

		for j, itemRaw := range resultItems {
			item := itemRaw.(map[string]interface{})

			user := item["user"].(map[string]interface{})
			ret.PosterInfo.Userid, err = strconv.ParseInt(user["userId"].(string), 10, 64)
			if err != nil {
				panic(err)
			}

			postId := item["id"].(string)

			isRestricted := item["isRestricted"].(bool)
			if isRestricted {
				utils.Logger.Info("restricted",
					zap.String("userName", username),
					zap.String("postId", postId),
					zap.String("url", fmt.Sprintf("https://%s.fanbox.cc/posts/%s", username, postId)))
				continue
			}

			cover := item["cover"].(map[string]interface{})
			ret.PostRef = append(ret.PostRef, PostRefType{
				PostId:      postId,
				Title:       item["title"].(string),
				Url:         fmt.Sprintf("https://%s.fanbox.cc/posts/%s", username, postId),
				PageId:      i,
				InterPageId: j,

				CoverUrl: "f" + cover["url"].(string),
			})
		}
	}

	return &ret
}
