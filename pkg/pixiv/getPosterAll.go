package pixiv

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"kdownloader/pkg/utils"
	"sort"
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

	PostRef []PostRefType
}

type PosterAll struct {
	PosterAllMeta *PostsInfo

	posterAllDataLinkMu sync.Mutex
	PosterAllDataLink   []*PostMeta
}

func GetPosterAll(userid int64) *PosterAll {
	ret := new(PosterAll)
	postIDs := SearchPoster(userid)

	newUUID, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	UUID := newUUID.String()

	ret.PosterAllMeta = new(PostsInfo)
	ret.PosterAllMeta.ID = UUID
	ret.PosterAllMeta.FetchTime = time.Now()

	hasPanic := false
	var wg sync.WaitGroup
	for _, postID := range postIDs {
		wg.Add(1)
		go func(postID string) {
			defer func() {
				if err := recover(); err != nil {
					utils.Logger.Error("GetPosterAll",
						zap.Any("error", err),
						zap.Stack("stack"),
					)
					hasPanic = true
				}
			}()
			defer wg.Done()

			ret.posterAllDataLinkMu.Lock()
			ret.PosterAllDataLink = append(ret.PosterAllDataLink, func(postMeta *PostMeta) *PostMeta {
				postMeta.PostsInfoID = UUID
				return postMeta
			}(GetMetaPost(postID)))
			ret.posterAllDataLinkMu.Unlock()
		}(postID)
	}

	wg.Wait()

	if hasPanic {
		panic("error happens in go routine")
	}

	sort.Slice(ret.PosterAllDataLink, func(i, j int) bool {
		return ret.PosterAllDataLink[i].PostInfoMeta.PostId < ret.PosterAllDataLink[j].PostInfoMeta.PostId
	})

	for i, v := range ret.PosterAllDataLink {
		if i == 0 {
			ret.PosterAllMeta.PosterInfo = v.PosterInfo
		}
		ret.PosterAllMeta.PostRef = append(ret.PosterAllMeta.PostRef, PostRefType{
			PostId:      v.PostInfoMeta.PostId,
			Title:       v.PostInfoMeta.PostTitle,
			Url:         v.Url,
			PageId:      0,
			InterPageId: i,
		})
	}

	return ret
}
