// +build manual

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/paulcager/aio"
)

const testFile = "/home/paul/Downloads/goland-173.3531.21.tar.gz"

func main() {
	BenchmarkGoGZIP()
	BenchmarkLinuxGZIP()
	BenchmarkGoGZIP()
	BenchmarkLinuxGZIP()
}

func BenchmarkGoGZIP() {
	start := time.Now()
	r, err := os.Open(testFile)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	gr, err := gzip.NewReader(r)
	if err != nil {
		panic(err)
	}

	buff := make([]byte, 2048)
	cnt := 0
	for {
		n, err := gr.Read(buff)
		cnt = cnt + n
		if err != nil {
			break
		}
	}
	fmt.Println("Go", cnt, time.Since(start))
}

func BenchmarkLinuxGZIP() {
	start := time.Now()
	r, err := os.Open(testFile)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	gr := aio.NewGZIPReader(r)

	cnt := 0
	buff := make([]byte, 2048)
	for {
		n, err := gr.Read(buff)
		cnt = cnt + n
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("pipe", cnt, time.Since(start))
}
