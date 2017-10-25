package io

import (
	"io"
	"archive/zip"
)

type ZipFlattener struct {
	r        *zip.ReadCloser
	err error
	files    []*zip.File
	currFile io.ReadCloser
}

func (z *ZipFlattener) Close() error {
	z.err =io.ErrClosedPipe
	return z.Close()
}

func (z *ZipFlattener) Read(p []byte) (n int, err error) {
	if z.err != nil{
		return 0, err
	}
	for {
		if z.currFile == nil {
			if len(z.files) == 0 {
				z.err = io.EOF
				return 0, io.EOF
			}
			nextFile := z.files[0]
			z.files = z.files[1:]
			z.currFile, err = nextFile.Open()
			if err != nil {
				z.err = err
				return 0, err
			}
		}

		n, err = z.currFile.Read(p)
		if err == io.EOF {
			z.currFile.Close()
			z.currFile = nil
			if n > 0 {
				return n, nil
			}
		} else {
			z.err = err
			return n, err
		}
	}
}

var _ io.ReadCloser = &ZipFlattener{}

func NewZipFlattener(name string) (z *ZipFlattener, err error) {
	z = &ZipFlattener{}
	z.r, err = zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	z.files = z.r.File
	return z, nil
}
