package utils

import (
	"io/ioutil"
	"kdownloader/utils"
	"testing"
)

func TestGetUrlBody(t *testing.T) {
	url := "https://www.fanbox.cc/@horosuke/posts/7013399"
	body := utils.GetUrlBody(url)

	// 将内容写入文件
	err := ioutil.WriteFile("output.html", body, 0644)
	if err != nil {
		panic(err)
	}
}
