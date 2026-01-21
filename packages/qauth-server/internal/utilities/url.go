package utilities

import (
	"net/url"
)

// PickHashParams 从 URL Hash 中筛选指定的参数并返回新的 URL
func PickHashParams(urlStr string, pick []string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// 解析 Hash 参数
	fragmentValues, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return "", err
	}

	// 筛选参数
	newValues := url.Values{}
	for _, key := range pick {
		if vals, ok := fragmentValues[key]; ok {
			newValues[key] = vals
		}
	}

	// 重新构造 Fragment 并返回 URL
	u.Fragment = newValues.Encode()
	return u.String(), nil
}
