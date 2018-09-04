package cache

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joyme123/cats/utils"
)

// Etag 用来生成文件的etag内容。使用文件的last_modified_time和length来生成
func Etag(filepath string) string {
	fileinfo, err := os.Stat(filepath)

	if err != nil {
		log.Println("can not find this file")
		return "*"
	}

	modtime := fileinfo.ModTime()
	len := fileinfo.Size()

	return fmt.Sprintf("\"%x-%x\"", modtime.Unix(), len)
}

// CompareETag 负责将请求头中的etag和文件的etag进行对比
// reqEtag 中包含W/开头,表示采用弱比较算法
func CompareETag(reqEtag string, fileEtag string) bool {
	return reqEtag == fileEtag
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

	modtime := GMTTime(filepath)

	return utils.FmtGMT(modtime)
}

// GMTTime 用来获取文件的GMT最后修改时间
func GMTTime(filepath string) time.Time {
	fileinfo, err := os.Stat(filepath)

	if err != nil {
		log.Println("can not find this file")
		return time.Now()
	}

	modtime := fileinfo.ModTime().UTC()

	return modtime
}

// CompareFileModifiedTime 比较文件的修改时间
// 返回一个bool型，true代表迟于给定时间，false代表早于或等于
func CompareFileModifiedTime(filepath string, timeStr string) bool {
	modTime := GMTTime(filepath)

	targetTime, err := ParseGMT(timeStr)

	if err != nil {
		log.Println("error when parse time")
		return true
	}

	return modTime.After(targetTime)
}

// GMTParseError 是当GMT解析错误时返回
type GMTParseError struct {
	Layout  string
	Value   string
	Message string
}

func (e *GMTParseError) Error() string {
	if e.Message == "" {
		return "parsing time " +
			e.Value + " as " +
			e.Layout + ": cannot parse "
	}
	return "parsing time " +
		e.Value + e.Message
}

// ParseGMT 用来解析HTTP协议中的GMT格式字符串
func ParseGMT(timeStr string) (time.Time, error) {
	items := strings.Split(timeStr, " ")

	if len(items) == 6 {

		day, err := strconv.Atoi(items[1])

		if err != nil {
			return time.Time{}, err
		}

		var month time.Month
		monthStr := items[2]

		// "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
		switch monthStr {
		case "Jan":
			month = time.January
			break

		case "Feb":
			month = time.February
			break

		case "Mar":
			month = time.March
			break

		case "Apr":
			month = time.April
			break

		case "May":
			month = time.May
			break

		case "Jun":
			month = time.June
			break

		case "Jul":
			month = time.July
			break

		case "Aug":
			month = time.August
			break

		case "Sep":
			month = time.September
			break

		case "Oct":
			month = time.October
			break

		case "Nov":
			month = time.November
			break

		case "Dec":
			month = time.December
			break

		default:
			return time.Time{}, &GMTParseError{"GMT", timeStr, "can not understand this month format"}
		}

		year, err := strconv.Atoi(items[3])

		if err != nil {
			return time.Time{}, err
		}

		hmsStr := items[4]

		hms := strings.Split(hmsStr, ":")

		if len(hms) != 3 {
			return time.Time{}, &GMTParseError{"GMT", timeStr, "can not understand this hms format"}
		}

		hour, err := strconv.Atoi(hms[0])
		if err != nil {
			return time.Time{}, err
		}
		min, err := strconv.Atoi(hms[1])
		if err != nil {
			return time.Time{}, err
		}
		sec, err := strconv.Atoi(hms[2])
		if err != nil {
			return time.Time{}, err
		}

		return time.Date(year, month, day, hour, min, sec, 0, time.UTC), nil

	} else {
		return time.Time{}, &GMTParseError{"GMT", timeStr, "can not understand this time format"}
	}
}
