package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net"
	"testing"
)

func init() {
	wg := startNetWorkDaemon()
	wg.Wait()
}

func BenchmarkT1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		require.NoError(b, err)
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatal("cannot read", err.Error())
		}
	}
}
