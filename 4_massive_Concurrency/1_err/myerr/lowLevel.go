package myerr

import (
	"os"
)

// 底层模块

type LowLevelErr struct {
	MyErr
}

// IsGloballyExec 判断是否为可执行文件
func IsGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{WarpError(err, err.Error())} // 封装错误
	}
	return info.Mode().Perm()&0100 == 0100, nil
}
