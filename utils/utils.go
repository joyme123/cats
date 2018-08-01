package utils

import (
	"net/url"
	"os"
	"strings"
)

func GetAbsolutePath(rootDir string, uri string, indexFiles []string) string {
	var filepath string
	if strings.HasPrefix(uri, "http") {
		u, err := url.Parse(uri)
		if err != nil {
			return "400"
		}

		filepath = u.Path

	} else {
		filepath = uri
	}

	filepath = rootDir + filepath

	// 文件夹结尾,自动加上index文件
	if strings.HasSuffix(filepath, "/") {

		if len(indexFiles) > 0 {
			for _, v := range indexFiles {
				_, err := os.Stat(filepath + v)
				if err == nil {
					// 文件存在
					filepath = filepath + v
					break
				}
			}
		} else {
			// 默认为index.html
			filepath = filepath + "index.html"
		}

	}

	return filepath
}
