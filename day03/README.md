# 路由树
关于路由树，我们使用的是前缀树，具体关于前缀树的相关信息这里就不再赘述，大家自行查阅资料
## 静态匹配
关于静态匹配，就是完完全全一样。并且它的匹配优先级是最高的

目前发现一个小问题
```text
注册路由：/login
请求地址：/login/123
像这种路由，我们需不需要支持，不考虑RESTFul风格
接下来我们看看Gin支持吗？
好吧，Gin不支持这种，那我们也不支持
```
解决方案：在命中到路由后，再判断匹配出来的node节点中的pattern字段和HTTP的URL是否匹配
一样就表示真的命中了，否则就匹配失败
```go
// 匹配路由
func (r *router) getRouter(method string, pattern string) *node {
	root, ok := r.roots[method]
	if !ok {
		// 路由树都不存在，直接返回nil
		return nil
	}
	// 切割pattern
	pattern = strings.Trim(pattern, "/")
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if part == "" {
			return nil
		}
		child := root.search(part)
		if child == nil {
			return nil
		}
		// TODO 关键就是这个判断
		if child.pattern != "" && child.pattern == fmt.Sprintf("/%s", pattern) {
			return child
		}
		root = child
	}
	return nil
}
```
## 参数路由匹配
参数路由匹配类似这种：`/study/:lang`、`/study/:lang/:id`等等

参数路由需要修改Trie数据的搜索逻辑，并且搜索还是有优先级的`精确匹配>模糊匹配[参数路由，模糊路由]`
```go
// 查询子节点是否含有part节点
// [/ , user, home]
func (n *node) search(part string) *node {
	if n.children == nil {
		n.children = make([]*node, 0)
	}
	for _, child := range n.children {
		// 重点需要关注这个优先级问题
		// 精确匹配，优先级高
		if child.part == part {
			return child
		}
		// 模糊匹配，优先级低
		if child.isWild {
			return child
		}
	}
	return nil
}
```

之前在实现注册路由和匹配路由的时候忘记处理根路由的情况，这里统一处理下
```go
// 注册路由
func addRouter(method, pattern string, handlerFunc HandlerFunc){
	......
    // 特殊处理根路由，必须放在处理下面两个逻辑之前
    if pattern == "/" {
        root.pattern = pattern
        key := fmt.Sprintf("%s-%s", method, pattern)
        r.handlers[key] = handlerFunc
        log.Printf("Add Router %4s - %s", method, pattern)
        return
    }
    // pattern 必须以 / 开头
    if !strings.HasPrefix(pattern, "/") {
        panic("web: 路由必须以 / 开头")
    }
    if strings.HasSuffix(pattern, "/") {
        panic("web: 路由不能以 / 结尾")
    }
	......
}

// 匹配路由
func getRouter(mathod, pattern string) (*node, map[string]string) {
    params := make(map[string]string)
    root, ok := r.roots[method]
    if !ok {
        // 路由树都不存在，直接返回nil
        return nil, params
    }
    // 特殊处理 根路由
    if pattern == "/" {
        return root, params
    }
	......
}
```

回到正题，如何支持参数路由匹配。前面我们已经改了Trie树的搜索逻辑了。接下来就是`router`树的注册和匹配逻辑了。

这里我们分成仔细想想，像这类的参数路由，和之前的静态路由有什么不同之处吗？或者说需要额外做哪些操作
1. 注册的时候，如果是带有`:`的part，需要给这个节点的`isWild`属性设置成`true`
    ```go
    // 内部核心API，仅共内部使用，用于注册路由
    // 注册路由
    // method = GET | pattern = "/user/home" => parts = [user, home]
    // method = GET | pattern = "/user/order" => parts = [user, order]
    func (r *router) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
        root, ok := r.roots[method]
        if !ok { // 根路由树不存在
            root = &node{
                part: "/",
            }
            r.roots[method] = root
        }
        // 特殊处理根路由
        if pattern == "/" {
            root.pattern = pattern
            key := fmt.Sprintf("%s-%s", method, pattern)
            r.handlers[key] = handlerFunc
            log.Printf("Add Router %4s - %s", method, pattern)
            return
        }
        // pattern 必须以 / 开头
        if !strings.HasPrefix(pattern, "/") {
            panic("web: 路由必须以 / 开头")
        }
        if strings.HasSuffix(pattern, "/") {
            panic("web: 路由不能以 / 结尾")
        }
        // 切割pattern
        parts := strings.Split(pattern[1:], "/")
        for _, part := range parts {
            if part == "" {
                panic("web: 路由不能连续出现 / ")
            }
            child := root.search(part)
            if child == nil { // 返回值是nil，说明当前节点下没有part的节点，需要新增
                child = &node{ // 这里不需要添加pattern字段，只有最底层的叶子节点才需要添加。我们统一在循环之后再添加
                    part: part,
                }
                // 如果part中带有 : 表示是模糊匹配，需要给节点的isWild设为true。核心就是这里
                if strings.HasPrefix(part, ":") {
                    child.isWild = true
                }
                root.children = append(root.children, child)
            }
            root = child // 查找到了child，沿着child继续查找
        }
        root.pattern = pattern
        key := fmt.Sprintf("%s-%s", method, pattern)
        r.handlers[key] = handlerFunc
        log.Printf("Add Router %4s - %s", method, pattern)
    }
    
    ```

2. 匹配的时候，需要判断child节点的`isWild`是否为`true`，如果是，需要拿到child的`part`字段和当前的part信息做映射，保存到map中
    ```go
    // 匹配路由
    func (r *router) getRouter(method string, pattern string) (*node, map[string]string) {
        params := make(map[string]string)
        root, ok := r.roots[method]
        if !ok {
            // 路由树都不存在，直接返回nil
            return nil, params
        }
        // 特殊处理 根路由
        if pattern == "/" {
            return root, params
        }
        // 切割pattern
        pattern = strings.Trim(pattern, "/")
        parts := strings.Split(pattern, "/")
        for _, part := range parts {
            if part == "" {
                return nil, params
            }
            child := root.search(part)
            if child == nil {
                return nil, params
            }
            // 将匹配到的带有:的路由添加到params中
            if child.isWild {
                params[child.part[1:]] = part
            }
            // 1. child.pattern != ""这是精确匹配，但是还不够，有点缺陷
            // 2. child.pattern == fmt.Sprintf("/%s", pattern) 这是精确匹配，搭配条件1才是完美
            // 上述两个条件只有配合使用才能精确匹配路由。/login 和 /login/123 两种路由都能匹配到
            // 3. child.isWild 这是模糊匹配，优先级低于精确匹配
            if (child.pattern != "" && child.pattern == fmt.Sprintf("/%s", pattern)) || child.isWild {
                return child, params
            }
            root = child
        }
        return nil, params
    }
    ```
3. 其实这个功能注册不难，难的是匹配

## '*'路由匹配
这个功能和上面的参数路由类似，复杂的是在匹配路由的时候需要额外在多做些操作
1. 路由注册
    ```go
    // 内部核心API，仅共内部使用，用于注册路由
    // 注册路由
    // method = GET | pattern = "/user/home" => parts = [user, home]
    // method = GET | pattern = "/user/order" => parts = [user, order]
    func (r *router) addRouter(method string, pattern string, handlerFunc HandlerFunc) {
        root, ok := r.roots[method]
        if !ok { // 根路由树不存在
            root = &node{
                part: "/",
            }
            r.roots[method] = root
        }
        // 特殊处理根路由
        if pattern == "/" {
            root.pattern = pattern
            key := fmt.Sprintf("%s-%s", method, pattern)
            r.handlers[key] = handlerFunc
            log.Printf("Add Router %4s - %s", method, pattern)
            return
        }
        // pattern 必须以 / 开头
        if !strings.HasPrefix(pattern, "/") {
            panic("web: 路由必须以 / 开头")
        }
        if strings.HasSuffix(pattern, "/") {
            panic("web: 路由不能以 / 结尾")
        }
        // 切割pattern
        parts := strings.Split(pattern[1:], "/")
        for _, part := range parts {
            if part == "" {
                panic("web: 路由不能连续出现 / ")
            }
            child := root.search(part)
            if child == nil { // 返回值是nil，说明当前节点下没有part的节点，需要新增
                child = &node{ // 这里不需要添加pattern字段，只有最底层的叶子节点才需要添加。我们统一在循环之后再添加
                    part: part,
                }
                // 如果part中带有 : 表示是模糊匹配，需要给节点的isWild设为true 关键在这里
                if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*") {
                    child.isWild = true
                }
                root.children = append(root.children, child)
            }
            if child.isWild {
                panic("web: 模糊匹配冲突")
            }
            root = child // 查找到了child，沿着child继续查找
        }
        root.pattern = pattern
        key := fmt.Sprintf("%s-%s", method, pattern)
        r.handlers[key] = handlerFunc
        log.Printf("Add Router %4s - %s", method, pattern)
    }
    
    ```
2. 路由匹配
    ```go
    // 匹配路由
    func (r *router) getRouter(method string, pattern string) (*node, map[string]string) {
        params := make(map[string]string)
        root, ok := r.roots[method]
        if !ok {
            // 路由树都不存在，直接返回nil
            return nil, params
        }
        // 特殊处理 根路由
        if pattern == "/" {
            return root, params
        }
        // 切割pattern
        pattern = strings.Trim(pattern, "/")
        parts := strings.Split(pattern, "/")
        for _, part := range parts {
            if part == "" {
                return nil, params
            }
            child := root.search(part)
            if child == nil {
                return nil, params
            }
            // 将匹配到的带有:的路由添加到params中
            // 关键在这里
            if child.isWild {
                // 完成*的模糊匹配
                // 这也是由优先级的
                // : 优先级高于 *
                if strings.HasPrefix(child.part, ":") {
                    params[child.part[1:]] = part
                } else if strings.HasPrefix(child.part, "*") {
                    params[child.part[1:]] = pattern[strings.Index(pattern, part):]
                }
            }
            // 1. child.pattern != ""这是精确匹配，但是还不够，有点缺陷
            // 2. child.pattern == fmt.Sprintf("/%s", pattern) 这是精确匹配，搭配条件1才是完美
            // 上述两个条件只有配合使用才能精确匹配路由。/login 和 /login/123 两种路由都能匹配到
            // 3. child.isWild 这是模糊匹配，优先级低于精确匹配
            if (child.pattern != "" && child.pattern == fmt.Sprintf("/%s", pattern)) || child.isWild {
                return child, params
            }
            root = child
        }
        return nil, params
    }
    
    ```
   
## 正则匹配
后续在实现。。。
## Bug修补
1. `/user/login`和`/user/:action`，可以同时注册，并且第一个优先级高于第二个
2. `/user/:action`和`/user/*filepath`，不可以同时注册，直接报错就好
3. `/user/:action`和`/assets/*filepath`，可以同时注册

情况1这种情况我们已经处理好了，就是添加一个优先级

情况2这种情况我们需要这样处理
1. `/user/:action`注册成功，注册`/user/*filepath`的时候会判断匹配到的节点的`isWild`是否是`true`，是就直接报错，因为之前已经注册过类似的

情况3是两个完全不一样的路由，没什么好处理的