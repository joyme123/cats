package config

type VHost struct {
	Addr      string
	Port      int
	ServeFile string
	Index     []string
}
