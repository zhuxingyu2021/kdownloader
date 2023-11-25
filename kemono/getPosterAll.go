package kemono

import (
	"fmt"
	"github.com/google/uuid"
	"runtime/debug"
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
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Panic caught: %v\n", err)
					debug.PrintStack()
					hasPanic = true
				}
			}()

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
