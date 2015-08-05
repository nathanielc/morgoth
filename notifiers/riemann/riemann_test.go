package riemann_test

import (
	"flag"
	//"github.com/nathanielc/morgoth"
	//"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	//"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	//"github.com/nathanielc/morgoth/notifiers/riemann"
	"testing"
	//"time"
)

func init() {
	flag.Parse()
	if testing.Verbose() {
		flag.Set("logtostderr", "1")
	}
}

//TODO inject riemann Dialer
//func TestNotify(t *testing.T) {
//	defer glog.Flush()
//	assert := assert.New(t)
//
//	r, err := riemann.New("localhost", 5555)
//	if !assert.Nil(err) {
//		t.FailNow()
//	}
//
//	stop := time.Now()
//	start := stop.Add(-time.Hour)
//	w := &morgoth.Window{
//		Start: start,
//		Stop:  stop,
//	}
//	r.Notify("mgof", w)
//
//	assert.True(false)
//}
