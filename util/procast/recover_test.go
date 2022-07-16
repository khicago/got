package procast_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/khicago/got/util/procast"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	err := func() (err error) {
		h := procast.GetRewriteErrHandler(&err)
		defer procast.Recover(h)
		// defer h.Recover()
		func() {
			panic("panic 1") // should keep at line 19.
		}()
		return
	}()

	assert.NotNil(t, err, "recover failed")
	assert.Equal(t, "panic 1 [ panic !!! [func:procast_test.TestRecover.func1.1:19] ]", err.Error(), "recover content failed")
}

func TestSafeGo(t *testing.T) {
	wg := &sync.WaitGroup{}
	err := func() (err error) {
		wg.Add(1)
		procast.GetRewriteErrHandler(&err).SafeGo(func() {
			defer wg.Done()
			panic("panic 2")
		})

		wg.Wait()
		time.Sleep(time.Millisecond)
		return
	}()

	if !assert.NotNil(t, err, "SafeGo handle error failed") {
		return
	}

	// fmt.Println(err.Error())
	assert.True(t, strings.HasPrefix(err.Error(), "panic 2 [ panic !!!"), "SafeGo handle error type wrong %s", err.Error())
}
