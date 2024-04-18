package wxclient

import (
	"fmt"
	"github.com/eatmoreapple/env"
	"os"
	"path/filepath"
)

// 转换为windows路径
func convertToWindows(filename string) (path string, err error) {
	tempDir := env.Name("TEMP_DIR").String()
	if _, err := os.Stat(filepath.Join(tempDir, filename)); err != nil {
		return "", err
	}
	return fmt.Sprintf("C:\\data\\%s", filename), nil
}
