package platform

import "strings"

func GetTag(platform string, postId string) ([]string, error) {
	switch strings.ToLower(platform) {
	case "fanbox":
		return GetTagFanbox(postId)
	}
	return nil, nil
}
