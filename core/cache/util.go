package cache

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Etag 用来生成文件的etag内容。使用文件的last_modified_time和length来生成
func Etag(filepath string) string {
	fileinfo, err := os.Stat(filepath)

	if err != nil {
		log.Println("can not find this file")
	}

	modtime := fileinfo.ModTime()
	len := fileinfo.Size()

	return fmt.Sprintf("%x-%x", modtime.Unix(), len)
}

// LastModified 用来生成 GMT 的最后修改时间
// Syntax
// Last-Modified: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
// Link to sectionDirectives
// <day-name>
// One of "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", or "Sun" (case-sensitive).
// <day>
// 2 digit day number, e.g. "04" or "23".
// <month>
// One of "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec" (case sensitive).
// <year>
// 4 digit year number, e.g. "1990" or "2016".
// <hour>
// 2 digit hour number, e.g. "09" or "23".
// <minute>
// 2 digit minute number, e.g. "04" or "59".
// <second>
// 2 digit second number, e.g. "04" or "59".
// GMT
// Greenwich Mean Time. HTTP dates are always expressed in GMT, never in local time.
// Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT
func LastModified(filepath string) string {
	fileinfo, err := os.Stat(filepath)

	if err != nil {
		log.Println("can not find this file")
	}

	modtime := fileinfo.ModTime().UTC()

	return fmtGMT(modtime)
}

func fmtGMT(time time.Time) string {
	return fmt.Sprintf("%s, %02d %s %d %02d:%02d:%02d GMT",
		time.Weekday().String()[0:3],
		time.Day(),
		time.Month().String()[0:3],
		time.Year(),
		time.Hour(),
		time.Minute(),
		time.Second())
}
