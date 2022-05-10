package myerr

import (
	"os/exec"
)

// 中间模块，调用底层模块

type IntermediateErr struct {
	MyErr
}

// RunJob 运行[logID:1]:19:57:44 /home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/main.go:13: myerr.IntermediateErr{error:myerr.MyErr{Inner:myerr.LowLevelErr{error:myerr.MyErr{Inner:(*fs.PathError)(0xc000100150), Message:"stat /bad/job/binary: no such file or directory", Stacktrace:"goroutine 1 [running]:\nruntime/debug.Stack()\n\t/home/raja/code/gosdk/go1.18.1/src/runtime/debug/stack.go:24 +0x65\nConcurrently/4_massive_Concurrency/1_err/myerr.WarpError({0x4d9990, 0xc000100150}, {0xc0000181b0?, 0x0?}, {0x0?, 0x0?, 0x3?})\n\t/home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/myerr/err.go:26 +0x85\nConcurrently/4_massive_Concurrency/1_err/myerr.IsGloballyExec({0x4b9cd3?, 0x4b3a60?})\n\t/home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/myerr/lowLevel.go:17 +0x6c\nConcurrently/4_massive_Concurrency/1_err/myerr.RunJob({0x4b8119, 0x1})\n\t/home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/myerr/intermediate.go:16 +0x4a\nmain.main()\n\t/home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/main.go:20 +0x56\n", Misc:map[string]interface {}{}}}, Message:"cannot run job binary \"1\":requisite binaries not available", Stacktrace:"goroutine 1 [running]:\nruntime/debug.Stack()\n\t/home/raja/code/gosdk/go1.18.1/src/runtime/debug/stack.go:24 +0x65\nConcurrently/4_massive_Concurrency/1_err/myerr.WarpError({0x4d9a90, 0xc0000102b0}, {0x4c00da?, 0xc00004a690?}, {0xc00007ce70?, 0x542a68?, 0x7fb5fa54e758?})\n\t/home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/myerr/err.go:26 +0x85\nConcurrently/4_massive_Concurrency/1_err/myerr.RunJob({0x4b8119, 0x1})\n\t/home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/myerr/intermediate.go:18 +0x1df\nmain.main()\n\t/home/raja/workspace/go/src/Concurrently/4_massive_Concurrency/1_err/main.go:20 +0x56\n", Misc:map[string]interface {}{}}}
func RunJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := IsGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{WarpError(err, "cannot run job binary %q:requisite binaries not available", id)} // 封装底层错误
	} else if isExecutable == false {
		return IntermediateErr{WarpError(nil, "job binary is not executable")} // 中间层错误
	}
	return exec.Command(jobBinPath, "--id"+id).Run()
}
