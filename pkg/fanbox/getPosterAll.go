package fanbox

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

func GetPosterAll(username string) *PosterAll {
	ret := new(PosterAll)
	ret.PosterAllMeta = SearchPoster(username)

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
		go func(url string, coverUrl string) {
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
				postMeta.PostFiles = append([]string{coverUrl}, postMeta.PostFiles...)
				return postMeta
			}(GetMetaPost(url)))
			ret.posterAllDataLinkMu.Unlock()
		}(postRef.Url, postRef.CoverUrl)
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
