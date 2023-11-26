package main

import (
	"kdownloader/api"
	"kdownloader/pkg/utils"
)

func main() {
	defer utils.LoggerSync()
	api.InitApi()
}
