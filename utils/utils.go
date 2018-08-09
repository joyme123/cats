package utils

import (
	"net/url"
	"os"
	"strings"
)

func CompleteURI(uri string, indexFiles []string) string {

	// 文件夹结尾,自动加上index文件
	if strings.HasSuffix(uri, "/") {

		if len(indexFiles) > 0 {
			for _, v := range indexFiles {
				_, err := os.Stat(uri + v)
				if err == nil {
					// 文件存在
					uri = uri + v
					break
				}
			}
		} else {
			// 默认为index.html
			uri = uri + "index.html"
		}

	}

	return uri
}

func GetAbsolutePath(rootDir string, uri string, indexFiles []string) (string, string) {
	var filepath string
	completeURI := uri

	if strings.HasPrefix(uri, "http") {
		u, err := url.Parse(uri)
		if err != nil {
			return "", "400"
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
					completeURI += v
					break
				}
			}
		} else {
			// 默认为index.html
			filepath = filepath + "index.html"
			completeURI += "index.html"
		}

	}

	return filepath, completeURI
}
