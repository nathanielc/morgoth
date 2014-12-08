package ymtn_test

import (
	"flag"
	"github.com/golang/glog"
	"github.com/nvcook42/morgoth/learners/ymtn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	flag.Parse()
	if testing.Verbose() {
		flag.Set("logtostderr", "true")
		flag.Set("vmodule", "motif=2,motif_test=3")
	}
}


func TestMotif(t *testing.T) {
	assert := assert.New(t)

	size := 100
	x := make([]float64, size)
	for i := range x {
		//Saw tooth pattern that has clear change points
		x[i] = float64((i+1) % 10)
	}
	scores := ymtn.RSST(x, 5, 4)

	motifs := ymtn.DetectMotifs(x, scores, 5, 25)
	assert.NotNil(motifs)
	glog.V(1).Infoln("motifs", motifs)
}
