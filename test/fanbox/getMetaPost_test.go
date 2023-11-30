package fanbox

import (
	"encoding/json"
	"kdownloader/pkg/fanbox"
	"os"
	"testing"
)

func TestGetMetaPost(t *testing.T) {
	postMeta := fanbox.GetMetaPost("https://richeonl0.fanbox.cc/posts/5831357")

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
