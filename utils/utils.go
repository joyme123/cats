package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
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

func FmtGMT(time time.Time) string {
	return fmt.Sprintf("%s, %02d %s %d %02d:%02d:%02d GMT",
		time.Weekday().String()[0:3],
		time.Day(),
		time.Month().String()[0:3],
		time.Year(),
		time.Hour(),
		time.Minute(),
		time.Second())
}
