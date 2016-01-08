package utilmagic

import (
	"net/url"
)

func MapToString(map_ map[string]string) string {
	values := url.Values{}
	for k, v := range map_ {
		values.Set(k, v)
	}
	return values.Encode()
}

func MapMerge(maps ...map[string]string) map[string]string {
	mergedMap := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			mergedMap[k] = v
		}
	}
	return mergedMap
}
