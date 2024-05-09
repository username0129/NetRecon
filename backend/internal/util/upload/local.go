package upload

import (
	"backend/internal/global"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Local struct{}

func EncryptMD5(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

// UploadFile 本地文件上传
func (l *Local) UploadFile(file *multipart.FileHeader) (string, error) {
	// 读取文件后缀
	ext := filepath.Ext(file.Filename)
	// 仅允许上传 jpg, png, txt 文件
	if ext != ".jpg" && ext != ".png" && ext != ".txt" {
		return "", fmt.Errorf("仅支持 .jpg .png .txt: %s", ext)
	}

	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = EncryptMD5([]byte(name))

	// 使用 UUID 生成新的文件名
	filename := uuid.Must(uuid.NewV4()).String() + ext
	// 拼接路径
	path := "uploads/file/" + filename

	if err := os.MkdirAll("uploads/file", os.ModePerm); err != nil {
		global.Logger.Error("创建上传目录失败", zap.Error(err))
		return "", errors.New("创建上传目录失败")
	}

	// 读取文件内容
	f, err := file.Open()
	if err != nil {
		global.Logger.Error("读取文件内容失败", zap.Error(err))
		return "", errors.New("读取文件内容失败")
	}
	defer f.Close()

	// 创建文件
	out, err := os.Create(path)
	if err != nil {
		global.Logger.Error("创建文件失败", zap.Error(err))
		return "", errors.New("创建文件失败")
	}
	defer out.Close()

	if _, err = io.Copy(out, f); err != nil {
		global.Logger.Error("写入文件内容失败", zap.Error(err))
		return "", errors.New("写入文件内容失败")
	}
	return path, nil
}

// DeleteFile 本地文件删除
func (l *Local) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		global.Logger.Error("本地文件删除失败: ", zap.Error(err))
		return errors.New("本地文件删除失败")
	}
	return nil
}
