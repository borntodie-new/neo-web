package neo

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type StaticFile struct {
	Dir  string // 需要开放的文件路径
	Path string // 参数地址
}

func NewStaticFile(dir string, path string) *StaticFile {
	return &StaticFile{Dir: dir, Path: path}
}

func (s *StaticFile) Handler() HandlerFunc {
	return func(ctx *Context) {
		// 拿到URL中的文件名
		fileName := ctx.Params(s.Path)
		// 拼接文件地址并打开
		filePath := filepath.Join(s.Dir, fileName)
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "服务器异常")
			return
		}
		ctx.String(http.StatusOK, string(data))
	}
}
