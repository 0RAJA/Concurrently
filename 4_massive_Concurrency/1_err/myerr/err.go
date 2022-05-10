package myerr

import (
	"fmt"
	"runtime/debug"
)

/*
	异常传递
	包括内容:
		发生了什么。发生在什么时间，发生在什么位置。自定义信息。堆栈调用信息
	封装异常，每层对异常进行封装，最外层进行判断，选取合适的异常信息展示给用户
*/

// MyErr 异常类型
type MyErr struct {
	Inner      error // 封装之前的err
	Message    string
	Stacktrace string
	Misc       map[string]interface{} // 存储各种杂项的字段，记录堆栈记录的hash，并发ID等
}

func WarpError(err error, msgf string, msgArgs ...interface{}) MyErr {
	return MyErr{
		Inner:      err,
		Message:    fmt.Sprintf(msgf, msgArgs...),
		Stacktrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

func (err MyErr) Error() string {
	return err.Message
}

func (err MyErr) Unwrap() error {
	return err.Inner
}
