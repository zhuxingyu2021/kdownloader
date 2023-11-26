package kemono

import (
	"encoding/json"
	"kdownloader/pkg/kemono"
	"os"
	"testing"
)

func TestGetMetaPost(t *testing.T) {
	// url := "https://kemono.su/patreon/user/93254587/post/90237820"
	url := `https://kemono.su/patreon/user/93254587/post/92144542`
	postMeta := kemono.GetMetaPost(url)

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
