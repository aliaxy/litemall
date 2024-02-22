package common

import "net/http"

// FilterHandle 声明处理函数
type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

// Filter 拦截器结构体
type Filter struct {
	filterMap map[string]FilterHandle
}

// NewFilter 新建拦截器实例
func NewFilter() *Filter {
	return &Filter{
		filterMap: make(map[string]FilterHandle),
	}
}

// RegisterFilterURI 注册拦截器
func (f *Filter) RegisterFilterURI(uri string, handle FilterHandle) {
	f.filterMap[uri] = handle
}

// GetFilterHandle 根据 URI 获取拦截器
func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

// WebHandle 声明 web 处理函数
type WebHandle func(rw http.ResponseWriter, req *http.Request) error

// Handle 执行拦截器 返回函数类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		for path, handle := range f.filterMap {
			if path == req.RequestURI {
				// 执行拦截业务
				err := handle(rw, req)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}

		// 正常执行
		webHandle(rw, req)
	}
}
