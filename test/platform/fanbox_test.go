package platform

import (
	"kdownloader/platform"
	"testing"
)

func TestGetTag(t *testing.T) {
	tags, err := platform.GetTagFanbox("6786130")

	if err != nil {
		panic(err)
	}

	print(tags)
}
