package config

// VHost，用来实例化一个server实例
type VHost struct {
	Addr       string
	Port       int
	ServerName string
	ServeFile  string
	Index      []string
}
