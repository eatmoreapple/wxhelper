package wxclient

import (
	"fmt"
	"github.com/eatmoreapple/env"
	"os"
	"path/filepath"
)

var _tempDir = env.Name("TEMP_DIR").StringOrElse(os.TempDir())

// for go:linkname
func tempDir() string {
	return _tempDir
}

// 转换为windows路径
func convertToWindows(filename string) (path string, err error) {
	if _, err := os.Stat(filepath.Join(tempDir(), filename)); err != nil {
		return "", err
	}
	return fmt.Sprintf("C:\\data\\%s", filename), nil
}
