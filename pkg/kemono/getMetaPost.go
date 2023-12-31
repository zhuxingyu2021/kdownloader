package kemono

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"kdownloader/pkg/platform"
	"kdownloader/pkg/utils"
	"net/http"
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

func getPoster(doc *goquery.Document, poster *Poster) {
	doc.Find("a.post__user-name").Each(func(index int, item *goquery.Selection) {
		// 提取href属性
		href, exists := item.Attr("href")
		if exists {
			var err error

			poster.Username = strings.TrimSpace(item.Text())

			hrefSplit := strings.Split(href, "/")
			poster.Platform = hrefSplit[len(hrefSplit)-3]
			poster.Userid, err = strconv.ParseInt(hrefSplit[len(hrefSplit)-1], 10, 64)

			if err != nil {
				panic(err)
			}
		}
	})
}

func getPostInfo(doc *goquery.Document, postinfo *PostInfo) {
	doc.Find(".post__info").Each(func(index int, item *goquery.Selection) {
		item.Find(".post__title").Each(func(index int, item *goquery.Selection) {
			postinfo.PostTitle = strings.TrimSpace(item.Text())
		})
		item.Find(".post__published").Each(func(index int, item *goquery.Selection) {
			var err error

			postinfo.PostPublished, err = utils.ExtractTime(item.Text())
			if err != nil {
				panic(err)
			}
		})
	})
}

func getPostTags(doc *goquery.Document, tags *[]string) {
	doc.Find(".post__info").Each(func(index int, item *goquery.Selection) {
		item.Find("#post-tags").Each(func(index int, item *goquery.Selection) {
			item.Find("a").Each(func(index int, item *goquery.Selection) {
				*tags = append(*tags, strings.TrimSpace(item.Text()))
			})
		})
	})
}

func getPostContent(doc *goquery.Document, postContent *string) {
	doc.Find(".post__content").Each(func(index int, item *goquery.Selection) {
		*postContent = strings.TrimSpace(item.Text())
	})
}

func getPostFiles(doc *goquery.Document, postFiles *[]string) {
	doc.Find(".post__files").Each(func(index int, item *goquery.Selection) {
		item.Find(".post__thumbnail").Each(func(index int, item *goquery.Selection) {
			item.Find("a").Each(func(index int, item *goquery.Selection) {
				// 提取href属性
				href, exists := item.Attr("href")
				if exists {
					*postFiles = append(*postFiles, href)
				}
			})
		})
	})
}

func getPostDownloads(doc *goquery.Document, postDownloads *[]string) {
	doc.Find(".post__attachments").Each(func(index int, item *goquery.Selection) {
		item.Find("a").Each(func(index int, item *goquery.Selection) {
			// 提取href属性
			href, exists := item.Attr("href")
			if exists {
				*postDownloads = append(*postDownloads, href)
			}
		})
	})
}

func GetMetaPost(url string) *PostMeta {
	ret := new(PostMeta)
	ret.Url = url

	urlSplit := strings.Split(url, "/")
	postId := urlSplit[len(urlSplit)-1]
	ret.PostInfoMeta.PostId = postId

	// 发送GET请求
	response, err := utils.GetHttpCLimit(url)
	if err != nil {
		panic(err)
	}
	defer response.Close()

	// Check server response
	if response.Resp.StatusCode != http.StatusOK {
		panic(errors.New("bad status: " + response.Resp.Status))
	}

	doc, err := goquery.NewDocumentFromReader(response.Resp.Body)
	if err != nil {
		panic(err)
	}

	getPostInfo(doc, &ret.PostInfoMeta)
	getPoster(doc, &ret.PosterInfo)
	getPostContent(doc, &ret.PostContent)
	getPostFiles(doc, &ret.PostFiles)
	getPostDownloads(doc, &ret.PostDownloads)
	getPostTags(doc, &ret.Tags)

	platformTag, err := platform.GetTag(ret.PosterInfo.Platform, ret.PostInfoMeta.PostId)
	if err != nil {
		utils.Logger.Info("PlatformError",
			zap.String("Action", "GetTag"),
			zap.String("Platform", ret.PosterInfo.Platform),
			zap.String("PostId", ret.PostInfoMeta.PostId))
	}
	ret.Tags = append(ret.Tags, platformTag...)

	return ret
}
