package pixiv

import (
	"encoding/json"
	"kdownloader/pkg/pixiv"
	"os"
	"testing"
)

func TestGetMetaPost(t *testing.T) {
	postMeta := pixiv.GetMetaPost(`113710002`)

	// 序列化结构体为 JSON
	jsonData, err := json.MarshalIndent(postMeta, "", "  ")
	if err != nil {
		panic(err)
	}

	// 创建输出文件
	file, err := os.Create("getMetaPost.json")
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
