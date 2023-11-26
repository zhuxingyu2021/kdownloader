package kemono

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"kdownloader/pkg/utils"
	"sort"
	"sync"
)

type PosterAll struct {
	PosterAllMeta *PostsInfo

	posterAllDataLinkMu sync.Mutex
	PosterAllDataLink   []*PostMeta
}

func GetPosterAll(platform string, userid int64) *PosterAll {
	ret := new(PosterAll)
	ret.PosterAllMeta = SearchPoster(platform, userid)

	newUUID, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	UUID := newUUID.String()

	ret.PosterAllMeta.ID = UUID

	hasPanic := false
	var wg sync.WaitGroup
	for _, postRef := range ret.PosterAllMeta.PostRef {
		wg.Add(1)
		go func(url string) {
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
			}(GetMetaPost(url)))
			ret.posterAllDataLinkMu.Unlock()
		}(postRef.Url)
	}

	wg.Wait()

	if hasPanic {
		panic("error happens in go routine")
	}

	sort.Slice(ret.PosterAllDataLink, func(i, j int) bool {
		return ret.PosterAllDataLink[i].PostInfoMeta.PostId < ret.PosterAllDataLink[j].PostInfoMeta.PostId
	})

	return ret
}
