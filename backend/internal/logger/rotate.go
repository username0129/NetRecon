package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Rotate struct {
	level    string        // 日志级别
	format   string        // 时间格式
	director string        // 日志目录
	file     *os.File      // 文件句柄
	mutex    *sync.RWMutex // 读写锁
}

func NewRotate(level string, format string, director string) *Rotate {
	return &Rotate{
		level:    level,
		format:   format,
		director: director,
		mutex:    new(sync.RWMutex)}
}

func (r *Rotate) Write(bytes []byte) (n int, err error) {
	r.mutex.Lock() // 上锁
	defer func() {
		if r.file != nil {
			_ = r.file.Close()
			r.file = nil
		}
		r.mutex.Unlock()
	}() // 解锁和释放句柄

	formats := make([]string, 0, 3)                        // 日志路径：目录/时间/文件名
	formats = append(formats, r.director)                  // 日志目录
	formats = append(formats, time.Now().Format(r.format)) // 时间
	formats = append(formats, r.level+".log")              // 日志文件名

	filename := filepath.Join(formats...)
	dirname := filepath.Dir(filename)
	if err = os.MkdirAll(dirname, 0755); err != nil {
		return 0, err
	} // 根据路径新建文件夹

	if r.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		return 0, err
	} // 打开文件句柄

	return r.file.Write(bytes)
}
