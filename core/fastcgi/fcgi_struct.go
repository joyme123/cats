// 在FastCGI协议中，传输的数据都是大端序的，因此在将uint32的值序列化成bytes时，要注意序列化的值的大小端序
// 可以参考这篇文章，使用go语言的系统包进行处理
// http://lihaoquan.me/2016/11/5/golang-byteorder.html

package fastcgi

import (
	"encoding/binary"
)

// 所有的记录都实现了ToBlob方法
type Record interface {
	ToBlob() []byte
}

// FCGIListenSockFileNo 是监听的文件描述符
const FCGIListenSockFileNo = 0

// FCGIHeader 是 FCGI Record的头部
type FCGIHeader struct {
	Version       uint8
	Type          uint8
	RequestID     uint16
	ContentLength uint16
	PaddingLength uint8
	Reserved      uint8
}

// FCGIHeaderLen 是FCGIHeader的字节数
const FCGIHeaderLen = 8

// FCGIVersion1 是FCGIHeader的Version
const FCGIVersion1 = 1

// FCGIHeader中type的取值
const (
	FCGIBeginRequest    = 1
	FCGIAbortRequest    = 2
	FCGIEndRequest      = 3
	FCGIParams          = 4
	FCGIStdin           = 5
	FCGIStdout          = 6
	FCGIStderr          = 7
	FCGIData            = 8
	FCGIGetValues       = 9
	FCGIGetValueResults = 10
	FCGIUnknownType     = 11
	FCGIMaxType         = 11
)

// FCGINullRequestID 是FCGIHeader中requestId的默认取值
const FCGINullRequestID = 0

// FCGIBeginRequestBody 是 FCGI 第一个请求的请求体
type FCGIBeginRequestBody struct {
	Role     uint16
	Flags    uint8
	Reserved [5]uint8
}

// FCGIBeginRequestRecord 是 FCGI的第一个请求
type FCGIBeginRequestRecord struct {
	Header FCGIHeader
	Body   FCGIBeginRequestBody
}

// FCGIKeepConn 是FCGIBeginRequestBody的Flags的取值
const FCGIKeepConn = 1

// FCGIBeginRequestBody的Role取值
const (
	FCGIResponder  = 1 // 响应器
	FCGIAuthorizer = 2 // 认证器
	FCGIFilter     = 3 // 过滤器
)

// FCGIEndRequestBody 是 FCGI的最后一个请求的请求体
type FCGIEndRequestBody struct {
	AppStatus      uint32
	ProtocolStatus uint8
	Reserved       [3]uint8
}

// FCGIEndRequestRecord 是 FCGI的最后一个请求记录
type FCGIEndRequestRecord struct {
	Header FCGIHeader
	Body   FCGIEndRequestBody
}

// FCGIEndRequestBody中ProtocalStatus的取值
const (
	FCGIRequestComplete = 0
	FCGICantMpxConn     = 1
	FCGIOverloaded      = 2
	FCGIUnknownRole     = 3
)

// FCGI_GET_VALUES / FCGI_GET_VALUES_RESULT 记录的名字
const (
	FCGIMaxConns = "FCGI_MAX_CONNS"
	FCGIMaxReqs  = "FCGI_MAX_REQS"
	FCGIMpxConns = "FCGI_MPXS_CONNS"
)

// FCGIUnknownTypeBody 未知响应响应体
type FCGIUnknownTypeBody struct {
	Type     uint8
	Reserved [7]uint8
}

// FCGIUnknownTypeRecord 未知响应记录
type FCGIUnknownTypeRecord struct {
	Header FCGIHeader
	Body   FCGIUnknownTypeBody
}

// FCGINameValuePair11 是FCGI_PARAMS的键值对封装，这是键和值长度都是在127之内的情况
type FCGINameValuePair11 struct {
	NameLength  uint8 /* nameLength  >> 7 == 0 */
	ValueLength uint8 /* valueLength >> 7 == 0 */
	NameData    []byte
	ValueData   []byte
}

// FCGINameValuePair14 是FCGI_PARAMS的键值对封装，这是键长在127之内,值长大于127的情况
type FCGINameValuePair14 struct {
	NameLength  uint8  /* nameLength  >> 7 == 0 */
	ValueLength uint32 /* valueLength >> 31 == 1 */

	NameData  []byte
	ValueData []byte
}

// FCGINameValuePair41 是FCGI_PARAMS的键值对封装，这是键长大于127,值长在127之内的情况
type FCGINameValuePair41 struct {
	NameLength  uint32 /* nameLength  >> 31 == 1 */
	ValueLength uint8  /* valueLength >> 7 == 0 */
	NameData    []byte
	ValueData   []byte
}

// FCGINameValuePair44 是FCGI_PARAMS的键值对封装，这是键值长大于127的情况
type FCGINameValuePair44 struct {
	NameLength  uint32 /* nameLength  >> 31 == 1 */
	ValueLength uint32 /* valueLength >> 31 == 1 */
	NameData    []byte
	ValueData   []byte
}

// FCGIParamsRecord 是 FCGI_PARAMS 记录
type FCGIParamsRecord struct {
	Header FCGIHeader
	Body   Record
}

// FCGIStdioRecord 是标准输入流
type FCGIStdioRecord struct {
	Header FCGIHeader
	Body   []byte
}

// ToBlob 将请求头转换成二进制流
func (header *FCGIHeader) ToBlob() []byte {
	headerBytes := make([]byte, 8, 8)
	headerBytes[0] = header.Version
	headerBytes[1] = header.Type
	binary.BigEndian.PutUint16(headerBytes[2:4], header.RequestID)
	binary.BigEndian.PutUint16(headerBytes[4:6], header.ContentLength)
	headerBytes[6] = header.PaddingLength
	headerBytes[7] = header.Reserved

	return headerBytes
}

// ToBlob 将请求开始记录转换成二进制流
func (record *FCGIBeginRequestRecord) ToBlob() []byte {
	headerBytes := record.Header.ToBlob()

	bodyBytes := make([]byte, 8, 8)

	binary.BigEndian.PutUint16(bodyBytes[0:2], record.Body.Role)
	bodyBytes[2] = record.Body.Flags

	blob := append(headerBytes, bodyBytes...)

	return blob
}

// ToBlob 将结束请求转换成二进制流
func (record *FCGIEndRequestRecord) ToBlob() []byte {
	headerBytes := record.Header.ToBlob()

	bodyBytes := make([]byte, 8, 8)

	binary.BigEndian.PutUint32(bodyBytes[0:4], record.Body.AppStatus)
	bodyBytes[4] = record.Body.ProtocolStatus

	blob := append(headerBytes, bodyBytes...)

	return blob
}

// ToBlob 将FCGI_PARAMS转换成二进制流
func (params *FCGIParamsRecord) ToBlob() []byte {
	headerBytes := params.Header.ToBlob()
	bodyBytes := params.Body.ToBlob()

	blob := append(headerBytes, bodyBytes...)

	return blob
}

// ToBlob FCGINameValuePair11
func (pair *FCGINameValuePair11) ToBlob() []byte {
	totalLen := pair.NameLength + pair.ValueLength + 2

	blob := make([]byte, 2, totalLen)

	blob[0] = pair.NameLength
	blob[1] = pair.ValueLength

	blob = append(blob, []byte(pair.NameData)...)

	blob = append(blob, []byte(pair.ValueData)...)

	return blob
}

// ToBlob FCGINameValuePair14
func (pair *FCGINameValuePair14) ToBlob() []byte {
	totalLen := uint32(pair.NameLength) + pair.ValueLength + 5

	blob := make([]byte, 5, totalLen)

	blob[0] = pair.NameLength
	binary.BigEndian.PutUint32(blob[1:5], pair.ValueLength)

	blob = append(blob, []byte(pair.NameData)...)

	blob = append(blob, []byte(pair.ValueData)...)

	return blob
}

// ToBlob FCGINameValuePair41
func (pair *FCGINameValuePair41) ToBlob() []byte {

	totalLen := pair.NameLength + uint32(pair.ValueLength) + 5

	blob := make([]byte, 5, totalLen)

	binary.BigEndian.PutUint32(blob[0:4], pair.NameLength)
	blob[4] = pair.ValueLength

	blob = append(blob, []byte(pair.NameData)...)

	blob = append(blob, []byte(pair.ValueData)...)

	return blob

}

// ToBlob FCGINameValuePair44
func (pair *FCGINameValuePair44) ToBlob() []byte {
	totalLen := pair.NameLength + uint32(pair.ValueLength) + 8

	blob := make([]byte, 8, totalLen)

	binary.BigEndian.PutUint32(blob[0:4], pair.NameLength)
	binary.BigEndian.PutUint32(blob[4:8], pair.ValueLength)

	blob = append(blob, []byte(pair.NameData)...)

	blob = append(blob, []byte(pair.ValueData)...)

	return blob
}

// ToBlob FCGIUnknownTypeRecord
func (record *FCGIUnknownTypeRecord) ToBlob() []byte {
	headerBytes := record.Header.ToBlob()

	bodyBytes := make([]byte, 8, 8)
	bodyBytes[0] = record.Body.Type

	blob := append(headerBytes, bodyBytes...)

	return blob
}

// ToBlob FCGIStdinRecord 标准输入记录
func (record *FCGIStdioRecord) ToBlob() []byte {
	headerBytes := record.Header.ToBlob()

	blob := append(headerBytes, record.Body...)

	return blob
}
