package fanbox

import (
	"fmt"
	"go.uber.org/zap"
	"kdownloader/pkg/utils"
	"sort"
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

func GetMetaPost(url string) *PostMeta {
	ret := new(PostMeta)
	ret.Url = url

	urlSplit := strings.Split(url, "/")
	postId := urlSplit[len(urlSplit)-1]

	ret.PostInfoMeta.PostId = postId

	respJson, err := getRequestFanbox(fmt.Sprintf("https://api.fanbox.cc/post.info?postId=%s", postId))
	if err != nil {
		panic(err)
	}

	resultBody := respJson["body"].(map[string]interface{})

	isRestricted := resultBody["isRestricted"].(bool)
	if isRestricted {
		utils.Logger.Error("DownloadRestricted",
			zap.String("url", url))
		return nil
	}

	ret.PostInfoMeta.PostTitle = resultBody["title"].(string)
	ret.PostInfoMeta.PostPublished, err = time.Parse(time.RFC3339, resultBody["publishedDatetime"].(string))
	if err != nil {
		panic(err)
	}

	ret.PosterInfo.Platform = PLATFORM
	userRaw := resultBody["user"].(map[string]interface{})
	ret.PosterInfo.Username = userRaw["name"].(string)
	ret.PosterInfo.Userid, err = strconv.ParseInt(userRaw["userId"].(string), 10, 64)
	if err != nil {
		panic(err)
	}

	body := resultBody["body"].(map[string]interface{})

	var useText bool
	ret.PostContent, useText = body["text"].(string)

	citeOrder := map[string]int{}
	if !useText {
		ret.PostContent = ""
		blocksRaw := body["blocks"].([]interface{})

		for _, blockRaw := range blocksRaw {
			block := blockRaw.(map[string]interface{})
			if block["type"].(string) == "p" {
				ret.PostContent += block["text"].(string) + "\n"
			} else if block["type"].(string) == "image" {
				imageId := block["imageId"].(string)
				ret.PostContent += fmt.Sprintf("\\cite{%s}", imageId)

				_, exists := citeOrder[imageId]
				if !exists {
					citeOrder[imageId] = len(citeOrder)
				}
			}
		}
	}

	imagesRaw, useImages := body["images"].([]interface{})
	if useImages {
		for _, v := range imagesRaw {
			v_ := v.(map[string]interface{})
			ret.PostFiles = append(ret.PostFiles, "f"+v_["originalUrl"].(string))
		}
	} else {
		imageMap := body["imageMap"].(map[string]interface{})

		type _t struct {
			key string
			url string
		}
		var _ts []_t

		for k, v := range imageMap {
			v_ := v.(map[string]interface{})
			_ts = append(_ts, _t{
				key: k,
				url: "f" + v_["originalUrl"].(string),
			})
		}

		sort.Slice(_ts, func(i, j int) bool {
			iOrder, iFound := citeOrder[_ts[i].key]
			jOrder, jFound := citeOrder[_ts[j].key]

			if iFound && jFound {
				return iOrder < jOrder // 如果两个元素都在order中，按order中的顺序排列
			}
			if iFound {
				return false // 只有i在order中，i排在j后面
			}
			if jFound {
				return true // 只有j在order中，i排在j前面
			}
			return i < j // 如果两个元素都不在order中，保持原始顺序
		})

		for _, v := range _ts {
			ret.PostFiles = append(ret.PostFiles, v.url)
		}
	}

	tagsRaw := resultBody["tags"].([]interface{})
	for _, strRaw := range tagsRaw {
		ret.Tags = append(ret.Tags, strRaw.(string))
	}

	return ret
}
