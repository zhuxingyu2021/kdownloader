package fanbox

import (
	"encoding/json"
	"kdownloader/pkg/fanbox"
	"os"
	"testing"
)

func TestSearchPoster(t *testing.T) {
	postsInfo := fanbox.SearchPoster("richeonl0")

	// 序列化结构体为 JSON
	jsonData, err := json.MarshalIndent(postsInfo, "", "  ")
	if err != nil {
		panic(err)
	}

	// 创建输出文件
	file, err := os.Create("searchPoster.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 写入数据到文件
	_, err = file.Write(jsonData)
	if err != nil {
		panic(err)
	}
}
