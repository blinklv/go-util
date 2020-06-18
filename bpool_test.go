// bpool_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-06-18
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-06-18

package util

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferPool(t *testing.T) {
	var (
		bp = NewBufferPool()
		g  = &Group{}
	)

	for i := 0; i < 64; i++ {
		g.Go(func() interface{} {
			rdata := make([]byte, 4096)
			for j := 0; j < 10000; j++ {
				b := bp.Get()
				if b.Len() > 0 {
					return fmt.Errorf("buffer size %d > 0", b.Len())
				}
				rand.Read(rdata)
				b.Write(rdata)
				bp.Put(b)
			}
			return nil
		})
	}

	assert.Nil(t, g.Result())
}

func TestBytesPool(t *testing.T) {
	var (
		bp httputil.BufferPool = NewBytesPool(4096)
		g                      = &Group{}
	)

	for i := 0; i < 64; i++ {
		g.Go(func() interface{} {
			var (
				obj   = make([]byte, 0, 64*1000)
				w     = &bytes.Buffer{}
				rdata = make([]byte, 64)
			)

			for j := 0; j < 1000; j++ {
				b := bp.Get()
				if len(b) != 4096 {
					return fmt.Errorf("incorrect byte size: %d", len(b))
				}

				rand.Read(rdata)
				io.CopyBuffer(&trivialWriter{w}, &trivialReader{bytes.NewReader(rdata)}, b)
				if !bytes.Equal(b[:len(rdata)], rdata) {
					return fmt.Errorf("byte content != random data")
				}
				obj = append(obj, rdata...)
			}

			if !bytes.Equal(w.Bytes(), obj) {
				return fmt.Errorf("output bytes != object")
			}

			return nil
		})
	}

	assert.Nil(t, g.Result())
}

type trivialReader struct {
	r io.Reader
}

func (tr *trivialReader) Read(b []byte) (int, error) {
	return tr.r.Read(b)
}

type trivialWriter struct {
	w io.Writer
}

func (tw *trivialWriter) Write(b []byte) (int, error) {
	return tw.w.Write(b)
}
