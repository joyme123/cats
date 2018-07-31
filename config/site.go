package config

// Site，用来实例化一个vhost实例
type Site struct {
	Addr       string
	Port       int
	ServerName string
	Root       string
	Index      []string
	FCGIPass   string
}
