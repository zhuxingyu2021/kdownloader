package downloader

import (
	"kdownloader/pkg/downloader"
	"os"
	"testing"
)

func TestPipeTask(t *testing.T) {
	URI, _ := os.LookupEnv("MONGO_URI")
	AccessKeyID, _ := os.LookupEnv("OSS_ACCESS_KEY_ID")
	AccessKeySecret, _ := os.LookupEnv("OSS_ACCESS_KEY_SECRET")
	config := downloader.GlobalConfig{
		URI:    URI,
		DBName: `kdb`,
		OSS: downloader.OSSConfig{
			BucketName:      `hk-test-zxy`,
			EndPoint:        `oss-cn-hongkong.aliyuncs.com`,
			AccessKeyID:     AccessKeyID,
			AccessKeySecret: AccessKeySecret,
		},
	}

	err := downloader.PipeTask(&config)
	if err != nil {
		panic(err)
	}
}
