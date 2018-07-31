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

// 一个字节编码的最长长度
const FCGIPairOneByte = 127

// 4个字节编码的最长长度
const FCGIPairFourBytes = 0x0fffffff

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
// 之所以是127，因为一个字节的首位是标志位
type FCGINameValuePair11 struct {
	NameLength  uint8 /* nameLength  >> 7 == 0 */
	ValueLength uint8 /* valueLength >> 7 == 0 */
	NameData    []byte
	ValueData   []byte
}

// FCGINameValuePair14 是FCGI_PARAMS的键值对封装，这是键长在127之内,值长大于127的情况
// 注意这里的ValueLength的范围是129~2^31次方，最高位是一个标志位
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
// Params是可以在多个Record中存在的，如果一个Record记录不下，把多出的部分放在下一个Record中即可
type FCGIParamsRecord struct {
	Header FCGIHeader
	Body   []byte // 因为contentLength是16位，所以body的最大长度是65536
}

// FCGIStdioRecord 是标准输入输出流
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

// New 初始化开始请求记录
// @param requestID  请求ID
// @param role       FastCGI程序扮演的角色
func (record *FCGIBeginRequestRecord) New(requestID uint16, role uint16) {
	record.Header.Version = FCGIVersion1
	record.Header.Type = FCGIBeginRequest
	record.Header.RequestID = requestID
	record.Header.ContentLength = 8
	record.Header.PaddingLength = 0
	record.Body.Role = role
	record.Body.Flags = FCGIKeepConn
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

// New 初始化结束请求记录
// @param version    		协议版本号，填1
// @param requestID  		请求ID
// @param appStatus     	应用程序的状态码，由应用程序自己定义
// @param protocolStatus	协议状态
func (record *FCGIEndRequestRecord) New(requestID uint16, appStatus uint32, protocolStatus uint8) {
	record.Header.Version = FCGIVersion1
	record.Header.Type = FCGIEndRequest
	record.Header.RequestID = requestID
	record.Header.ContentLength = 8
	record.Header.PaddingLength = 0
	record.Body.AppStatus = appStatus
	record.Body.ProtocolStatus = protocolStatus
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

// New FCGI参数记录的初始化函数
func (record *FCGIParamsRecord) New(requestID uint16, pair map[string]string) {
	record.Header.Version = FCGIVersion1
	record.Header.Type = FCGIParams
	record.Header.RequestID = requestID
	// TODO: padding length 可能要修改
	record.Header.PaddingLength = 0

	bodyBytes := make([]byte, 0, 0)

	// 处理记录的body
	for key, value := range pair {
		keyBytes := []byte(key)
		valueBytes := []byte(value)

		keyLen := len(keyBytes)
		valueLen := len(valueBytes)

		if keyLen <= 127 {
			// pair1x
			if valueLen <= 127 {
				// pair11
				var pair11 FCGINameValuePair11
				pair11.NameLength = uint8(keyLen)
				pair11.ValueLength = uint8(valueLen)
				pair11.NameData = []byte(key)
				pair11.ValueData = []byte(value)
				bodyBytes = append(bodyBytes, pair11.ToBlob()...)
			} else {
				// pair14
				var pair14 FCGINameValuePair14
				pair14.NameLength = uint8(keyLen)
				pair14.ValueLength = uint32(valueLen)
				pair14.NameData = []byte(key)
				pair14.ValueData = []byte(value)
				bodyBytes = append(bodyBytes, pair14.ToBlob()...)
			}
		} else {
			// pair4x
			if valueLen <= 127 {
				// pair41
				var pair41 FCGINameValuePair41
				pair41.NameLength = uint32(keyLen)
				pair41.ValueLength = uint8(valueLen)
				pair41.NameData = []byte(key)
				pair41.ValueData = []byte(value)
				bodyBytes = append(bodyBytes, pair41.ToBlob()...)
			} else {
				// pair44
				var pair44 FCGINameValuePair44
				pair44.NameLength = uint32(keyLen)
				pair44.ValueLength = uint32(valueLen)
				pair44.NameData = []byte(key)
				pair44.ValueData = []byte(value)
				bodyBytes = append(bodyBytes, pair44.ToBlob()...)
			}
		}
	} // end for keyvalue

	// FIXME: 判断bodyBytes的长度，根据该长度去决定是否生成多条FCGIParamsRecord。
	// 单条FCGIParamsRecord的ContentLength最多65536
	record.Body = bodyBytes
	record.Header.ContentLength = uint16(len(bodyBytes))
}

// ToBlob 将FCGI_PARAMS转换成二进制流
func (record *FCGIParamsRecord) ToBlob() []byte {
	headerBytes := record.Header.ToBlob()
	bodyBytes := record.Body

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
	binary.BigEndian.PutUint32(blob[1:5], pair.ValueLength|0x80000000) // 或0x80000000是为了将首位置1

	blob = append(blob, []byte(pair.NameData)...)

	blob = append(blob, []byte(pair.ValueData)...)

	return blob
}

// ToBlob FCGINameValuePair41
func (pair *FCGINameValuePair41) ToBlob() []byte {

	totalLen := pair.NameLength + uint32(pair.ValueLength) + 5

	blob := make([]byte, 5, totalLen)

	binary.BigEndian.PutUint32(blob[0:4], pair.NameLength|0x80000000)
	blob[4] = pair.ValueLength

	blob = append(blob, []byte(pair.NameData)...)

	blob = append(blob, []byte(pair.ValueData)...)

	return blob

}

// ToBlob FCGINameValuePair44
func (pair *FCGINameValuePair44) ToBlob() []byte {
	totalLen := pair.NameLength + uint32(pair.ValueLength) + 8

	blob := make([]byte, 8, totalLen)

	binary.BigEndian.PutUint32(blob[0:4], pair.NameLength|0x80000000)
	binary.BigEndian.PutUint32(blob[4:8], pair.ValueLength|0x80000000)

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

// New FCGIStdioRecord的初始化函数
func (record *FCGIStdioRecord) New(requestID uint16, data []byte) {
	record.Header.Version = FCGIVersion1
	record.Header.Type = FCGIStdin
	record.Header.RequestID = requestID
	// TODO: padding length 可能要修改
	record.Header.PaddingLength = 0
	record.Header.ContentLength = uint16(len(data))

	record.Body = data
}

// ToBlob FCGIStdinRecord 标准输入记录
func (record *FCGIStdioRecord) ToBlob() []byte {
	headerBytes := record.Header.ToBlob()

	blob := append(headerBytes, record.Body...)

	return blob
}
