package wxclient

import (
	"github.com/eatmoreapple/env"
	"io"
	"os"
)

func readerToFile(reader io.Reader) (file *os.File, cb func(), err error) {
	var ok bool
	if file, ok = reader.(*os.File); ok {
		return file, func() {}, nil
	}
	file, err = os.CreateTemp(env.Name("TEMP_DIR").String(), "*")
	if err != nil {
		return nil, nil, err
	}
	cb = func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		cb()
		return nil, nil, err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		cb()
		return nil, nil, err
	}
	return file, cb, nil
}
