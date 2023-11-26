package utils

import (
	"io/ioutil"
	"kdownloader/pkg/utils"
	"os"
	"testing"
)

func TestGetUrlBody(t *testing.T) {
	url := "https://kemono.su/patreon/user/93254587/post/92144542"
	body := utils.GetUrlBody(url)

	// 将内容写入文件
	err := ioutil.WriteFile("output.html", body, 0644)
	if err != nil {
		panic(err)
	}
}

func TestZipDirectoryToOSS(t *testing.T) {
	path := `/Users/zhuxingyu/Desktop/mypaper`
	bucketName := `hk-test-zxy`
	accessKeyId, exists := os.LookupEnv("OSS_ACCESS_KEY_ID")
	if !exists {
		panic("AccessKey not exists")
	}
	accessKeySecret, exists := os.LookupEnv("OSS_ACCESS_KEY_SECRET")
	if !exists {
		panic("AccessKeySecret not exists")
	}
	err := utils.ZipDirectoryToOSS(path, bucketName, `mypaper`, `oss-cn-hongkong.aliyuncs.com`,
		accessKeyId, accessKeySecret)

	if err != nil {
		panic(err)
	}
}
