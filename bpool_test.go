// bpool_test.go
//
// Author: blinklv <blinklv@icloud.com>
// Create Time: 2020-06-18
// Maintainer: blinklv <blinklv@icloud.com>
// Last Change: 2020-06-18

package util

import (
	"fmt"
	"math/rand"
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
