// +build manual

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/paulcager/aio"
	"github.com/ulikunitz/xz"
)

const testFile = "/home/paul/Downloads/goland-173.3531.21.tar"

func main() {
	BenchmarkGoGZIP()
	BenchmarkPipe(testFile+".gz", "zcat")
	BenchmarkGoXZ()
	BenchmarkPipe(testFile+".xz", "xzcat")

}

func BenchmarkGoGZIP() {
	start := time.Now()
	r, err := os.Open(testFile + ".gz")
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

func BenchmarkPipe(file string, cmd string) {
	start := time.Now()
	r, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	gr := aio.NewPipeReader(r, cmd)

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

func BenchmarkGoXZ() {
	start := time.Now()
	r, err := os.Open(testFile + ".xz")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	gr, err := xz.NewReader(r)
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
