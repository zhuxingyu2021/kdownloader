package pixiv

import (
	"kdownloader/pkg/pixiv"
	"testing"
)

func TestSearchPoster(t *testing.T) {
	postsInfo := pixiv.SearchPoster(86929043)

	println("count: ", len(postsInfo))
	for _, v := range postsInfo {
		println(v)
	}
}
