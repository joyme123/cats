package serveFile

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Range struct {
	start int
	end   int
}

type RangeData struct {
	unit  string
	parts []Range
}

// parseRange 负责解析range请求头的指令
func parseRange(rangeValue []byte) (RangeData, error) {
	var rangeData RangeData
	var token bytes.Buffer
	var part = Range{0, -1}
	var err error
	for _, b := range rangeValue {
		switch b {
		case '=':
			//token匹配结束
			rangeData.unit = token.String()
			token.Reset()
			break

		case '-':
			part.start, err = strconv.Atoi(token.String())
			token.Reset()
			if err != nil {
				// 解析有错
				log.Printf("解析range出错：%v\n", token.String())
				return rangeData, errors.New("parse range start error")
			}
			break

		case ',':
			part.end, err = strconv.Atoi(token.String())
			token.Reset()
			if err != nil {
				// 解析有错
				log.Printf("解析range出错：%v\n", token.String())
				return rangeData, errors.New("parse range end error")
			}
			rangeData.parts = append(rangeData.parts, part)
			part = Range{0, -1}
			break
		default:
			token.WriteByte(b)
		}
	}
	if token.Len() != 0 {
		part.end, err = strconv.Atoi(token.String())
		token.Reset()
		if err != nil {
			// 解析有错
			log.Printf("解析range出错：%v\n", token.String())
			return rangeData, errors.New("parse range end error")
		}
	}

	rangeData.parts = append(rangeData.parts, part)

	return rangeData, nil
}

// 根据整段的内容以及range数组生成响应体
// data 是请求的内容字节数组
// parts 是请求的分段
// contentType 内容mime
func mergeMultiRange(data []byte, parts []Range, contentType string) ([]byte, string) {
	var boundary = "dasdasdads"
	var res bytes.Buffer
	bytesLen := len(data)

	var firstLoop bool = true

	for _, part := range parts {

		if firstLoop {
			res.WriteString("--" + boundary + "\r\n")
			firstLoop = false
		} else {
			res.WriteString("\r\n--" + boundary + "\r\n")
		}

		res.WriteString("Content-Type: " + contentType + "\r\n")

		contentRange := fmt.Sprintf("Content-Range: bytes %d-%d/%d\r\n\r\n", part.start, part.end, bytesLen)

		res.WriteString(contentRange)
		if part.end == -1 {
			res.Write(data[part.start:])
		} else {
			res.Write(data[part.start:part.end])
		}
	}

	res.WriteString("\r\n--" + boundary + "--\r\n")

	return res.Bytes(), boundary
}
