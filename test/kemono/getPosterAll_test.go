package kemono

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"kdownloader/pkg/kemono"
	"os"
	"testing"
)

func TestGetPosterAll(t *testing.T) {
	posterAll := kemono.GetPosterAll("fanbox", 3316400)

	assert.Equal(t, int64(len(posterAll.PosterAllMeta.PostRef)), int64(len(posterAll.PosterAllDataLink)))

	// 序列化结构体为 JSON
	jsonData, err := json.MarshalIndent(posterAll, "", "  ")
	if err != nil {
		panic(err)
	}

	// 创建输出文件
	file, err := os.Create("/tmp/3316400.json")
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
