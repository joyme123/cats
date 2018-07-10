//Package config 封装的是http的配置
package config

//Config 是http的配置文件
//addr是vhost的地址
//port是端口号
type Config struct {
	VHosts []VHost
}
