package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

var mimes map[string]string

func init() {

	mimes = make(map[string]string)

	mimes["aac"] = "audio/aac"
	mimes["abw"] = "application/x-abiword"
	mimes["arc"] = "application/octet-stream"
	mimes["avi"] = "video/x-msvideo"
	mimes["azw"] = "application/vnd.amazon.ebook"
	mimes["bin"] = "application/octet-stream"
	mimes["bz"] = "application/x-bzip"
	mimes["bz2"] = "application/x-bzip2"
	mimes["csh"] = "application/x-csh"
	mimes["css"] = "text/css"
	mimes["csv"] = "text/csv"
	mimes["doc"] = "application/msword"
	mimes["epub"] = "application/epub+zip"
	mimes["gif"] = "image/gif"
	mimes["htm"] = "text/html"
	mimes["html"] = "text/html"
	mimes["ico"] = "image/x-icon"
	mimes["ics"] = "text/calendar"
	mimes["jar"] = "application/java-archive"
	mimes["jpeg"] = "image/jpeg"
	mimes["jpg"] = "image/jpeg"
	mimes["js"] = "application/javascript"
	mimes["json"] = "application/json"
	mimes["mid"] = "audio/midi"
	mimes["midi"] = "audio/midi"
	mimes["mpeg"] = "video/mpeg"
	mimes["mpkg"] = "application/vnd.apple.installer+xml"
	mimes["odp"] = "application/vnd.oasis.opendocument.presentation"
	mimes["ods"] = "application/vnd.oasis.opendocument.spreadsheet"
	mimes["odt"] = "application/vnd.oasis.opendocument.text"
	mimes["oga"] = "audio/ogg"
	mimes["ogv"] = "video/ogg"
	mimes["ogx"] = "application/ogg"
	mimes["pdf"] = "application/pdf"
	mimes["ppt"] = "application/vnd.ms-powerpoint"
	mimes["rar"] = "application/x-rar-compressed"
	mimes["rtf"] = "application/rtf"
	mimes["sh"] = "application/x-sh"
	mimes["svg"] = "image/svg+xml"
	mimes["swf"] = "application/x-shockwave-flash"
	mimes["tar"] = "application/x-tar"
	mimes["tif"] = "image/tiff"
	mimes["tiff"] = "image/tiff"
	mimes["ttf"] = "application/x-font-ttf"
	mimes["vsd"] = "application/vnd.visio"
	mimes["wav"] = "audio/x-wav"
	mimes["weba"] = "audio/webm"
	mimes["webm"] = "video/webm"
	mimes["webp"] = "image/webp"
	mimes["woff"] = "application/x-font-woff"
	mimes["xhtml"] = "application/xhtml+xml"
	mimes["xls"] = "application/vnd.ms-excel"
	mimes["xml"] = "application/xml"
	mimes["xul"] = "application/vnd.mozilla.xul+xml"
	mimes["zip"] = "application/zip"
	mimes["3gp"] = "video/3gpp"
	mimes["3g2"] = "video/3gpp2"
	mimes["7z"] = "application/x-7z-compressed"
}

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

func GetMimeByExt(ext string) (string, bool) {
	if ctype, ok := mimes[ext]; ok {
		return ctype, true
	} else {
		return "", false
	}
}
