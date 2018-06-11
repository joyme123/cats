package http

type Hub struct {
	container []*Component
}

// 注册插件
func (hub *Hub) Register(comp *Component) {
	hub.container = append(hub.container, comp)
}
