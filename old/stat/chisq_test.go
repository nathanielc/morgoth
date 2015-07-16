package stat_test

import (
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nathanielc/morgoth/stat"
	"testing"
)

func TestSpecificChiSqInc(t *testing.T) {

	assert := assert.New(t)

	ppf := stat.Xsquare_InvCDF(60)
	assert.Equal(79.66881012367774, ppf(0.95))
}
