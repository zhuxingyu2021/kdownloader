package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"kdownloader/pkg/downloader"
	"kdownloader/pkg/utils"
	"net/http"
	"os"
	"sync"
)

var runPipe bool = false
var runPipeMu sync.RWMutex

func getPipe(c *gin.Context) {
	runPipeMu.RLock()
	c.JSON(http.StatusOK, map[string]bool{
		"RUNNING": runPipe,
	})
	runPipeMu.RUnlock()
}

func startPipe(c *gin.Context) {
	runPipeMu.Lock()
	if runPipe {
		runPipeMu.Unlock()
		c.JSON(http.StatusServiceUnavailable, map[string]bool{
			"RUNNING": true,
		})
		return
	}
	runPipe = true
	runPipeMu.Unlock()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
		return
	}

	var conf downloader.GlobalConfig
	_ = json.Unmarshal(body, &conf)

	URI, exists := os.LookupEnv("MONGO_URI")
	if exists {
		conf.URI = URI
	}

	AccessKeyID, exists := os.LookupEnv("OSS_ACCESS_KEY_ID")
	if exists {
		conf.OSS.AccessKeyID = AccessKeyID
	}

	AccessKeySecret, exists := os.LookupEnv("OSS_ACCESS_KEY_SECRET")
	if exists {
		conf.OSS.AccessKeySecret = AccessKeySecret
	}

	c.String(http.StatusOK, "{}")

	go func(c downloader.GlobalConfig) {
		utils.Logger.Info("PipeTask",
			zap.String("action", "Start"),
			zap.Any("config", c))
		err := downloader.PipeTask(&c)

		if err != nil {
			utils.Logger.Error("PipeTask",
				zap.String("action", "Error"),
				zap.Error(err))
		}

		utils.Logger.Info("PipeTask",
			zap.String("action", "Done"))

		runPipeMu.Lock()
		runPipe = false
		runPipeMu.Unlock()
	}(conf)
}
