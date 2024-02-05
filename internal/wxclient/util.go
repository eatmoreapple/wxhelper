package wxclient

import (
	"fmt"
	"github.com/eatmoreapple/env"
	"io"
	"os"
	"path/filepath"
)

func saveToLocal(reader io.Reader) (path string, err error) {
	file, err := os.CreateTemp(env.Name("TEMP_DIR").String(), "*")
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	if _, err = io.Copy(file, reader); err != nil {
		return "", err
	}
	filename := filepath.Base(file.Name())
	// 转换为windows路径
	return fmt.Sprintf("C:\\data\\%s", filename), nil
}
