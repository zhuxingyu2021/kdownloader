package db

import (
	kemono2 "kdownloader/pkg/kemono"
	"time"
)

type DBPostMeta struct {
	FileInDataBase bool

	PostsInfoID string

	Url string

	PosterInfo   kemono2.Poster
	PostInfoMeta kemono2.PostInfo

	PostContent string

	PostFiles     []string
	PostDownloads []string

	Tags []string
}

type PostRefType struct {
	PostId string
	Title  string
	Url    string
}

type DBPosterMeta struct {
	ID string

	PosterInfo kemono2.Poster

	FetchTime time.Time

	PostRef []PostRefType
}

type DBLinkQueryResult struct {
	DBQueryID string

	PostFiles     []string
	PostDownloads []string
}

func DBTypeConvert(allMeta *kemono2.PosterAll) (*DBPosterMeta, []*DBPostMeta) {
	retPosterMeta := &DBPosterMeta{
		ID:         allMeta.PosterAllMeta.ID,
		PosterInfo: allMeta.PosterAllMeta.PosterInfo,
		FetchTime:  allMeta.PosterAllMeta.FetchTime,
		PostRef:    make([]PostRefType, len(allMeta.PosterAllMeta.PostRef)),
	}

	for k, v := range allMeta.PosterAllMeta.PostRef {
		retPosterMeta.PostRef[k] = PostRefType{
			PostId: v.PostId,
			Title:  v.Title,
			Url:    v.Url,
		}
	}

	retDBPostMeta := make([]*DBPostMeta, len(allMeta.PosterAllDataLink))

	for k, v := range allMeta.PosterAllDataLink {
		retDBPostMeta[k] = &DBPostMeta{
			FileInDataBase: false,
			PostsInfoID:    v.PostsInfoID,
			Url:            v.Url,
			PosterInfo:     v.PosterInfo,
			PostInfoMeta:   v.PostInfoMeta,
			PostContent:    v.PostContent,
			PostFiles:      v.PostFiles,
			PostDownloads:  v.PostDownloads,
			Tags:           v.Tags,
		}
	}

	return retPosterMeta, retDBPostMeta
}
