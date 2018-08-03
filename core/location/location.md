location在nginx中，是一个非常强大的组件。它会根据路由规则，调用不同的其他的模块，产生响应的输出

比如匹配.php为后缀的请求，将这种请求一律调用fastcgi_pass组件，使用php-fpm产生结果并输出。而其他文件都一律使用serveFile来输出静态的文件内容

为了实现这种机制。我的代码组织应该有以下几点：

1. serveFile组件和fastcgi组件将不由vhost进行直接调用。而是应该由location组件调用。这样location组件就有权力去决定哪个组件被调用
2. location组件拥有Hub容器。可以向其中注入组件
3. location组件可以解析正则