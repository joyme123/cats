package hub

type Box struct {
	container []*Controller
}

// 注册插件
func (box *Box) Register(plugin *Controller) {
	box.container = append(box.container, plugin)
}
