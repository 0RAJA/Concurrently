package main

import (
	"Concurrently/3_pattern/common"
	"bufio"
	"io"
	"log"
	"os"
	"testing"
)

//演示缓冲写入队列和未缓冲写入队列的区别

func BenchmarkUnBufferedWrite(b *testing.B) {
	performWrite(b, tmpFileOrFail())
}

func BenchmarkBufferedWrite(b *testing.B) {
	bufferFile := bufio.NewWriter(tmpFileOrFail())
	performWrite(b, bufferFile)
}

func tmpFileOrFail() *os.File {
	file, err := os.CreateTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func performWrite(b *testing.B, write io.Writer) {
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for bt := range common.ToType[byte](done, common.Take(done, common.Repeat(done, byte(1)), b.N)) {
		_, _ = write.Write([]byte{bt})
	}
}
