package stat_test

import (
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nvcook42/morgoth/stat"
	"testing"
)

func TestSpecificChiSqInc(t *testing.T) {

	assert := assert.New(t)

	ppf := stat.Xsquare_InvCDF(60)
	assert.Equal(79.66881012367774, ppf(0.95))
}
