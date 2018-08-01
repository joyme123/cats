package http

type Hub struct {
	// FIXME: 这里应该注入指针，这样才能保证是同一个组件对象
	container []Component
}

// 注册插件
func (hub *Hub) Register(comp Component) {
	hub.container = append(hub.container, comp)
}
