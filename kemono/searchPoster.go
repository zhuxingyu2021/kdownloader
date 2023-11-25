package kemono

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"kdownloader/utils"
	"math"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PostRefType struct {
	PostId      string
	Title       string
	Url         string
	PageId      int
	InterPageId int
}

type PostsInfo struct {
	ID string

	PosterInfo Poster

	FetchTime time.Time

	postRefMu sync.Mutex
	PostRef   []PostRefType
}

func getPostCount(doc *goquery.Document) int64 {
	div := doc.Find("div.paginator#paginator-top")

	if div.Length() == 0 {
		panic("getPostCount error: paginator not found")
	}

	small := div.Find("small")
	if small.Length() == 0 {
		panic("getPostCount error: small not found")
	}

	countText := strings.Split(strings.TrimSpace(small.Text()), " ")
	count, err := strconv.ParseInt(countText[len(countText)-1], 10, 64)

	if err != nil {
		panic(err)
	}

	return count
}

func getPosterS(doc *goquery.Document, poster *Poster) {
	doc.Find("a.user-header__profile").Each(func(index int, item *goquery.Selection) {
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

func searchPosterPages(doc *goquery.Document, info *PostsInfo, page int) {
	iid := 0
	doc.Find(".card-list__items").Each(func(index int, item *goquery.Selection) {
		item.Find("a").Each(func(index int, item *goquery.Selection) {
			// 提取href属性
			href, exists := item.Attr("href")
			if exists {
				var postRef PostRefType

				hrefSplit := strings.Split(href, "/")
				postRef.PostId = hrefSplit[len(hrefSplit)-1]

				postRef.Url = href
				if !strings.HasPrefix(postRef.Url, "http") {
					postRef.Url = fmt.Sprintf("https://kemono.su%s", href)
				}

				item.Find("header.post-card__header").Each(func(index int, item *goquery.Selection) {
					postRef.Title = strings.TrimSpace(item.Text())
				})
				postRef.PageId = page
				postRef.InterPageId = iid
				iid++

				info.postRefMu.Lock()
				info.PostRef = append(info.PostRef, postRef)
				info.postRefMu.Unlock()
			}
		})
	})
}

func searchPosterInternal(url string) (*PostsInfo, int64) {
	ret := new(PostsInfo)
	ret.FetchTime = time.Now()

	// 发送GET请求
	response, err := utils.GetHttpCLimit(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}

	count := getPostCount(doc)
	pages := int(math.Ceil(float64(count) / 50.0))

	var wg sync.WaitGroup

	hasPanic := false
	for i := 1; i < pages; i++ {
		urlInternal := fmt.Sprintf("%s?o=%d", url, i*50)
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Panic caught: %v\n", err)
					debug.PrintStack()
					hasPanic = true
				}
			}()
			// 发送GET请求
			response, err := utils.GetHttpCLimit(urlInternal)
			if err != nil {
				panic(err)
			}
			defer response.Body.Close()

			doc, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				panic(err)
			}

			searchPosterPages(doc, ret, page)
		}(i)
	}
	getPosterS(doc, &ret.PosterInfo)
	searchPosterPages(doc, ret, 0)

	wg.Wait()

	if hasPanic {
		panic("error happens in go routine")
	}
	sort.Slice(ret.PostRef, func(i, j int) bool {
		return ret.PostRef[i].PostId < ret.PostRef[j].PostId
	})

	return ret, count
}

func SearchPoster(platform string, userid int64) (*PostsInfo, int64) {
	url := fmt.Sprintf("https://kemono.su/%s/user/%d", platform, userid)

	return searchPosterInternal(url)
}
